package cron

import (
	"context"
	"time"

	"github.com/Ptt-official-app/go-openbbsmiddleware/api"
	"github.com/Ptt-official-app/go-openbbsmiddleware/boardd"
	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/sirupsen/logrus"
)

func RetryLoadGeneralArticles(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			logrus.Infof("RetryLoadGeneralArticles: to LoadGeneralArticles")
			_ = LoadGeneralArticles()
			select {
			case <-ctx.Done():
				return nil
			default:
				logrus.Infof("RetryLoadGeneralArticles: to sleep 1 min")
				time.Sleep(1 * time.Minute)
			}
		}
	}
}

func LoadGeneralArticles() (err error) {
	nextBrdname := ""
	count := 0

	for {
		boardIDs, err := schema.GetBoardIDs(nextBrdname, false, N_BOARDS+1, true)
		if err != nil {
			return err
		}

		newNextBrdname := ""
		if len(boardIDs) == N_BOARDS+1 {
			newNextBoardID := boardIDs[N_BOARDS]
			newNextBrdname = newNextBoardID.Brdname
			boardIDs = boardIDs[:N_BOARDS]
		}

		for _, each := range boardIDs {
			err = loadGeneralArticles(each.BBoardID)
			if err == nil {
				count++
			}
		}

		if newNextBrdname == "" {
			logrus.Infof("cron.LoadGeneralArticle: load %v boards", count)
			return nil

		}

		nextBrdname = newNextBrdname
	}
}

func loadGeneralArticles(boardID bbs.BBoardID) (err error) {
	nextIdx := int32(0)
	count := 0

	for {
		articleSummaries, newNextIdx, err := loadGeneralArticlesCore(boardID, nextIdx)
		if err != nil {
			logrus.Errorf("cron.LoadGeneralArticles: unable to loadGeneralArticles: nextIdx: %v e: %v", nextIdx, err)
			return err
		}
		count += len(articleSummaries)

		logrus.Infof("cron.LoadGeneralArticles: bid: %v count: %v", boardID, count)

		if newNextIdx == INVALID_LOAD_GENERAL_ARTICLES_NEXT_IDX {
			logrus.Infof("cron.LoadGeneralArticles: bid: %v load %v articles", boardID, count)
			break
		}

		nextIdx = newNextIdx
	}

	err = loadBottomArticles(boardID)
	if err != nil {
		logrus.Errorf("loadGeneralArticles: unable to loadBottomArticles: e: %v", err)
		return err
	}

	return nil
}

func loadGeneralArticlesCore(boardID bbs.BBoardID, startIdx int32) (articleSummaries []*schema.ArticleSummaryWithRegex, nextIdx int32, err error) {
	nextIdx = INVALID_LOAD_GENERAL_ARTICLES_NEXT_IDX
	brdnameStr := boardID.ToBrdname()
	// backend load-general-articles
	ctx := context.Background()
	brdname := &boardd.BoardRef_Name{Name: brdnameStr}
	req := &boardd.ListRequest{
		Ref:          &boardd.BoardRef{Ref: brdname},
		IncludePosts: true,
		Offset:       startIdx,
		Length:       N_ARTICLES + 1,
	}
	resp, err := boardd.Cli.List(ctx, req)
	if err != nil {
		return nil, INVALID_LOAD_GENERAL_ARTICLES_NEXT_IDX, err
	}

	posts := resp.Posts
	if len(posts) == N_ARTICLES+1 {
		nextIdx = startIdx + N_ARTICLES
		posts = posts[:N_ARTICLES]
	}

	// deserialize
	updateNanoTS := types.NowNanoTS()
	articleSummaries, err = api.DeserializePBArticles(boardID, posts, updateNanoTS, false)
	if err != nil {
		return nil, INVALID_LOAD_GENERAL_ARTICLES_NEXT_IDX, err
	}

	// articleIDs
	articleIDs := make([]bbs.ArticleID, len(articleSummaries))
	for idx, each := range articleSummaries {
		articleIDs[idx] = each.ArticleID
	}

	// origArticleSummaries
	origArticleSummaries, err := schema.GetArticleSummariesByArticleIDs(boardID, articleIDs)
	if err != nil {
		return nil, INVALID_LOAD_GENERAL_ARTICLES_NEXT_IDX, err
	}

	articleSummaries, _, err = compareArticleSummaries(articleSummaries, origArticleSummaries)
	if err != nil {
		return nil, INVALID_LOAD_GENERAL_ARTICLES_NEXT_IDX, err
	}

	if len(articleSummaries) > 0 {
		err = schema.UpdateArticleSummaryWithRegexes(articleSummaries, updateNanoTS)
		if err != nil {
			return nil, INVALID_LOAD_GENERAL_ARTICLES_NEXT_IDX, err
		}
	}

	return articleSummaries, nextIdx, nil
}

