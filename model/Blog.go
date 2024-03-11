package model

import (
	"errors"
	"fmt"
	"time"
)

type BlogStatus string

const (
	Draft     BlogStatus = "draft"
	Published BlogStatus = "published"
	Closed    BlogStatus = "closed"
	Active    BlogStatus = "active"
	Famous    BlogStatus = "famous"
)

type BlogVisibilityPolicy string

const (
	PublicBlog  BlogVisibilityPolicy = "public"
	PrivateBlog BlogVisibilityPolicy = "private"
)

type Blog struct {
	Id          int        `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Date        time.Time  `json:"date"`
	Status      BlogStatus `json:"status"`
	AuthorId    int64      `json:"authorId"`
	//ClubId        *int64               `json:"clubId,omitempty"` // Optional club ID
	//Comments []Comment `json:"comments"`
	//Votes         []Vote               `json:"votes"`
	Visibility    BlogVisibilityPolicy `json:"visibility"`
	VoteCount     int64                `json:"voteCount"`
	UpvoteCount   int64                `json:"upvoteCount"`
	DownvoteCount int64                `json:"downvoteCount"`
}

func NewBlog(title string, description string, date time.Time, status BlogStatus, authorId int64, visibility BlogVisibilityPolicy) (*Blog, error) {
	if title == "" {
		return nil, fmt.Errorf("title cannot be empty or null")
	}

	blog := &Blog{
		Title:       title,
		Description: description,
		Date:        date,
		Status:      status,
		AuthorId:    authorId,
		Visibility:  visibility,
		//ClubId:      nil,
	}
	blog.calculateVoteCounts()
	return blog, nil
}

func (b *Blog) calculateVoteCounts() {
	//b.VoteCount = 0
	//b.UpvoteCount = 0
	//b.DownvoteCount = 0
	// for _, vote := range b.Votes {
	// 	if vote.VoteType == "UPVOTE" {
	// 		b.VoteCount++
	// 		b.UpvoteCount++
	// 	} else {
	// 		b.VoteCount--
	// 		b.DownvoteCount++
	// 	}
	// }
}

func (b *Blog) Validate() error {
	if b.Title == "" {
		return errors.New("Title can't be empty.")
	}

	if b.Description == "" {
		return errors.New("Description can't be empty.")
	}

	if b.Status == "" {
		return errors.New("Status can't be empty.")
	}

	if b.Visibility == "" {
		return errors.New("Visibility can't be empty.")
	}

	if b.AuthorId < 0 {
		return errors.New("AuthorId can't be less than 0.")
	}

	return nil
}
