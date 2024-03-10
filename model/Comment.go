package model

import (
	"time"
)

type Comment struct {
	Id        int        `json:"id"`
	AuthorId  int64      `json:"authorId"`
	BlogId    int64      `json:"blogId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Text      string     `json:"text"`
}

func NewComment(authorId, blogId int64, createdAt time.Time, updatedAt *time.Time, text string) (*Comment, error) {
	return &Comment{
		AuthorId:  authorId,
		BlogId:    blogId,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		Text:      text,
	}, nil
}
