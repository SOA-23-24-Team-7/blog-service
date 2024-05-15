package server

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"BlogApplication/service"
	"context"
	"fmt"
	"log"
	"time"

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

func (s *BlogMicroservice) FindBlogsByType(ctx context.Context, req *TypeRequest) (*BlogListResponse, error) {
	topicType, _ := model.ParseBlogTopicType(req.Type)
	blogsByTopic, err := s.BlogService.GetBlogsByTopic(topicType)

	var blogs []*BlogResponse
	for _, b := range blogsByTopic {

		var comments []*CommentResponse
		for _, c := range b.Comments {
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
		for _, v := range b.Votes {
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

		blogs = append(blogs, &BlogResponse{
			Id:            int32(b.Id),
			Title:         b.Title,
			Description:   b.Description,
			Date:          timestamppb.New(b.Date),
			Status:        string(b.Status),
			AuthorId:      b.AuthorId,
			Comments:      comments,
			Votes:         votes,
			Visibility:    string(b.Visibility),
			VoteCount:     b.VoteCount,
			UpvoteCount:   b.UpvoteCount,
			DownvoteCount: b.DownvoteCount,
			BlogTopic:     string(b.BlogTopic),
		})
	}
	if blogs == nil {
		blogs = []*BlogResponse{}
	}

	return &BlogListResponse{
		Blogs: blogs,
	}, err
}

func (s *BlogMicroservice) FindPublishedBlogs(ctx context.Context, req *Empty) (*BlogListResponse, error) {
	blogsByTopic, err := s.BlogService.FindAllPublished()

	var blogs []*BlogResponse
	for _, b := range blogsByTopic {

		var comments []*CommentResponse
		for _, c := range b.Comments {
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
		for _, v := range b.Votes {
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

		blogs = append(blogs, &BlogResponse{
			Id:            int32(b.Id),
			Title:         b.Title,
			Description:   b.Description,
			Date:          timestamppb.New(b.Date),
			Status:        string(b.Status),
			AuthorId:      b.AuthorId,
			Comments:      comments,
			Votes:         votes,
			Visibility:    string(b.Visibility),
			VoteCount:     b.VoteCount,
			UpvoteCount:   b.UpvoteCount,
			DownvoteCount: b.DownvoteCount,
			BlogTopic:     string(b.BlogTopic),
		})
	}
	if blogs == nil {
		blogs = []*BlogResponse{}
	}

	return &BlogListResponse{
		Blogs: blogs,
	}, err
}

func (s *BlogMicroservice) FindBlogsByAuthor(ctx context.Context, req *AuthorIdRequest) (*BlogListResponse, error) {
	blogsByTopic, err := s.BlogService.FindAllByAuthor(req.Id)

	var blogs []*BlogResponse
	for _, b := range blogsByTopic {

		var comments []*CommentResponse
		for _, c := range b.Comments {
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
		for _, v := range b.Votes {
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

		blogs = append(blogs, &BlogResponse{
			Id:            int32(b.Id),
			Title:         b.Title,
			Description:   b.Description,
			Date:          timestamppb.New(b.Date),
			Status:        string(b.Status),
			AuthorId:      b.AuthorId,
			Comments:      comments,
			Votes:         votes,
			Visibility:    string(b.Visibility),
			VoteCount:     b.VoteCount,
			UpvoteCount:   b.UpvoteCount,
			DownvoteCount: b.DownvoteCount,
			BlogTopic:     string(b.BlogTopic),
		})
	}
	if blogs == nil {
		blogs = []*BlogResponse{}
	}

	return &BlogListResponse{
		Blogs: blogs,
	}, err
}

func (s *BlogMicroservice) CreateBlog(ctx context.Context, in *BlogCreationRequest) (*StringMessage, error) {
	blog := &model.Blog{
		Title:       in.Title,
		Description: in.Description,
		AuthorId:    in.AuthorId,
		BlogTopic:   model.BlogTopicType(in.BlogTopic),
		Date: time.Now(),
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

func (s *BlogMicroservice) DeleteBlog(ctx context.Context, in *BlogIdRequest) (*StringMessage, error) {

	err := s.BlogService.Delete(in.Id)

	if err != nil {
		fmt.Println("Error while deleting a blog:", err)
		message := &StringMessage{Message: "Error while deleting a blog"}
		return message, err
	}

	message := &StringMessage{Message: "Successfully deleted a blog"}
	return message, err
}

func (s *BlogMicroservice) BlockBlog(ctx context.Context, in *BlogIdRequest) (*StringMessage, error) {

	err := s.BlogService.Block(in.Id)

	if err != nil {
		fmt.Println("Error while blocking a blog:", err)
		message := &StringMessage{Message: "Error while blocking a blog"}
		return message, err
	}

	message := &StringMessage{Message: "Successfully blocked a blog"}
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

func (s *BlogMicroservice) CreateReport(ctx context.Context, in *ReportRequest) (*StringMessage, error) {
	report := &model.Report{
		BlogId: int(in.BlogId),
		UserId: int(in.UserId),
		Reason: in.Reason,
	}

	err := s.ReportService.Create(report)

	if err != nil {
		fmt.Println("Error while creating a new report:", err)
		message := &StringMessage{Message: "Error while creating report"}
		return message, err
	}

	message := &StringMessage{Message: "Successfully created report"}
	return message, err
}

func (s *BlogMicroservice) FindReportsByBlog(ctx context.Context, req *BlogIdRequest) (*ReportListResponse, error) {
	reportsByBlog, err := s.ReportService.FindAllByBlog(req.Id)

	var reports []*ReportResponse
	for _, r := range reportsByBlog {

		reports = append(reports, &ReportResponse{
			Id:     int64(r.Id),
			UserId: int64(r.UserId),
			BlogId: int64(r.BlogId),
			Reason: r.Reason,
		})
	}
	if reports == nil {
		reports = []*ReportResponse{}
	}

	return &ReportListResponse{
		Reports: reports,
	}, err
}

func (s *BlogMicroservice) Vote(ctx context.Context, in *VoteRequest) (*StringMessage, error) {
	_, err := s.BlogService.SetVote(in.BlogId, in.UserId, model.ParseVoteType(in.VoteType))
	message := &StringMessage{Message: "Successfully created voted"}
	return message, err
}
