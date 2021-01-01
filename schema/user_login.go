package schema

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"go.mongodb.org/mongo-driver/bson"
)

type UserLogin struct {
	UserID    bbs.UUserID  `bson:"user_id"`
	LastLogin types.NanoTS `bson:"last_login_ts"`
}

var (
	EMPTY_USER_SUMMARY = &UserLogin{}
	userLoginFields    = getFields(EMPTY_USER, EMPTY_USER_SUMMARY)
)

func UpdateUserLogin(userLogin *UserLogin) (err error) {

	query := bson.M{
		USER_USER_ID_b: userLogin.UserID,
	}

	r, err := User_c.CreateOnly(query, userLogin)
	if err != nil {
		return err
	}
	if r.UpsertedCount > 0 {
		return nil
	}

	query = bson.M{
		"$or": bson.A{
			bson.M{
				USER_USER_ID_b: userLogin.UserID,
				USER_LAST_LOGIN_b: bson.M{
					"$exists": false,
				},

				USER_IS_DELETED_b: bson.M{"$exists": false},
			},
			bson.M{
				USER_USER_ID_b: userLogin.UserID,
				USER_LAST_LOGIN_b: bson.M{
					"$lt": userLogin.LastLogin,
				},

				USER_IS_DELETED_b: bson.M{"$exists": false},
			},
		},
	}

	r, err = User_c.UpdateOneOnly(query, userLogin)
	if err != nil {
		return err
	}
	if r.MatchedCount == 0 {
		return ErrNoMatch
	}
	return nil
}
