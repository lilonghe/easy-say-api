package viewModel

type MessageForm struct {
	Content string `json:"content" form:"Content" binding:"required"`
}

type LikeMessageForm struct {
	MessageId *string `form:"message_id" json:"message_id" binding:"required"`
	Unlike    *bool   `form:"unlike" json:"unlike" binding:"required"`
}

type MessageCommentForm struct {
	MessageId string `json:"message_id" form:"message_id" binding:"required"`
	Content   string `json:"content" form:"content" binding:"required"`
}

type DelMessageCommentForm struct {
	CommentId string `json:"comment_id" form:"comment_id" binding:"required"`
	//MessageId string `json:"message_id" binding:"required"`
}
