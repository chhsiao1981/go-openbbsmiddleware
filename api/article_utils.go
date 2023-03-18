package api

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/boardd"
	"github.com/Ptt-official-app/go-openbbsmiddleware/dbcs"
	"github.com/Ptt-official-app/go-openbbsmiddleware/queue"
	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-openbbsmiddleware/utils"
	pttbbsapi "github.com/Ptt-official-app/go-pttbbs/api"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/cmsys"
	pttbbstypes "github.com/Ptt-official-app/go-pttbbs/types"
	"github.com/gin-gonic/gin"
)

func UpdateArticleContentInfo(
	boardID bbs.BBoardID,
	articleID bbs.ArticleID,
	content [][]*types.Rune,
	contentPrefix [][]*types.Rune,
	contentMD5 string,
	ip string,
	host string,
	bbs string,
	signatureMD5 string,
	signatureDBCS []byte,
	updateNanoTS types.NanoTS,
) (err error) {
	if contentMD5 == "" {
		return nil
	}

	contentID, contentBlocks := dbcs.ParseContentBlocks(boardID, articleID, content, contentMD5, updateNanoTS)

	err = schema.UpdateContentBlocks(contentBlocks, updateNanoTS)
	if err != nil {
		return err
	}

	contentInfo := &schema.ArticleContentInfo{
		ContentMD5: contentMD5,

		ContentID:     contentID,
		ContentPrefix: contentPrefix,
		IP:            ip,
		Host:          host,
		BBS:           bbs,

		SignatureDBCS: signatureDBCS,
		SignatureMD5:  signatureMD5,

		ContentUpdateNanoTS: updateNanoTS,
	}

	err = schema.UpdateArticleContentInfo(boardID, articleID, contentInfo)

	return err
}

func DeserializeArticlesAndUpdateDB(articleSummaries_b []*bbs.ArticleSummary, updateNanoTS types.NanoTS) (articleSummaries []*schema.ArticleSummaryWithRegex, err error) {
	if len(articleSummaries_b) == 0 {
		return nil, nil
	}
	articleSummaries = make([]*schema.ArticleSummaryWithRegex, len(articleSummaries_b))
	for idx, each_b := range articleSummaries_b {
		articleSummaries[idx] = schema.NewArticleSummaryWithRegex(each_b, updateNanoTS)
	}

	err = schema.UpdateArticleSummaryWithRegexes(articleSummaries, updateNanoTS)
	if err != nil {
		return nil, err
	}

	return articleSummaries, nil
}

func deserializeArticlesAndUpdateDB(userID bbs.UUserID, bboardID bbs.BBoardID, articleSummaries_b []*bbs.ArticleSummary, updateNanoTS types.NanoTS) (articleSummaries []*schema.ArticleSummaryWithRegex, userReadArticleMap map[bbs.ArticleID]bool, err error) {
	if len(articleSummaries_b) == 0 {
		return nil, nil, nil
	}

	articleSummaries, err = DeserializeArticlesAndUpdateDB(articleSummaries_b, updateNanoTS)
	if err != nil {
		return nil, nil, err
	}

	userReadArticles := make([]*schema.UserReadArticle, 0, len(articleSummaries_b))
	userReadArticleMap = make(map[bbs.ArticleID]bool)
	for _, each_b := range articleSummaries_b {
		if each_b.Read {
			each_db := &schema.UserReadArticle{
				UserID:       userID,
				ArticleID:    each_b.ArticleID,
				UpdateNanoTS: updateNanoTS,
			}
			userReadArticles = append(userReadArticles, each_db)

			userReadArticleMap[each_db.ArticleID] = true
		}
	}

	err = schema.UpdateUserReadArticles(userReadArticles, updateNanoTS)
	if err != nil {
		return nil, nil, err
	}

	// get n-comments
	updateArticleNComments(bboardID, articleSummaries)

	return articleSummaries, userReadArticleMap, err
}

func DeserializePBArticlesAndUpdateDB(boardID bbs.BBoardID, articleSummaries_b []*boardd.Post, updateNanoTS types.NanoTS, isBottom bool) (articleSummaries []*schema.ArticleSummaryWithRegex, err error) {
	if len(articleSummaries_b) == 0 {
		return nil, nil
	}
	articleSummaries = make([]*schema.ArticleSummaryWithRegex, len(articleSummaries_b))
	for idx, each_b := range articleSummaries_b {
		articleSummaries[idx] = schema.NewArticleSummaryWithRegexFromPBArticle(boardID, each_b, updateNanoTS, isBottom)
	}

	err = schema.UpdateArticleSummaryWithRegexes(articleSummaries, updateNanoTS)
	if err != nil {
		return nil, err
	}

	return articleSummaries, nil
}

