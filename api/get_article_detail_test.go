package api

import (
	"context"
	"sort"
	"sync"
	"testing"

	"github.com/Ptt-official-app/go-openbbsmiddleware/apitypes"
	"github.com/Ptt-official-app/go-openbbsmiddleware/boardd"
	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/testutil"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetArticleDetail(t *testing.T) {
	setupTest()
	defer teardownTest()

	boardSummaries_b := []*bbs.BoardSummary{testBoardSummaryWhoAmI_b}
	_, _, _ = deserializeBoardsAndUpdateDB("SYSOP", boardSummaries_b, 123456890000000000)

	// load articles
	ctx := context.Background()
	brdname := &boardd.BoardRef_Name{Name: "WhoAmI"}
	req := &boardd.ListRequest{
		Ref:          &boardd.BoardRef{Ref: brdname},
		IncludePosts: true,
		Offset:       0,
		Length:       100 + 1,
	}
	resp, _ := boardd.Cli.List(ctx, req)

	posts := resp.Posts

	logrus.Infof("TestGetArticleDetail: posts: %v", len(posts))

	updateNanoTS := types.NowNanoTS()
	_, _ = DeserializePBArticlesAndUpdateDB("10_WhoAmI", posts, updateNanoTS, false)

	// params
	params := &GetArticleDetailParams{}
	path0 := &GetArticleDetailPath{
		FBoardID:   apitypes.FBoardID("WhoAmI"),
		FArticleID: apitypes.FArticleID("M.1608386280.A.BC9"),
	}

	expectedResult0 := &GetArticleDetailResult{
		BBoardID:   apitypes.FBoardID("WhoAmI"),
		ArticleID:  apitypes.FArticleID("M.1608386280.A.BC9"),
		Owner:      bbs.UUserID("SYSOP"),
		CreateTime: types.Time8(1608386280),
		MTime:      types.Time8(1608386280),

		URL:  "http://localhost:3457/bbs/board/WhoAmI/article/M.1608386280.A.BC9",
		Read: true,

		Title:         "測試一下特殊字～",
		Money:         0,
		Recommend:     9,
		Class:         "心得",
		Brdname:       "WhoAmI",
		Content:       testContent3Utf8[4:],
		ContentPrefix: testContent3Utf8[:4],
		IP:            "172.22.0.1",
		Host:          "",
		BBS:           "批踢踢 docker(pttdocker.test)",
	}

	expectedArticleDetailSummary0 := &schema.ArticleDetailSummary{
		BBoardID:   bbs.BBoardID("10_WhoAmI"),
		ArticleID:  bbs.ArticleID("1VtWRel9"),
		ContentMD5: "L6QISYJFt-Y5g4Thl-roaw",
		// ContentMTime:          types.NanoTS(1608386280000000000),
		FirstCommentsLastTime: types.NanoTS(0),

		CreateTime: types.NanoTS(1608386280000000000),
		MTime:      types.NanoTS(1608386280000000000),
		Recommend:  9,
		Owner:      "SYSOP",
		Title:      "測試一下特殊字～",
		Money:      0,
		Class:      "心得",

		IP:  "172.22.0.1",
		BBS: "批踢踢 docker(pttdocker.test)",

		Idx: "1608386280@1VtWRel9",
	}

	expectedArticleDetailSummary02 := &schema.ArticleDetailSummary{
		BBoardID:   bbs.BBoardID("10_WhoAmI"),
		ArticleID:  bbs.ArticleID("1VtWRel9"),
		ContentMD5: "L6QISYJFt-Y5g4Thl-roaw",
		// ContentMTime:          types.NanoTS(1608386300000000000),
		FirstCommentsLastTime: types.NanoTS(0),

		CreateTime: types.NanoTS(1608386280000000000),
		MTime:      types.NanoTS(1608386280000000000),
		Recommend:  9,
		Owner:      "SYSOP",
		Title:      "測試一下特殊字～",
		Money:      0,
		Class:      "心得",

		IP:  "172.22.0.1",
		BBS: "批踢踢 docker(pttdocker.test)",

		Idx: "1608386280@1VtWRel9",
	}

	path1 := &GetArticleDetailPath{
		FBoardID:   apitypes.FBoardID("WhoAmI"),
		FArticleID: apitypes.FArticleID("M.1608388506.A.85D"),
	}

	expectedResult1 := &GetArticleDetailResult{
		BBoardID:   apitypes.FBoardID("WhoAmI"),
		ArticleID:  apitypes.FArticleID("M.1608388506.A.85D"),
		Owner:      bbs.UUserID("SYSOP"),
		CreateTime: types.Time8(1608388506),
		MTime:      types.Time8(1608386280),

		URL:  "http://localhost:3457/bbs/board/WhoAmI/article/M.1608388506.A.85D",
		Read: true,

		Title:         "所以特殊字真的是有綠色的∼",
		Money:         0,
		Recommend:     9,
		Class:         "閒聊",
		Brdname:       "WhoAmI",
		Content:       testContent4Utf8[4:],
		ContentPrefix: testContent4Utf8[:4],
		IP:            "172.22.0.1",
		Host:          "",
		BBS:           "批踢踢 docker(pttdocker.test)",
		NComments:     3,
	}

	expectedArticleDetailSummary1 := &schema.ArticleDetailSummary{
		BBoardID:  bbs.BBoardID("10_WhoAmI"),
		ArticleID: bbs.ArticleID("1VtW-QXT"),
		// ContentMTime: types.NanoTS(1608388624000000000),
		ContentMD5: "riiRuKCZzG0gAGpQiq4GJA",
		Owner:      "SYSOP",

		FirstCommentsMD5: "3fjMk__1yvzpuEgq8jfdmg",
		NComments:        0,

		CreateTime: types.NanoTS(1608388506000000000),
		MTime:      types.NanoTS(1608386280000000000),

		Title:     "所以特殊字真的是有綠色的∼",
		Money:     0,
		Recommend: 9,
		Class:     "閒聊",

		Host: "",
		IP:   "172.22.0.1",
		BBS:  "批踢踢 docker(pttdocker.test)",

		Idx: "1608388506@1VtW-QXT",
	}
	c := &gin.Context{}
	// c.Request = &http.Request{URL: &url.URL{Path: "/api/boards/WhoAmI/article/M.1607937174.A.081"}}
	type args struct {
		remoteAddr string
		userID     bbs.UUserID
		params     interface{}
		path       interface{}
		c          *gin.Context
		boardID    bbs.BBoardID
		articleID  bbs.ArticleID
	}
	tests := []struct {
		name                         string
		args                         args
		expectedResult               interface{}
		expectedFirstComments        []*schema.Comment
		expectedArticleDetailSummary *schema.ArticleDetailSummary

		expectedStatusCode int
		wantErr            bool
		toSoonNanoTS       types.NanoTS
	}{
		// TODO: Add test cases.
		{
			name: "0th-BC9",
			args: args{
				remoteAddr: "localhost",
				userID:     "chhsiao123",
				params:     params,
				path:       path0,
				c:          c,
				boardID:    "10_WhoAmI",
				articleID:  "1VtWRel9",
			},
			expectedArticleDetailSummary: expectedArticleDetailSummary0,
			expectedResult:               expectedResult0,
			expectedStatusCode:           200,
		},
		{
			name: "1st-BC9",
			args: args{
				remoteAddr: "localhost",
				userID:     "chhsiao123",
				params:     params,
				path:       path0,
				c:          c,
				boardID:    "10_WhoAmI",
				articleID:  "1VtWRel9",
			},
			expectedArticleDetailSummary: expectedArticleDetailSummary0,
			expectedResult:               expectedResult0,
			expectedStatusCode:           200,
		},
		{
			name: "2st-BC9",
			args: args{
				remoteAddr: "localhost",
				userID:     "chhsiao123",
				params:     params,
				path:       path0,
				c:          c,
				boardID:    "10_WhoAmI",
				articleID:  "1VtWRel9",
			},
			expectedArticleDetailSummary: expectedArticleDetailSummary02,
			expectedResult:               expectedResult0,
			expectedStatusCode:           200,
			toSoonNanoTS:                 1,
		},
		/*
			{
				name: "0th-081",
				args: args{
					remoteAddr: "localhost",
					userID:     "chhsiao123",
					params:     params,
					path:       path1,
					c:          c,
					boardID:    "10_WhoAmI",
					articleID:  "1VrooM21",
				},
				expectedFirstComments:        testUtf8FullFirstComments5,
				expectedResult:               expectedResult1,
				expectedStatusCode:           200,
				expectedArticleDetailSummary: expectedArticleDetailSummary1,
			},
		*/
		{
			name: "1st-1VtW-QXT",
			args: args{
				remoteAddr: "localhost",
				userID:     "chhsiao123",
				params:     params,
				path:       path1,
				c:          c,
				boardID:    "10_WhoAmI",
				articleID:  "1VtW-QXT",
			},
			expectedResult:               expectedResult1,
			expectedStatusCode:           200,
			expectedArticleDetailSummary: expectedArticleDetailSummary1,
			expectedFirstComments:        testFullFirstComments4,
		},
	}

	var wg sync.WaitGroup
	for _, tt := range tests {
		wg.Add(1)
		t.Run(tt.name, func(t *testing.T) {
			defer wg.Done()
			origTooSoonNanoTS := GET_ARTICLE_CONTENT_INFO_TOO_SOON_NANO_TS
			defer func() {
				GET_ARTICLE_CONTENT_INFO_TOO_SOON_NANO_TS = origTooSoonNanoTS
			}()

			if tt.toSoonNanoTS != 0 {
				GET_ARTICLE_CONTENT_INFO_TOO_SOON_NANO_TS = tt.toSoonNanoTS
			}

			gotResult, gotStatusCode, err := GetArticleDetail(tt.args.remoteAddr, tt.args.userID, tt.args.params, tt.args.path, tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArticleDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			testutil.TDeepEqual(t, "result", gotResult, tt.expectedResult)

			if gotStatusCode != tt.expectedStatusCode {
				t.Errorf("GetArticleDetail() gotStatusCode = %v, want %v", gotStatusCode, tt.expectedStatusCode)
			}

			query := bson.M{
				schema.COMMENT_BBOARD_ID_b:  tt.args.boardID,
				schema.COMMENT_ARTICLE_ID_b: tt.args.articleID,
			}
			var gotComments []*schema.Comment
			_ = schema.Comment_c.Find(query, 0, &gotComments, nil, nil)

			for _, each := range gotComments {
				each.UpdateNanoTS = 0
			}

			sort.SliceStable(gotComments, func(i, j int) bool {
				return gotComments[i].SortTime <= gotComments[j].SortTime
			})

			testutil.TDeepEqual(t, "comments", gotComments, tt.expectedFirstComments)

			gotArticleDetailSummary, err := schema.GetArticleDetailSummary(tt.args.boardID, tt.args.articleID)
			logrus.Infof("GetArticleDetail: after GetArticleDetailSummary: e: %v", err)
			if gotArticleDetailSummary != nil {
				gotArticleDetailSummary.ContentUpdateNanoTS = 0
				gotArticleDetailSummary.FirstCommentsUpdateNanoTS = 0
				gotArticleDetailSummary.CommentsUpdateNanoTS = 0
				gotArticleDetailSummary.NComments = 0
				gotArticleDetailSummary.UpdateNanoTS = 0

				gotArticleDetailSummary.ContentID = ""
				gotArticleDetailSummary.ContentMTime = 0
			}
			testutil.TDeepEqual(t, "article-detail-summary", gotArticleDetailSummary, tt.expectedArticleDetailSummary)
		})
		wg.Wait()
	}
}
