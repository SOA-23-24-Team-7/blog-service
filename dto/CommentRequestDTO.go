package dto

import "time"

type CommentRequestDTO struct {
	AuthorId  int64     `json:"authorId"`
	BlogId    int64     `json:"blogId"`
	CreatedAt time.Time `json:"createdAt"`
	Text      string    `json:"text"`
}