func updateArticleNComments(bboardID bbs.BBoardID, articleSummaries []*schema.ArticleSummaryWithRegex) {
	if len(articleSummaries) == 0 {
		return
	}

	articleIDs := make([]bbs.ArticleID, len(articleSummaries))
	for idx, each := range articleSummaries {
		articleIDs[idx] = each.ArticleID
	}

	articleNComments, err := schema.GetArticleNCommentsByArticleIDs(bboardID, articleIDs)
	if err != nil {
		return
	}

	nCommentsByArticleIDMap := make(map[bbs.ArticleID]*schema.ArticleNComments)
	for _, each := range articleNComments {
		nCommentsByArticleIDMap[each.ArticleID] = each
	}

	for _, each := range articleSummaries {
		eachArticleNComments := nCommentsByArticleIDMap[each.ArticleID]
		if eachArticleNComments == nil {
			continue
		}

		each.NComments = eachArticleNComments.NComments
		each.Rank = eachArticleNComments.Rank
	}
}

func ArticleLockKey(boardID bbs.BBoardID, articleID bbs.ArticleID) (key string) {
	return "a:" + string(boardID) + ":" + string(articleID)
}

// TryGetArticleContentInfo
//
// 嘗試拿到 article-content
//
//  1. 根據 article-id 得到相對應的 filename, ownerid, create-time.
//  2. 嘗試從 schema 拿到 db summary 資訊. (create-time)
//  3. 如果可以從 schema 拿到 db 資訊:
//     3.1. 如果已經 deleted: return deleted.
//     3.2. 如果距離上次跟 pttbbs 拿的時間太近: 從 schema 拿到 content, return schema-content.
//  4. 嘗試做 lock.
//     4.1. 如果 lock 失敗: 從 schema 拿到 content, return schema-content.
//  5. 從 pttbbs 拿到 article
//  6. 如果從 pttbbs 拿到的時間比 schema 裡拿到的時間舊的話: return schema-content.
//  7. parse article 為 content / comments.
//  8. 將 comments parse 為 firstComments / theRestComments.
//  9. 將 theRestComments 丟進 queue 裡.
func TryGetArticleContentInfo(userID bbs.UUserID, bboardID bbs.BBoardID, articleID bbs.ArticleID, c *gin.Context, isSystem bool, isHash bool, isContent bool) (content [][]*types.Rune, contentPrefix [][]*types.Rune, contentMD5 string, ip string, host string, bbs string, signatureMD5 string, signatureDBCSByte []byte, articleDetailSummary *schema.ArticleDetailSummary, fileSize int, hash cmsys.Fnv64_t, statusCode int, err error) {
	updateNanoTS := types.NanoTS(0)
	// set user-read-article-id
	defer func() {
		if err == nil && !isSystem {
			setUserReadArticle(content, userID, articleID, updateNanoTS)
		}
	}()

	isForce := false
	isQueue := true

	// if isHash => force to receive the new article
	//              and re-calc hash. (and no queue.)
	if isHash {
		isForce = true
		isQueue = false
	}

	// if isSystem (cron) => no queue.
	if isSystem {
		isQueue = false
	}

	// get article-info (ensuring valid article-id)
	articleFilename := articleID.ToRaw()
	articleCreateTime, err := articleFilename.CreateTime()
	if err != nil {
		return nil, nil, "", "", "", "", "", nil, nil, 0, 0, 500, err
	}

	articleCreateTimeNanoTS := types.Time4ToNanoTS(articleCreateTime)

	// get from backend with content-mtime
	// estimated max 500ms + 3 seconds
	articleDetailSummary, statusCode, err = tryGetArticleDetailSummary(userID, bboardID, articleID, articleCreateTime, c, isSystem)
	if err != nil {
		return nil, nil, "", "", "", "", "", nil, nil, 0, 0, statusCode, err
	}

	// preliminarily checking article-detail-summary.
	if articleDetailSummary.IsDeleted {
		return nil, nil, "", "", "", "", "", nil, nil, 500, 0, 0, ErrAlreadyDeleted
	}

	if articleDetailSummary.CreateTime == 0 {
		articleDetailSummary.CreateTime = articleCreateTimeNanoTS
	}

	// already got the most updated content.
	if !isForce && tryGetArticleContentInfoTooSoon(articleDetailSummary.ContentUpdateNanoTS) {
		contentInfo, err := schema.GetArticleContentInfo(bboardID, articleID, isContent)
		if err != nil {
			return nil, nil, "", "", "", "", "", nil, nil, 0, 0, 500, err
		}
		return contentInfo.Content, contentInfo.ContentPrefix, contentInfo.ContentMD5, contentInfo.IP, contentInfo.Host, contentInfo.BBS, contentInfo.SignatureMD5, contentInfo.SignatureDBCS, articleDetailSummary, 0, 0, 200, nil
	}

	ownerID := articleDetailSummary.Owner

	// 4. do lock. if failed, return the data in db.
	lockKey := ArticleLockKey(bboardID, articleID)
	err = schema.TryLock(lockKey, ARTICLE_LOCK_TS_DURATION)
	if err != nil {
		if isForce {
			return nil, nil, "", "", "", "", "", nil, nil, 0, 0, 500, err
		}

		// unable to do lock. get the most updated content.
		contentInfo, err := schema.GetArticleContentInfo(bboardID, articleID, isContent)
		if err != nil {
			return nil, nil, "", "", "", "", "", nil, nil, 0, 0, 500, err
		}
		updateNanoTS = types.NowNanoTS()
		return contentInfo.Content, contentInfo.ContentPrefix, contentInfo.ContentMD5, contentInfo.IP, contentInfo.Host, contentInfo.BBS, contentInfo.SignatureMD5, contentInfo.SignatureDBCS, articleDetailSummary, 0, 0, 200, nil
	}
	defer func() { _ = schema.Unlock(lockKey) }()

	// 5. get article from pttbbs
	theParams_b := &pttbbsapi.GetArticleParams{
		RetrieveTS: articleDetailSummary.ContentMTime.ToTime4(),
		IsSystem:   isSystem,
		IsHash:     isHash,
	}
	var result_b *pttbbsapi.GetArticleResult

	urlMap := map[string]string{
		"bid": string(bboardID),
		"aid": string(articleID),
	}

	url := utils.MergeURL(urlMap, pttbbsapi.GET_ARTICLE_R)
	statusCode, err = utils.BackendGet(c, url, theParams_b, nil, &result_b)
	if err != nil {
		return nil, nil, "", "", "", "", "", nil, nil, 0, 0, statusCode, err
	}

	fileSize = len(result_b.Content)
	hash = result_b.Hash

	// 6. check content-mtime (no modify from backend, no need to parse again)
	contentMTime := types.Time4ToNanoTS(result_b.MTime)
	if articleDetailSummary.ContentMTime >= contentMTime {
		contentInfo, err := schema.GetArticleContentInfo(bboardID, articleID, isContent)
		if err != nil {
			return nil, nil, "", "", "", "", "", nil, nil, 0, 0, 500, err
		}
		return contentInfo.Content, contentInfo.ContentPrefix, contentInfo.ContentMD5, contentInfo.IP, contentInfo.Host, contentInfo.BBS, contentInfo.SignatureMD5, contentInfo.SignatureDBCS, articleDetailSummary, 0, 0, 200, nil
	}

	if result_b.Content == nil { // XXX possibly the article is deleted. Need to check error-code and mark the article as deleted.
		return nil, nil, "", "", "", "", "", nil, nil, 0, 0, 500, ErrNoArticle
	}

	// 7. parse article as content / commentsDBCS
	updateNanoTS = types.NowNanoTS()

	content, contentPrefix, contentMD5, ip, host, bbs, signatureMD5, signatureDBCS, commentsDBCS := dbcs.ParseContent(result_b.Content, articleDetailSummary.ContentMD5)

	signatureDBCSByte = signatureDBCS

	// update article
	// we need update-article-content be the 1st to upload,
	// because it's possible that there is no first-comments.
	// only article-content is guaranteed.

	err = UpdateArticleContentInfo(bboardID, articleID, content, contentPrefix, contentMD5, ip, host, bbs, signatureMD5, signatureDBCSByte, updateNanoTS)

	if err != nil {
		return nil, nil, "", "", "", "", "", nil, nil, 0, 0, 500, err
	}

	if contentMD5 == "" {
		contentInfo, err := schema.GetArticleContentInfo(bboardID, articleID, isContent)
		if err != nil {
			return nil, nil, "", "", "", "", "", nil, nil, 0, 0, 500, err
		}
		content = contentInfo.Content
		contentMD5 = contentInfo.ContentMD5
		contentPrefix = contentInfo.ContentPrefix
		ip = contentInfo.IP
		host = contentInfo.Host
		bbs = contentInfo.BBS
		signatureMD5 = contentInfo.SignatureMD5
		signatureDBCSByte = contentInfo.SignatureDBCS
	}

	if isQueue {
		// 8. parse comments as firstComments and theRestComments
		firstComments, firstCommentsMD5, _, err := dbcs.ParseFirstComments(
			bboardID,
			articleID,
			ownerID,
			articleCreateTimeNanoTS,
			contentMTime,
			commentsDBCS,
			articleDetailSummary.FirstCommentsMD5,
		)
		if err != nil {
			return content, contentPrefix, contentMD5, ip, host, bbs, signatureMD5, signatureDBCSByte, articleDetailSummary, fileSize, hash, 200, nil
		}

		// update first-comments
		// possibly err because the data is too old.
		// we don't need to queue and update content-mtime if the data is too old.
		err = tryUpdateFirstComments(firstComments, firstCommentsMD5, updateNanoTS, articleDetailSummary)
		if err != nil {
			//if failed update: we still send the content back.
			//(no updating the content in db,
			// so the data will be re-processed again next time).
			return content, contentPrefix, contentMD5, ip, host, bbs, signatureMD5, signatureDBCSByte, articleDetailSummary, fileSize, hash, 200, nil
		}

		// 9. enqueue and n_comments
		err = queue.QueueCommentDBCS(bboardID, articleID, ownerID, commentsDBCS, articleCreateTimeNanoTS, contentMTime, updateNanoTS)
		if err != nil {
			return content, contentPrefix, contentMD5, ip, host, bbs, signatureMD5, signatureDBCSByte, articleDetailSummary, fileSize, hash, 200, nil
		}

		if articleDetailSummary.NComments == 0 {
			articleDetailSummary.NComments = len(firstComments)
		}
	} else {
		commentQueue := &queue.CommentQueue{
			BBoardID:          bboardID,
			ArticleID:         articleID,
			OwnerID:           ownerID,
			CommentDBCS:       commentsDBCS,
			ArticleCreateTime: articleCreateTimeNanoTS,
			ArticleMTime:      contentMTime,
			UpdateNanoTS:      updateNanoTS,
		}

		_ = queue.ProcessCommentQueue(commentQueue)
	}

	// everything is good, update content-mtime
	_ = schema.UpdateArticleContentMTime(bboardID, articleID, contentMTime)

	return content, contentPrefix, contentMD5, ip, host, bbs, signatureMD5, signatureDBCSByte, articleDetailSummary, fileSize, hash, 200, nil
}

