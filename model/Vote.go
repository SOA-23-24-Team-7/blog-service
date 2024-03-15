package model

import "fmt"

type VoteType string

const (
	Downvote VoteType = "DOWNVOTE"
	Upvote   VoteType = "UPVOTE"
)

type Vote struct {
	Id       int      `json:"id"`
	UserId   int64    `json:"userId"`
	BlogId   int64    `json:"blogId"`
	VoteType VoteType `json:"voteType"`
}

func NewVote(userId int64, voteType VoteType) *Vote {
	return &Vote{
		UserId:   userId,
		VoteType: voteType,
	}
}

func (v *Vote) Validate() error {
	if v.UserId <= 0 {
		return fmt.Errorf("user ID must be a positive integer")
	}

	if v.VoteType != Downvote && v.VoteType != Upvote {
		return fmt.Errorf("invalid vote type: %s, allowed values are 'DOWNVOTE' or 'UPVOTE'", v.VoteType)
	}

	return nil
}
