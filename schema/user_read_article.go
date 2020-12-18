package schema

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/db"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

var (
	UserReadArticle_c *db.Collection
)

type UserReadArticle struct {
	//已讀文章紀錄

	UserID       bbs.UUserID   `bson:"user_id"`
	ArticleID    bbs.ArticleID `bson:"aid"`
	UpdateNanoTS types.NanoTS  `bson:"update_nano_ts"`
}

const USER_READ_ARTICLE_USER_ID_b = "user_id"
const USER_READ_ARTICLE_ARTICLE_ID_b = "aid"
