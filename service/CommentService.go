package service

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"BlogApplication/repository"
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

type CommentService struct {
	CommentRepo *repository.CommentRepository
}

func (service *CommentService) FindById(ctx context.Context, id int) (*model.Comment, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "FindById")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	comment, err := service.CommentRepo.FindById(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, "FindById failed")
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("comment with id %d not found", id)
		}
		return nil, err
	}

	span.SetStatus(codes.Ok, "FindById successful")
	return &comment, nil
}

func (service *CommentService) Create(ctx context.Context, commentRequest *dto.CommentRequestDTO) (*model.Comment, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()

	reqData, err := json.Marshal(commentRequest)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	comment := model.Comment{
		AuthorId:  commentRequest.AuthorId,
		BlogId:    commentRequest.BlogId,
		CreatedAt: commentRequest.CreatedAt,
		Text:      commentRequest.Text,
	}

	err = comment.Validate()
	if err != nil {
		span.SetStatus(codes.Error, "Create failed")
		return nil, fmt.Errorf("error validating comment: %w", err)
	}

	createdComment, err := service.CommentRepo.Create(ctx, &comment)
	if err != nil {
		span.SetStatus(codes.Error, "Create failed")
		return nil, fmt.Errorf("error creating comment: %w", err)
	}

	span.SetStatus(codes.Ok, "Create successful")
	return createdComment, nil
}

func (service *CommentService) Update(ctx context.Context, comment *dto.CommentUpdateDto) error {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()

	reqData, err := json.Marshal(comment)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	err = service.CommentRepo.Update(ctx, comment)
	if err != nil {
		span.SetStatus(codes.Error, "Update failed")
		return fmt.Errorf("error updating comment: %w", err)
	}

	span.SetStatus(codes.Ok, "Update successful")
	return nil
}

func (service *CommentService) Delete(ctx context.Context, id int64) error {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	err := service.CommentRepo.Delete(ctx, int64(id))
	if err != nil {
		span.SetStatus(codes.Error, "Delete failed")
		return fmt.Errorf("error deleting comment: %w", err)
	}

	span.SetStatus(codes.Ok, "Delete successful")
	return nil
}

func (service *CommentService) GetAll(ctx context.Context) ([]model.Comment, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "GetAll")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{}"))

	comments, err := service.CommentRepo.GetAll(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "GetAll failed")
		return nil, fmt.Errorf("error fetching all comments: %w", err)
	}

	span.SetStatus(codes.Ok, "GetAll successful")
	return comments, nil
}

func (service *CommentService) GetAllBlogComments(ctx context.Context, blogID int64) ([]model.Comment, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "GetAllBlogComments")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(blogID)+" }"))

	comments, err := service.CommentRepo.GetAllByBlog(ctx, int64(blogID))
	if err != nil {
		span.SetStatus(codes.Error, "GetAllBlogComments failed")
		return nil, fmt.Errorf("error fetching comments for blog ID %d: %w", blogID, err)
	}

	span.SetStatus(codes.Ok, "GetAllBlogComments successful")
	return comments, nil
}
