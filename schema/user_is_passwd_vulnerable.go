package schema

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"go.mongodb.org/mongo-driver/bson"
)

type UserIsPasswdVulnerable struct {
	UserID                         bbs.UUserID  `bson:"user_id"`
	IsPasswdVulnerable             bool         `bson:"pw_vulnerable"`
	IsPasswdVulnerableUpdateNanoTS types.NanoTS `bson:"pw_vulnerable_nano_ts"`
}

var (
	EMPTY_USER_IS_PASSWD_VULNERABLE = &UserIsPasswdVulnerable{}
	userIsPasswdVulnerableFields    = getFields(EMPTY_USER, EMPTY_USER_IS_PASSWD_VULNERABLE)
)

func UpdateUserIsPasswdVulnerable(userIsPasswdVulnerable *UserIsPasswdVulnerable) (err error) {

	query := bson.M{
		"$or": bson.A{
			bson.M{
				USER_USER_ID_b: userIsPasswdVulnerable.UserID,
				USER_UPDATE_NANO_TS_b: bson.M{
					"$exists": false,
				},

				USER_IS_DELETED_b: bson.M{"$exists": false},
			},
			bson.M{
				USER_USER_ID_b: userIsPasswdVulnerable.UserID,
				USER_IS_PASSWD_VULNERABLE_UPDATE_NANO_TS_b: bson.M{
					"$lt": userIsPasswdVulnerable.IsPasswdVulnerableUpdateNanoTS,
				},

				USER_IS_DELETED_b: bson.M{"$exists": false},
			},
		},
	}

	r, err := User_c.UpdateOneOnly(query, userIsPasswdVulnerable)
	if err != nil {
		return err
	}
	if r.MatchedCount == 0 {
		return ErrNoMatch
	}
	return nil
}
