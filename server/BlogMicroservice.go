package server

import (
	"BlogApplication/model"
	"BlogApplication/service"
	"context"
	"fmt"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type BlogMicroservice struct {
	UnimplementedBlogMicroserviceServer
	BlogService    *service.BlogService
	CommentService *service.CommentService
	ReportService  *service.ReportService
}

func (s *BlogMicroservice) FindBlogById(ctx context.Context, req *BlogIdRequest) (*BlogResponse, error) {
	blog, err := s.BlogService.Find(req.Id)

	var comments []*CommentResponse
	for _, c := range blog.Comments {
		comments = append(comments, &CommentResponse{
			Id:        int32(c.Id),
			AuthorId:  c.AuthorId,
			BlogId:    c.BlogId,
			CreatedAt: timestamppb.New(c.CreatedAt),
			UpdatedAt: timestamppb.New(c.UpdatedAt),
			Text:      c.Text,
		})
	}
	if comments == nil {
		comments = []*CommentResponse{}
	}

	var votes []*VoteResponse
	for _, v := range blog.Votes {
		votes = append(votes, &VoteResponse{
			Id:       int32(v.Id),
			UserId:   v.UserId,
			BlogId:   v.BlogId,
			VoteType: string(v.VoteType),
		})
	}
	if votes == nil {
		votes = []*VoteResponse{}
	}

	return &BlogResponse{
		Id:            int32(blog.Id),
		Title:         blog.Title,
		Description:   blog.Description,
		Date:          timestamppb.New(blog.Date),
		Status:        string(blog.Status),
		AuthorId:      blog.AuthorId,
		Comments:      comments,
		Votes:         votes,
		Visibility:    string(blog.Visibility),
		VoteCount:     blog.VoteCount,
		UpvoteCount:   blog.UpvoteCount,
		DownvoteCount: blog.DownvoteCount,
		BlogTopic:     string(blog.BlogTopic),
	}, err
}

func (s *BlogMicroservice) CreateBlog(ctx context.Context, in *BlogCreationRequest) (*StringMessage, error) {
	blog := &model.Blog{
		Title:       in.Title,
		Description: in.Description,
		AuthorId:    in.AuthorId,
		BlogTopic:   model.BlogTopicType(in.BlogTopic),
	}

	err := s.BlogService.Create(blog)

	if err != nil {
		fmt.Println("Error while creating a new blog:", err)
		message := &StringMessage{Message: "Error while creating blog"}
		return message, err
	}

	message := &StringMessage{Message: "Successfully created blog"}
	return message, err
}
