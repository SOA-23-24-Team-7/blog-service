package server

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"BlogApplication/service"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type BlogMicroservice struct {
	UnimplementedBlogMicroserviceServer
	BlogService    *service.BlogService
	CommentService *service.CommentService
	ReportService  *service.ReportService
}

func (s *BlogMicroservice) FindBlogById(ctx context.Context, req *BlogIdRequest) (*BlogResponse, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "FindBlogById")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	blog, err := s.BlogService.Find(ctx, req.Id)

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

	span.SetStatus(codes.Ok, "FindBlogById successful")
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
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "FindBlogsByType")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	topicType, _ := model.ParseBlogTopicType(req.Type)
	blogsByTopic, err := s.BlogService.GetBlogsByTopic(ctx, topicType)

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

	span.SetStatus(codes.Ok, "FindBlogsByType successful")
	return &BlogListResponse{
		Blogs: blogs,
	}, err
}

func (s *BlogMicroservice) FindPublishedBlogs(ctx context.Context, req *Empty) (*BlogListResponse, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "FindPublishedBlogs")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	blogsByTopic, err := s.BlogService.FindAllPublished(ctx)

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

	span.SetStatus(codes.Ok, "FindPublishedBlogs successful")
	return &BlogListResponse{
		Blogs: blogs,
	}, err
}

func (s *BlogMicroservice) FindBlogsByAuthor(ctx context.Context, req *AuthorIdRequest) (*BlogListResponse, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "FindBlogsByAuthor")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	blogsByTopic, err := s.BlogService.FindAllByAuthor(ctx, req.Id)

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

	span.SetStatus(codes.Ok, "FindBlogsByAuthor successful")
	return &BlogListResponse{
		Blogs: blogs,
	}, err
}

func (s *BlogMicroservice) CreateBlog(ctx context.Context, req *BlogCreationRequest) (*StringMessage, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "CreateBlog")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	blog := &model.Blog{
		Title:       req.Title,
		Description: req.Description,
		AuthorId:    req.AuthorId,
		BlogTopic:   model.BlogTopicType(req.BlogTopic),
		Date:        time.Now(),
	}

	err = s.BlogService.Create(ctx, blog)

	if err != nil {
		fmt.Println("Error while creating a new blog:", err)
		message := &StringMessage{Message: "Error while creating blog"}
		span.SetStatus(codes.Error, "CreateBlog failed")
		return message, err
	}

	span.SetStatus(codes.Ok, "CreateBlog successful")
	message := &StringMessage{Message: "Successfully created blog"}
	return message, err
}

func (s *BlogMicroservice) DeleteBlog(ctx context.Context, req *BlogIdRequest) (*StringMessage, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "DeleteBlog")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	err = s.BlogService.Delete(ctx, req.Id)

	if err != nil {
		fmt.Println("Error while deleting a blog:", err)
		message := &StringMessage{Message: "Error while deleting a blog"}
		span.SetStatus(codes.Error, "DeleteBlog failed")
		return message, err
	}

	message := &StringMessage{Message: "Successfully deleted a blog"}
	span.SetStatus(codes.Ok, "DeleteBlog successful")
	return message, err
}

func (s *BlogMicroservice) BlockBlog(ctx context.Context, req *BlogIdRequest) (*StringMessage, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "BlockBlog")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	err = s.BlogService.Block(ctx, req.Id)

	if err != nil {
		fmt.Println("Error while blocking a blog:", err)
		message := &StringMessage{Message: "Error while blocking a blog"}
		span.SetStatus(codes.Error, "BlockBlog failed")
		return message, err
	}

	message := &StringMessage{Message: "Successfully blocked a blog"}
	span.SetStatus(codes.Ok, "BlockBlog successful")
	return message, err
}

