package api

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/apitypes"
	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

const LOAD_USER_COMMENTS_R = "/user/:user_id/comments"

type LoadUserCommentsParams struct {
	StartIdx   string `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	Descending bool   `json:"desc,omitempty"  form:"desc,omitempty" url:"desc,omitempty"`
	Max        int    `json:"limit,omitempty" form:"limit,omitempty" url:"limit,omitempty"`
}

type LoadUserCommentsPath struct {
	UserID bbs.UUserID `uri:"user_id"`
}

type LoadUserCommentsResult struct {
	List    []*apitypes.ArticleComment `json:"list"`
	NextIdx string                     `json:"next_idx"`
}

func NewLoadUserCommentsParams() *LoadUserCommentsParams {
	return &LoadUserCommentsParams{
		Descending: DEFAULT_DESCENDING,
		Max:        DEFAULT_MAX_LIST,
	}
}

func LoadUserCommentsWrapper(c *gin.Context) {
	params := NewLoadUserCommentsParams()
	path := &LoadUserCommentsPath{}
	LoginRequiredPathQuery(LoadUserComments, params, path, c)
}

func LoadUserComments(remoteAddr string, userID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, statusCode int, err error) {
	theParams, ok := params.(*LoadUserCommentsParams)
	if !ok {
		return nil, 400, ErrInvalidParams
	}

	thePath, ok := path.(*LoadUserCommentsPath)
	if !ok {
		return nil, 400, ErrInvalidPath
	}

	commentSummaries_db, nextIdx, err := loadUserComments(thePath.UserID, theParams.StartIdx, theParams.Descending, theParams.Max)
	if err != nil {
		return nil, 500, err
	}

	articleSummaryMap, userReadArticleMap, err := getArticleMapFromCommentSummaries(userID, commentSummaries_db)
	if err != nil {
		return nil, 500, err
	}

	r := NewLoadUserCommentsResult(commentSummaries_db, articleSummaryMap, userReadArticleMap, nextIdx)

	return r, 200, nil
}

func loadUserComments(ownerID bbs.UUserID, startIdx string, descending bool, max int) (commentSummaries_db []*schema.CommentSummary, nextIdx string, err error) {
	commentSummaries_db = make([]*schema.CommentSummary, 0, max+1)

	var nextSortTime types.NanoTS
	if startIdx != "" {
		_, startIdx, err = apitypes.DeserializeArticleCommentIdx(startIdx)
		if err != nil {
			return nil, "", err
		}
		nextSortTime, _ = apitypes.DeserializeCommentIdx(startIdx)
	}

	isEndLoop := false
	remaining := max
	for !isEndLoop && remaining > 0 {
		eachCommentSummaries_db, err := schema.GetBasicCommentSummariesByOwnerID(ownerID, nextSortTime, descending, max+1)
		if err != nil {
			return nil, "", err
		}

		// check is-last query
		if len(eachCommentSummaries_db) < max+1 {
			isEndLoop = true
			nextSortTime = 0
		} else {
			// setup next
			nextCommentSummary := eachCommentSummaries_db[len(eachCommentSummaries_db)-1]
			eachCommentSummaries_db = eachCommentSummaries_db[:len(eachCommentSummaries_db)-1]

			nextSortTime = nextCommentSummary.SortTime
			nextIdx = apitypes.SerializeCommentIdx(nextSortTime, nextCommentSummary.CommentID)
		}
		if len(eachCommentSummaries_db) == 0 {
			break
		}

		// is-valid
		validCommentSummaries_db, err := isValidCommentSummaries(eachCommentSummaries_db)
		if err != nil {
			return nil, "", err
		}

		// append
		if len(validCommentSummaries_db) > remaining {
			nextCommentSummary := validCommentSummaries_db[remaining]
			validCommentSummaries_db = validCommentSummaries_db[:remaining]

			nextSortTime = nextCommentSummary.SortTime
			nextIdx = apitypes.SerializeCommentIdx(nextSortTime, nextCommentSummary.CommentID)
		}

		commentSummaries_db = append(commentSummaries_db, validCommentSummaries_db...)
		remaining -= len(validCommentSummaries_db)
	}

	return commentSummaries_db, nextIdx, nil
}

// isValidCommentSummaries
// XXX TODO
func isValidCommentSummaries(commentSummaries_db []*schema.CommentSummary) ([]*schema.CommentSummary, error) {
	return commentSummaries_db, nil
}

func getArticleMapFromCommentSummaries(userID bbs.UUserID, commentSummaries_db []*schema.CommentSummary) (articleSummaryMap map[bbs.BBoardID]map[bbs.ArticleID]*schema.ArticleSummary, userReadArticleMap map[bbs.BBoardID]map[bbs.ArticleID]types.NanoTS, err error) {
	articleIDMap := make(map[bbs.BBoardID][]bbs.ArticleID)
	for _, each := range commentSummaries_db {
		_, ok := articleIDMap[each.BBoardID]
		if !ok {
			articleIDMap[each.BBoardID] = make([]bbs.ArticleID, 0, len(commentSummaries_db))
		}

		articleIDMap[each.BBoardID] = append(articleIDMap[each.BBoardID], each.ArticleID)
	}

	// article summary map

	articleSummaryMap = make(map[bbs.BBoardID]map[bbs.ArticleID]*schema.ArticleSummary)
	var articleSummaries []*schema.ArticleSummary
	for boardID, articleIDs := range articleIDMap {
		articleSummaries, err = schema.GetArticleSummariesByArticleIDs(boardID, articleIDs)
		if err != nil {
			continue
		}
		eachArticleSummaryMap := make(map[bbs.ArticleID]*schema.ArticleSummary)
		for _, each := range articleSummaries {
			eachArticleSummaryMap[each.ArticleID] = each
		}
		articleSummaryMap[boardID] = eachArticleSummaryMap
	}

	// user read article map
	userReadArticleMap = make(map[bbs.BBoardID]map[bbs.ArticleID]types.NanoTS)
	for boardID, articleIDs := range articleIDMap {
		userReadArticles, err := schema.FindUserReadArticlesByArticleIDs(userID, boardID, articleIDs)
		if err != nil {
			continue
		}
		eachUserReadArticleMap := make(map[bbs.ArticleID]types.NanoTS)
		for _, each := range userReadArticles {
			eachUserReadArticleMap[each.ArticleID] = each.UpdateNanoTS
		}
		userReadArticleMap[boardID] = eachUserReadArticleMap
	}

	return articleSummaryMap, userReadArticleMap, nil
}

func NewLoadUserCommentsResult(
	commentSummaries_db []*schema.CommentSummary,
	articleSummaryMap map[bbs.BBoardID]map[bbs.ArticleID]*schema.ArticleSummary,
	userReadArticleMap map[bbs.BBoardID]map[bbs.ArticleID]types.NanoTS,
	nextIdx string,
) (result *LoadUserCommentsResult) {
	comments := make([]*apitypes.ArticleComment, len(commentSummaries_db))
	for idx, each := range commentSummaries_db {
		articleSummaryMapByBoardID, ok := articleSummaryMap[each.BBoardID]
		if !ok {
			continue
		}
		articleSummary, ok := articleSummaryMapByBoardID[each.ArticleID]
		if !ok {
			continue
		}
		comments[idx] = apitypes.NewArticleCommentFromComment(articleSummary, each)

		// read
		userReadArticleMapByBoardID, ok := userReadArticleMap[each.BBoardID]
		if !ok {
			continue
		}
		readNanoTS, ok := userReadArticleMapByBoardID[each.ArticleID]
		if !ok {
			continue
		}

		if readNanoTS > articleSummary.MTime {
			comments[idx].Read = types.READ_STATUS_MTIME
		} else if readNanoTS > each.SortTime {
			comments[idx].Read = types.READ_STATUS_COMMENT_TIME
		} else if readNanoTS > each.CreateTime {
			comments[idx].Read = types.READ_STATUS_CREATE_TIME
		} else {
			comments[idx].Read = types.READ_STATUS_UNREAD
		}
	}

	nextIdx = apitypes.SerializeArticleCommentIdx(apitypes.ARTICLE_COMMENT_TYPE_COMMENT, nextIdx)

	result = &LoadUserCommentsResult{
		List:    comments,
		NextIdx: nextIdx,
	}
	return result
}
