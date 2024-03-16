package model

import (
	"errors"
	"fmt"
	"time"
)

type BlogTopicType string

const (
	BlogTopicTypeBiking      BlogTopicType = "biking"
	BlogTopicTypeFood        BlogTopicType = "food"
	BlogTopicTypeMuseums     BlogTopicType = "museums"
	BlogTopicTypeNature      BlogTopicType = "nature"
	BlogTopicTypeCulture     BlogTopicType = "culture"
	BlogTopicTypeHistory     BlogTopicType = "history"
	BlogTopicTypeBackpacking BlogTopicType = "backpacking"
	BlogTopicTypeSoloTravel  BlogTopicType = "soloTravel"
	BlogTopicTypeAdventure   BlogTopicType = "adventure"
	BlogTopicTypeArt         BlogTopicType = "art"
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
	Comments      []Comment            `json:"comments"`
	Votes         []Vote               `json:"votes" gorm:"foreignKey:BlogId"`
	Visibility    BlogVisibilityPolicy `json:"visibility"`
	VoteCount     int64                `json:"voteCount"`
	UpvoteCount   int64                `json:"upvoteCount"`
	DownvoteCount int64                `json:"downvoteCount"`
	BlogTopic     BlogTopicType        `json:"blogTopic"`
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
	b.VoteCount = 0
	b.UpvoteCount = 0
	b.DownvoteCount = 0
	for _, vote := range b.Votes {
		if vote.VoteType == "UPVOTE" {
			b.VoteCount++
			b.UpvoteCount++
		} else {
			b.VoteCount--
			b.DownvoteCount++
		}
	}
}

func (b *Blog) Validate() error {
	if b.Title == "" {
		return errors.New("title can't be empty")
	}

	if b.Description == "" {
		return errors.New("description can't be empty")
	}

	if b.Status == "" {
		return errors.New("status can't be empty")
	}

	if b.Visibility == "" {
		return errors.New("visibility can't be empty")
	}

	if _, err := ParseBlogTopicType(string(b.BlogTopic)); err != nil {
		return fmt.Errorf("invalid blog topic type: %w", err)
	}

	return nil
}

func (b *Blog) UpdateBlogStatus() {
	switch {
	case b.VoteCount < -2:
		b.Status = "closed"
	case b.VoteCount >= 3 && len(b.Comments) >= 3:
		b.Status = "famous"
	case b.VoteCount >= 2 && len(b.Comments) >= 2:
		b.Status = "active"
	default:
		b.Status = "published"
	}
}

func (b *Blog) SetVote(userID int64, voteType VoteType) error {
	for i, vote := range b.Votes {
		if vote.UserId == userID {
			if vote.VoteType == voteType {
				return nil
			}
			t := b.Votes[i]
			b.Votes[i] = Vote{Id: t.Id, UserId: userID, BlogId: int64(b.Id), VoteType: voteType} // Replace vote directly
			b.calculateVoteCounts()
			b.UpdateBlogStatus()
			return nil
		}
	}

	newVote := Vote{UserId: userID, BlogId: int64(b.Id), VoteType: voteType}
	b.Votes = append(b.Votes, newVote)
	b.calculateVoteCounts()
	b.UpdateBlogStatus()

	return nil
}
func ParseBlogTopicType(topicTypeStr string) (BlogTopicType, error) {
	switch topicTypeStr {
	case string(BlogTopicTypeBiking):
		return BlogTopicTypeBiking, nil
	case string(BlogTopicTypeFood):
		return BlogTopicTypeFood, nil
	case string(BlogTopicTypeMuseums):
		return BlogTopicTypeMuseums, nil
	case string(BlogTopicTypeNature):
		return BlogTopicTypeNature, nil
	case string(BlogTopicTypeCulture):
		return BlogTopicTypeCulture, nil
	case string(BlogTopicTypeHistory):
		return BlogTopicTypeHistory, nil
	case string(BlogTopicTypeBackpacking):
		return BlogTopicTypeBackpacking, nil
	case string(BlogTopicTypeSoloTravel):
		return BlogTopicTypeSoloTravel, nil
	case string(BlogTopicTypeAdventure):
		return BlogTopicTypeAdventure, nil
	case string(BlogTopicTypeArt):
		return BlogTopicTypeArt, nil
	default:
		return "", fmt.Errorf("invalid blog topic type: %s", topicTypeStr)
	}
}
