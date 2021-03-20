package types

type CommentType uint8

const (
	COMMENT_TYPE_RECOMMEND CommentType = 1
	COMMENT_TYPE_BOO       CommentType = 2
	COMMENT_TYPE_COMMENT   CommentType = 3
	COMMENT_TYPE_REPLY     CommentType = 4
	COMMENT_TYPE_EDIT      CommentType = 5
	COMMENT_TYPE_FORWARD   CommentType = 6
	COMMENT_TYPE_DELETED   CommentType = 7
)