func (s *BlogMicroservice) CreateComment(ctx context.Context, req *CommentCreationRequest) (*CommentResponse, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "CreateComment")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	comment := dto.CommentRequestDTO{
		AuthorId:  req.AuthorId,
		BlogId:    req.BlogId,
		CreatedAt: req.CreatedAt.AsTime(),
		Text:      req.Text,
	}
	createdComment, err := s.CommentService.Create(ctx, &comment)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		span.SetStatus(codes.Error, "CreateComment failed")
		return nil, err
	}

	span.SetStatus(codes.Ok, "CreateComment successful")
	return &CommentResponse{
		Id:        int32(createdComment.Id),
		AuthorId:  int64(createdComment.AuthorId),
		BlogId:    int64(createdComment.BlogId),
		CreatedAt: timestamppb.New(createdComment.CreatedAt),
		UpdatedAt: timestamppb.New(createdComment.UpdatedAt),
		Text:      createdComment.Text,
	}, err
}
func (s *BlogMicroservice) UpdateComment(ctx context.Context, req *CommentUpdateRequest) (*StringMessage, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "UpdateComment")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	comment := &dto.CommentUpdateDto{
		ID:   req.Id,
		Text: req.Text,
	}
	err = s.CommentService.Update(ctx, comment)
	if err != nil {
		log.Printf("Error updating comment: %v", err)
		span.SetStatus(codes.Error, "UpdateComment failed")
		return &StringMessage{Message: "Error updating comment"}, err
	}

	span.SetStatus(codes.Ok, "UpdateComment successful")
	return &StringMessage{Message: "Successfully updated comment"}, nil
}

func (s *BlogMicroservice) DeleteComment(ctx context.Context, req *CommentIdRequest) (*StringMessage, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "DeleteComment")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	err = s.CommentService.Delete(ctx, req.Id)
	if err != nil {
		log.Printf("Error deleting comment: %v", err)
		span.SetStatus(codes.Error, "DeleteComment failed")
		return &StringMessage{Message: "Error deleting comment"}, err
	}

	span.SetStatus(codes.Ok, "DeleteComment successful")
	return &StringMessage{Message: "Successfully deleted comment"}, nil
}

func (s *BlogMicroservice) GetAllComments(ctx context.Context, req *Empty) (*CommentListResponse, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "GetAllComments")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	comments, err := s.CommentService.GetAll(ctx)
	if err != nil {
		log.Printf("Error fetching all comments: %v", err)
		span.SetStatus(codes.Error, "GetAllComments failed")
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

	span.SetStatus(codes.Ok, "GetAllComments successful")
	return &CommentListResponse{Comments: response}, nil
}

func (s *BlogMicroservice) GetAllBlogComments(ctx context.Context, req *BlogIdRequest) (*CommentListResponse, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "GetAllBlogComments")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	comments, err := s.CommentService.GetAllBlogComments(ctx, req.Id)
	if err != nil {
		log.Printf("Error fetching comments for blog: %v", err)
		span.SetStatus(codes.Error, "GetAllBlogComments failed")
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

	span.SetStatus(codes.Ok, "GetAllBlogComments successful")
	return &CommentListResponse{Comments: response}, nil
}

func (s *BlogMicroservice) CreateReport(ctx context.Context, req *ReportRequest) (*StringMessage, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "CreateReport")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	report := &model.Report{
		BlogId: int(req.BlogId),
		UserId: int(req.UserId),
		Reason: req.Reason,
	}

	err = s.ReportService.Create(ctx, report)

	if err != nil {
		fmt.Println("Error while creating a new report:", err)
		message := &StringMessage{Message: "Error while creating report"}
		span.SetStatus(codes.Error, "CreateReport failed")
		return message, err
	}

	message := &StringMessage{Message: "Successfully created report"}
	span.SetStatus(codes.Ok, "CreateReport successful")
	return message, err
}

func (s *BlogMicroservice) FindReportsByBlog(ctx context.Context, req *BlogIdRequest) (*ReportListResponse, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "FindReportsByBlog")
	defer span.End()

	reqData, err := json.Marshal(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	reportsByBlog, err := s.ReportService.FindAllByBlog(ctx, req.Id)

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

	span.SetStatus(codes.Ok, "FindReportsByBlog successful")
	return &ReportListResponse{
		Reports: reports,
	}, err
}

func (s *BlogMicroservice) Vote(ctx context.Context, req *VoteRequest) (*StringMessage, error) {
	tracer := otel.Tracer("controller")
	ctx, span := tracer.Start(ctx, "Vote")
	defer span.End()

	reqData, jsonError := json.Marshal(req)
	if jsonError != nil {
		span.RecordError(jsonError)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, jsonError
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	_, err := s.BlogService.SetVote(ctx, req.BlogId, req.UserId, model.ParseVoteType(req.VoteType))

	message := &StringMessage{Message: "Successfully created voted"}
	span.SetStatus(codes.Ok, "Vote successful")
	return message, err
}
