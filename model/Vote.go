package model

type VoteType string

const (
	Downvote VoteType = "DOWNVOTE"
	Upvote   VoteType = "UPVOTE"
)

type Vote struct {
	Id       int      `json:"id"`
	UserId   int64    `json:"userId"`
	VoteType VoteType `json:"voteType"`
}

func NewVote(userId int64, voteType VoteType) *Vote {
	return &Vote{
		UserId:   userId,
		VoteType: voteType,
	}
}
