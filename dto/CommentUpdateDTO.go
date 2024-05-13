package dto

type CommentUpdateDto struct {
	ID   int64  `json:"id"`
	Text string `json:"text"`
}