func tryGetArticleContentInfoTooSoon(updateNanoTS types.NanoTS) bool {
	nowNanoTS := types.NowNanoTS()
	return nowNanoTS-updateNanoTS < GET_ARTICLE_CONTENT_INFO_TOO_SOON_NANO_TS
}

func tryGetArticleDetailSummary(userID bbs.UUserID, boardID bbs.BBoardID, articleID bbs.ArticleID, articleCreateTime pttbbstypes.Time4, c *gin.Context, isSystem bool) (articleDetailSummary *schema.ArticleDetailSummary, statusCode int, err error) {
	articleDetailSummary, err = schema.GetArticleDetailSummary(boardID, articleID)
	if err != nil { // something went wrong with db.
		return nil, 500, err
	}
	if articleDetailSummary == nil {
		return nil, 500, ErrNoArticle
	}

	return articleDetailSummary, 200, nil
}

func setUserReadArticle(content [][]*types.Rune, userID bbs.UUserID, articleID bbs.ArticleID, updateNanoTS types.NanoTS) {
	if content == nil {
		return
	}

	// user read article
	userReadArticle := &schema.UserReadArticle{
		UserID:       userID,
		ArticleID:    articleID,
		UpdateNanoTS: updateNanoTS,
	}
	_ = schema.UpdateUserReadArticle(userReadArticle)
}

func editArticleGetArticleContentInfo(userID bbs.UUserID, boardID bbs.BBoardID, articleID bbs.ArticleID, c *gin.Context, isContent bool) (oldContent [][]*types.Rune, oldContentPrefix [][]*types.Rune, signatureDBCS []byte, articleDetailSummary_db *schema.ArticleDetailSummary, sz int, hash cmsys.Fnv64_t, statusCode int, err error) {
	oldContent, oldContentPrefix, _, _, _, _, _, signatureDBCS, articleDetailSummary_db, sz, hash, statusCode, err = TryGetArticleContentInfo(userID, boardID, articleID, c, false, true, isContent)

	return oldContent, oldContentPrefix, signatureDBCS, articleDetailSummary_db, sz, hash, statusCode, err
}
