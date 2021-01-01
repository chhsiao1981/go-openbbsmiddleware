package schema

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/db"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
)

var (
	User_c *db.Collection
)

type User struct {
	UserID   bbs.UUserID `bson:"user_id"`
	Username string      `bson:"username"`
	Realname string      `bson:"realtime"`
	Nickname string      `bson:"nickname"`
	Avatar   []byte      `bson:"avatar"`

	UFlag        ptttype.UFlag `bson:"flag"`
	UserLevel    ptttype.PERM  `bson:"perm"`
	NumLoginDays int           `bson:"login_days"` /* 考慮透過 db-count, 但是可能無法跟以前版本相容 */
	NumPosts     int           `bson:"posts"`      /* 考慮透過 db-count, 因為 post 可能會被消失, 但是這個欄位不應該因為被消失的而減少? (poster 主動消失的需要減少)? */
	FirstLogin   types.NanoTS  `bson:"first_login_ts"`
	LastLogin    types.NanoTS  `bson:"last_login_ts"` /* 考慮透過 db-max, 但是可能在拉 user-detail 時會花很多時間. */
	LastIP       string        `bson:"last_ip"`
	LastHost     string        `bson:"last_host"` //last-ip 的中文呈現, 外國則為國家.

	Money            int    `bson:"money"`
	Email            string `bson:"email"`
	EmailVerified    bool   `bson:"email_verified"`
	TWIDHash         string `bson:"twidhash"`
	TWIDHashVerified bool   `bson:"twidhash_verified"`

	TwoFactorEnabled bool `bson:"twofactor_enabled"`

	Justify string `bson:"justify"`
	Over18  bool   `bson:"over18"`

	PagerUIType int               `bson:"pager_ui"` /* 呼叫器界面類別 (was: WATER_*) */
	Pager       ptttype.PagerMode `bson:"pager"`    /* 呼叫器狀態 */
	Invisible   bool              `bson:"invisible"`
	ExMailbox   int               `bson:"exmail"`

	Career string `bson:"career"`
	Role   int    `bson:"role"`

	LastSeen      types.NanoTS `bson:"last_seen_nano_ts"`
	TimeSetAngel  types.NanoTS `bson:"time_set_angel_nano_ts"`
	TimePlayAngel types.NanoTS `bson:"time_play_angel_nano_ts"`

	LastSong  types.NanoTS `bson:"last_song_nano_ts"`
	LoginView int          `bson:"login_view"`

	VlCount   int   `bson:"violation"`
	FiveWin   int   `bson:"five_win"`
	FiveLose  int   `bson:"five_lose"`
	FiveTie   int   `bson:"five_tie"`
	ChcWin    int   `bson:"chc_win"`
	ChcLose   int   `bson:"chc_lose"`
	ChcTie    int   `bson:"chc_tie"`
	Conn6Win  int   `bson:"conn6_win"`
	Conn6Lose int   `bson:"conn6_lose"`
	Conn6Tie  int   `bson:"conn6_tie"`
	GoWin     int   `bson:"go_win"`
	GoLose    int   `bson:"go_lose"`
	GoTie     int   `bson:"go_tie"`
	DarkWin   int   `bson:"dark_win"`
	DarkLose  int   `bson:"dark_lose"`
	DarkTie   int   `bson:"dark_tie"` /* 暗棋戰績 和 */
	UaVersion uint8 `bson:"ua_version"`

	Signature uint8  `bson:"signaure"` /* 慣用簽名檔 */
	BadPost   int    `bson:"bad_post"` /* 評價為壞文章數 */
	MyAngel   string `bson:"angel"`    /* 我的小天使 */

	ChessEloRating int `bson:"chess_rank"` /* 象棋等級 */

	TimeRemoveBadPost types.NanoTS `bson:"time_remove_bad_post_nano_ts"`
	TimeViolateLaw    types.NanoTS `bson:"time_violate_law_nano_ts"`

	IsDeleted bool `bson:"deleted,omitempty"`

	IsPasswdVulnerable             bool         `bson:"pw_vulnerable"`
	IsPasswdVulnerableUpdateNanoTS types.NanoTS `bson:"pw_vulnerable_nano_ts"`

	UpdateNanoTS types.NanoTS `bson:"update_nano_ts"`

	//NFriend int `bson:"n_friend"` /* 需要透過 db-count */
}

var (
	EMPTY_USER = &User{}
)

