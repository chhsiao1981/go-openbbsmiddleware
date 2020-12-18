package api

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"

	"github.com/Ptt-official-app/go-pttbbs/types/proto"
)

const GET_ARTICLE_R = "/board/:bid/article/:aid"

type GetArticleDetailParams struct {
	Fields string `json:"fields,omitempty" form:"fields,omitempty" url:"fields,omitempty"`
}

type GetArticleDetailPath struct {
	BBoardID  bbs.BBoardID  `uri:"bid"`
	ArticleID bbs.ArticleID `uri:"aid"`
}

type GetArticleDetailResult struct {
	*types.ArticleSummary

	Brdname string        `json:"brdname"`
	Content proto.Content `json:"content"`
	IP      string        `json:"ip"`
	Host    string        `json:"host"` //ip 的中文呈現, 外國則為國家.
	BBS     string        `json:"bbs"`
}

func GetArticleDetail(remoteAddr string, userID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, statusCode int, err error) {
	thePath, ok := path.(*GetArticleDetailPath)
	if !ok {
		return nil, 400, ErrInvalidParams
	}

	result, err = tryGetArticleContentInfo(userID, thePath.BBoardID, thePath.ArticleID)
	if err != nil {
		return nil, 400, err
	}

	return result, 200, nil
}

func tryGetArticleContentInfo(userID bbs.UUserID, bboardID bbs.BBoardID, articleID bbs.ArticleID) (result *GetArticleDetailResult, err error) {

	return nil, nil
}
