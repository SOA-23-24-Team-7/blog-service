package server

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"BlogApplication/service"
	"context"
	"fmt"
	"log"

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
func (s *BlogMicroservice) CreateComment(ctx context.Context, in *CommentCreationRequest) (*CommentResponse, error) {

	comment := dto.CommentRequestDTO{
		AuthorId:  in.AuthorId,
		BlogId:    in.BlogId,
		CreatedAt: in.CreatedAt.AsTime(),
		Text:      in.Text,
	}
	createdComment, err := s.CommentService.Create(&comment)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
	}

	return &CommentResponse{
		Id:        int32(createdComment.Id),
		AuthorId:  int64(createdComment.AuthorId),
		BlogId:    int64(createdComment.BlogId),
		CreatedAt: timestamppb.New(createdComment.CreatedAt),
		UpdatedAt: timestamppb.New(createdComment.UpdatedAt),
		Text:      createdComment.Text,
	}, err
}
func (s *BlogMicroservice) UpdateComment(ctx context.Context, in *CommentUpdateRequest) (*StringMessage, error) {
	comment := &dto.CommentUpdateDto{
		ID:   in.Id,
		Text: in.Text,
	}
	err := s.CommentService.Update(comment)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		return &StringMessage{Message: "Error updating comment"}, err
	}
	return &StringMessage{Message: "Successfully updated comment"}, nil
}

func (s *BlogMicroservice) DeleteComment(ctx context.Context, in *CommentIdRequest) (*StringMessage, error) {
	err := s.CommentService.Delete(in.Id)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		return &StringMessage{Message: "Error deleting comment"}, err
	}
	return &StringMessage{Message: "Successfully deleted comment"}, nil
}

func (s *BlogMicroservice) GetAllComments(ctx context.Context, in *Empty) (*CommentListResponse, error) {
	comments, err := s.CommentService.GetAll()
	if err != nil {
		log.Printf("Error fetching all comments: %v", err)
		return nil, err
	}

	var response []*CommentResponse
	for _, comment := range comments {
		response = append(response, &CommentResponse{
			Id:        int32(comment.Id),
			AuthorId:  comment.AuthorId,
			BlogId:    comment.BlogId,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			UpdatedAt: timestamppb.New(comment.UpdatedAt),
			Text:      comment.Text,
		})
	}
	return &CommentListResponse{Comments: response}, nil
}

func (s *BlogMicroservice) GetAllBlogComments(ctx context.Context, in *BlogIdRequest) (*CommentListResponse, error) {
	comments, err := s.CommentService.GetAllBlogComments(in.Id)
	if err != nil {
		log.Printf("Error fetching comments for blog: %v", err)
		return nil, err
	}

	var response []*CommentResponse
	for _, comment := range comments {
		response = append(response, &CommentResponse{
			Id:        int32(comment.Id),
			AuthorId:  comment.AuthorId,
			BlogId:    comment.BlogId,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			UpdatedAt: timestamppb.New(comment.UpdatedAt),
			Text:      comment.Text,
		})
	}
	return &CommentListResponse{Comments: response}, nil
}
