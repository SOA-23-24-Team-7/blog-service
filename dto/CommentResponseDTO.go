package dto

import "time"

type CommentResponseDto struct {
	ID        int        `json:"id"`
	AuthorID  int64      `json:"authorId"`
	BlogID    int64      `json:"blogId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Text      string     `json:"text"`
}