var (
	USER_USER_ID_b  = getBSONName(EMPTY_USER, "UserID")
	USER_USERNAME_b = getBSONName(EMPTY_USER, "Username")
	USER_REALNAME_b = getBSONName(EMPTY_USER, "Realname")
	USER_NICKNAME_b = getBSONName(EMPTY_USER, "Nickname")
	USER_AVATAR_b   = getBSONName(EMPTY_USER, "Avatar")

	USER_UFLAG_b          = getBSONName(EMPTY_USER, "UFlag")
	USER_USER_LEVEL_b     = getBSONName(EMPTY_USER, "UserLevel")
	USER_NUM_LOGIN_DAYS_b = getBSONName(EMPTY_USER, "NumLoginDays")
	USER_NUM_POSTS_b      = getBSONName(EMPTY_USER, "NumPosts")
	USER_FIRST_LOGIN_b    = getBSONName(EMPTY_USER, "FirstLogin")
	USER_LAST_LOGIN_b     = getBSONName(EMPTY_USER, "LastLogin")
	USER_LAST_IP_b        = getBSONName(EMPTY_USER, "LastIP")
	USER_LAST_HOST_b      = getBSONName(EMPTY_USER, "LastHost")

	USER_MONEY_b              = getBSONName(EMPTY_USER, "Money")
	USER_EMAIL_b              = getBSONName(EMPTY_USER, "Email")
	USER_EMAIL_VERIFIED_b     = getBSONName(EMPTY_USER, "EmailVerified")
	USER_TWID_HASH_b          = getBSONName(EMPTY_USER, "TWIDHash")
	USER_TWID_HASH_VERIFIED_b = getBSONName(EMPTY_USER, "TWIDHashVerified")

	USER_TWO_FACTOR_ENABLED_b = getBSONName(EMPTY_USER, "TwoFactorEnabled")

	USER_JUSTIFY_b = getBSONName(EMPTY_USER, "Justify")
	USER_OVER18_b  = getBSONName(EMPTY_USER, "Over18")

	USER_PAGER_UI_TYPE_b = getBSONName(EMPTY_USER, "PagerUIType")
	USER_PAGER_b         = getBSONName(EMPTY_USER, "Pager")
	USER_INVISIBLE_b     = getBSONName(EMPTY_USER, "Invisible")
	USER_EX_MAILBOX_b    = getBSONName(EMPTY_USER, "ExMailbox")

	USER_CAREER_b          = getBSONName(EMPTY_USER, "Career")
	USER_ROLE_b            = getBSONName(EMPTY_USER, "Role")
	USER_LAST_SEEN_b       = getBSONName(EMPTY_USER, "LastSeen")
	USER_TIME_SET_ANGEL_b  = getBSONName(EMPTY_USER, "TimeSetAngel")
	USER_TIME_PLAY_ANGEL_b = getBSONName(EMPTY_USER, "TimePlayAngel")
	USER_LAST_SONG_b       = getBSONName(EMPTY_USER, "LastSong")
	USER_LOGIN_VIEW_b      = getBSONName(EMPTY_USER, "LoginView")

	USER_VL_COUNT_b   = getBSONName(EMPTY_USER, "VlCount")
	USER_FIVE_WIN_b   = getBSONName(EMPTY_USER, "FiveWin")
	USER_FIVE_LOSE_b  = getBSONName(EMPTY_USER, "FiveLose")
	USER_FIVE_TIE_b   = getBSONName(EMPTY_USER, "FiveTie")
	USER_CHC_WIN_b    = getBSONName(EMPTY_USER, "ChcWin")
	USER_CHC_LOSE_b   = getBSONName(EMPTY_USER, "ChcLose")
	USER_CHC_TIE_b    = getBSONName(EMPTY_USER, "ChcTie")
	USER_CONN6_WIN_b  = getBSONName(EMPTY_USER, "Conn6Win")
	USER_CONN6_LOSE_b = getBSONName(EMPTY_USER, "Conn6Lose")
	USER_CONN6_TIE_b  = getBSONName(EMPTY_USER, "Conn6Tie")
	USER_GO_WIN_b     = getBSONName(EMPTY_USER, "GoWin")
	USER_GO_LOSE_b    = getBSONName(EMPTY_USER, "GoLose")
	USER_GO_TIE_b     = getBSONName(EMPTY_USER, "GoTie")
	USER_DARK_WIN_b   = getBSONName(EMPTY_USER, "DarkWin")
	USER_DARK_LOSE_b  = getBSONName(EMPTY_USER, "DarkLose")
	USER_DARK_TIE_b   = getBSONName(EMPTY_USER, "DarkTie")

	USER_UA_VERSION_b = getBSONName(EMPTY_USER, "UaVersion")

	USER_SIGNATURE_b        = getBSONName(EMPTY_USER, "Signature")
	USER_BAD_POST_b         = getBSONName(EMPTY_USER, "BadPost")
	USER_MY_ANGEL_b         = getBSONName(EMPTY_USER, "MyAngel")
	USER_CHESS_ELO_RATING_b = getBSONName(EMPTY_USER, "ChessEloRating")

	USER_TIME_REMOVE_BAD_POST_b = getBSONName(EMPTY_USER, "TimeRemoveBadPost")
	USER_TIME_VIOLATE_LAW_b     = getBSONName(EMPTY_USER, "TimeViolateLaw")

	USER_IS_DELETED_b = getBSONName(EMPTY_USER, "IsDeleted")

	USER_IS_PASSWD_VULNERABLE_b                = getBSONName(EMPTY_USER, "IsPasswdVulnerable")
	USER_IS_PASSWD_VULNERABLE_UPDATE_NANO_TS_b = getBSONName(EMPTY_USER, "IsPasswdVulnerableUpdateNanoTS")

	USER_UPDATE_NANO_TS_b = getBSONName(EMPTY_USER, "UpdateNanoTS")
)

func assertUserFields() error {
	if err := assertFields(EMPTY_USER, EMPTY_USER_QUERY); err != nil {
		return err
	}

	if err := assertFields(EMPTY_USER, EMPTY_USER_SUMMARY); err != nil {
		return err
	}

	if err := assertFields(EMPTY_USER, EMPTY_USER_IS_PASSWD_VULNERABLE); err != nil {
		return err
	}

	return nil
}

type UserQuery struct {
	UserID bbs.UUserID `bson:"user_id"`
}

var (
	EMPTY_USER_QUERY = &UserQuery{}
)