func loadBottomArticles(boardID bbs.BBoardID) (err error) {
	brdnameStr := boardID.ToBrdname()
	// backend load-general-articles
	ctx := context.Background()
	brdname := &boardd.BoardRef_Name{Name: brdnameStr}
	req := &boardd.ListRequest{
		Ref:            &boardd.BoardRef{Ref: brdname},
		IncludeBottoms: true,
		Offset:         0,
		Length:         N_ARTICLES + 1,
	}
	resp, err := boardd.Cli.List(ctx, req)
	if err != nil {
		return err
	}

	updateNanoTS := types.NowNanoTS()
	articleSummaries, err := api.DeserializePBArticles(boardID, resp.Bottoms, updateNanoTS, false)
	if err != nil {
		return err
	}

	// origArticleSummaries
	origArticleSummaries, err := schema.GetBottomArticleSummaries(boardID)
	if err != nil {
		return err
	}

	articleSummaries, toRemoveArticleIDs, err := compareArticleSummaries(articleSummaries, origArticleSummaries)
	if err != nil {
		return err
	}

	if len(toRemoveArticleIDs) > 0 {
		err = schema.ResetArticleIsBottom(boardID, toRemoveArticleIDs)
		if err != nil {
			return err
		}
	}

	if len(articleSummaries) > 0 {
		err = schema.UpdateArticleSummaryWithRegexes(articleSummaries, updateNanoTS)
		if err != nil {
			return err
		}
	}

	return nil
}

func compareArticleSummaries(articleSummaries []*schema.ArticleSummaryWithRegex, origArticleSummaries []*schema.ArticleSummary) (newArticleSummaries []*schema.ArticleSummaryWithRegex, toRemoveArticleIDs []bbs.ArticleID, err error) {
	origArticleSummaryMap := make(map[bbs.ArticleID]*schema.ArticleSummary)
	for _, each := range origArticleSummaries {
		origArticleSummaryMap[each.ArticleID] = each
	}

	newArticleSummaries = make([]*schema.ArticleSummaryWithRegex, 0, len(articleSummaries))
	for _, each := range articleSummaries {
		eachOrig, ok := origArticleSummaryMap[each.ArticleID]
		if !ok {
			newArticleSummaries = append(newArticleSummaries, each)
			continue
		}
		if each.IsDeleted != eachOrig.IsDeleted {
			newArticleSummaries = append(newArticleSummaries, each)
			continue
		}
		if each.MTime != eachOrig.MTime {
			newArticleSummaries = append(newArticleSummaries, each)
			continue
		}
		if each.Recommend != eachOrig.Recommend {
			newArticleSummaries = append(newArticleSummaries, each)
			continue
		}
		if each.FullTitle != eachOrig.FullTitle {
			newArticleSummaries = append(newArticleSummaries, each)
			continue
		}
		if each.Filemode != eachOrig.Filemode {
			newArticleSummaries = append(newArticleSummaries, each)
			continue
		}
		if each.IsBottom != eachOrig.IsBottom {
			newArticleSummaries = append(newArticleSummaries, each)
			continue
		}
	}

	articleSummarySet := make(map[bbs.ArticleID]bool)
	for _, each := range articleSummaries {
		articleSummarySet[each.ArticleID] = true
	}

	toRemoveArticleIDs = make([]bbs.ArticleID, 0, len(origArticleSummaries))
	for _, each := range origArticleSummaries {
		_, ok := articleSummarySet[each.ArticleID]
		if !ok {
			toRemoveArticleIDs = append(toRemoveArticleIDs, each.ArticleID)
		}
	}

	return newArticleSummaries, toRemoveArticleIDs, nil
}
