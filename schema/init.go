package schema

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/db"
)

func Init() (err error) {
	if client != nil {
		return nil
	}
	client, err = db.NewClient(MONGO_PROTOCOL, MONGO_HOST, MONGO_PORT, MONGO_DBNAME)
	if err != nil {
		return err
	}

	AccessToken_c = client.Collection("access_token")
	Article_c = client.Collection("article")
	Board_c = client.Collection("board")
	BoardBanuser_c = client.Collection("board_banuser")
	BoardChildren_c = client.Collection("board_children")
	BoardFriend_c = client.Collection("board_friend")
	Client_c = client.Collection("client")
	Comment_c = client.Collection("comment")

	User_c = client.Collection("user")
	UserAloha_c = client.Collection("user_aloha")
	UserFriend_c = client.Collection("user_friend")

	UserReadArticle_c = client.Collection("user_read_article")
	UserReadBoard_c = client.Collection("user_read_board")

	UserReject_c = client.Collection("user_reject")

	return nil
}

//Close
//
//XXX do not really close to avoid db connection-error in tests.
func Close() (err error) {
	return nil
	/*
		err = client.Close()
		if err != nil {
			log.Errorf("schema.Close: unable to close mongo: e: %v", err)
		}

		client = nil
		Client_c = nil
		User_c = nil
		AccessToken_c = nil
		UserReadArticle_c = nil
		UserReadBoard_c = nil

		return nil
	*/
}
