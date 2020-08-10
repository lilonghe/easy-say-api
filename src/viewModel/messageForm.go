package viewModel

type MessageForm struct {
	Content string `json:"content"`
}

type LikeMessageForm struct {
	MessageId string `json:"message_id"`
	Unlike    bool   `json:"unlike"`
}
