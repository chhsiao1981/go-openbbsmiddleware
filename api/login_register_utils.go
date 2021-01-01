package api

import (
	"time"

	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/gin-gonic/gin"
)

func setTokenToCookie(c *gin.Context, accessToken string) {
	if c == nil {
		return
	}

	expires := time.Now().Add(3 * 86400 * time.Second)

	expiresStr := expires.Format("Mon, Jan 2 2006 15:04:05 MST")
	c.Header("Set-Cookie", "token="+accessToken+";Expires="+expiresStr+";SameSite=Strict;"+types.ACCESS_TOKEN_COOKIE_SUFFIX)
}

func checkIsPasswdVulnerable(passwd string) bool {
	return len(passwd) < MIN_LEN_PASSWD
}

func serializeUserLoginAndUpdateDB(userID bbs.UUserID, updateNanoTS types.NanoTS) (err error) {

	userLogin := &schema.UserLogin{
		UserID:    userID,
		LastLogin: updateNanoTS,
	}

	return schema.UpdateUserLogin(userLogin)
}

func serializeUserIsPasswdVulnerableAddUpdateDB(userID bbs.UUserID, updateNanoTS types.NanoTS) (err error) {

	userIsPasswdVulnerable := &schema.UserIsPasswdVulnerable{
		UserID:                         userID,
		IsPasswdVulnerable:             true,
		IsPasswdVulnerableUpdateNanoTS: updateNanoTS,
	}

	return schema.UpdateUserIsPasswdVulnerable(userIsPasswdVulnerable)
}
