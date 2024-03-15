package dto

import "BlogApplication/model"

type VoteRequest struct {
	UserId   int64          `json:"userId"`
	BlogId   int64          `json:"blogId"`
	VoteType model.VoteType `json:"voteType"`
}
