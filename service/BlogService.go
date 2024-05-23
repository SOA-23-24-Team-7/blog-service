package service

import (
	"BlogApplication/model"
	"BlogApplication/repository"
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type BlogService struct {
	BlogRepository *repository.BlogRepository
}

func (service *BlogService) Find(ctx context.Context, id int64) (*model.Blog, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Find")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	blog, err := service.BlogRepository.Find(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, "Find failed")
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %d not found", id))
	}

	span.SetStatus(codes.Ok, "Find successful")
	return &blog, nil
}

func (service *BlogService) FindAllPublished(ctx context.Context) ([]model.Blog, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "FindAllPublished")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{}"))

	blogs, _ := service.BlogRepository.FindAllPublished(ctx)

	span.SetStatus(codes.Ok, "FindAllPublished successful")
	return blogs, nil
}

func (service *BlogService) FindAllByAuthor(ctx context.Context, id int64) ([]model.Blog, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "FindAllByAuthor")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	blogs, _ := service.BlogRepository.FindAllByAuthor(ctx, id)

	span.SetStatus(codes.Ok, "FindAllByAuthor successful")
	return blogs, nil
}

func (service *BlogService) Create(ctx context.Context, blog *model.Blog) error {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()

	reqData, err := json.Marshal(blog)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	blog.DownvoteCount = 0
	blog.UpvoteCount = 0
	blog.VoteCount = 0
	blog.Status = "published"
	blog.Visibility = "public"
	blog.Votes = []model.Vote{}
	blog.Comments = []model.Comment{}
	err = blog.Validate()
	if err != nil {
		span.SetStatus(codes.Error, "Create failed")
		return err
	}
	err = service.BlogRepository.Create(ctx, blog)
	if err != nil {
		span.SetStatus(codes.Error, "Create failed")
		return err
	}
	span.SetStatus(codes.Ok, "Create successful")

	span.SetStatus(codes.Ok, "Create successful")
	return nil
}

func (service *BlogService) Update(ctx context.Context, id int64, blog *model.Blog) error {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()

	reqData, err := json.Marshal(blog)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	oldBlog, err := service.BlogRepository.Find(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, "Update failed")
		return fmt.Errorf(fmt.Sprintf("menu item with id %d not found", id))
	}
	oldBlog.Title = blog.Title
	oldBlog.Description = blog.Description
	err = oldBlog.Validate()
	if err != nil {
		span.SetStatus(codes.Error, "Update failed")
		return err
	}
	err = service.BlogRepository.Update(ctx, &oldBlog)
	if err != nil {
		span.SetStatus(codes.Error, "Update failed")
		return err
	}

	span.SetStatus(codes.Ok, "Update successful")
	return nil
}

func (service *BlogService) Block(ctx context.Context, id int64) error {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Block")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	oldBlog, err := service.BlogRepository.Find(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, "Block failed")
		return fmt.Errorf(fmt.Sprintf("menu item with id %d not found", id))
	}
	oldBlog.Visibility = "private"
	err = oldBlog.Validate()
	if err != nil {
		span.SetStatus(codes.Error, "Block failed")
		return err
	}
	err = service.BlogRepository.Update(ctx, &oldBlog)
	if err != nil {
		span.SetStatus(codes.Error, "Block failed")
		return err
	}

	span.SetStatus(codes.Ok, "Block successful")
	return nil
}

func (service *BlogService) Delete(ctx context.Context, id int64) error {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	err := service.BlogRepository.Delete(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, "Delete failed")
		return fmt.Errorf("error deleting blog: %w", err)
	}

	span.SetStatus(codes.Ok, "Delete successful")
	return nil
}

func (service *BlogService) SetVote(ctx context.Context, blogID int64, userID int64, voteType model.VoteType) (*model.Blog, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "SetVote")
	defer span.End()

	blog, err := service.BlogRepository.Find(ctx, blogID)
	if err != nil {
		span.SetStatus(codes.Error, "SetVote failed")
		return nil, err
	}

	err = blog.SetVote(userID, voteType)
	if err != nil {
		span.SetStatus(codes.Error, "SetVote failed")
		return nil, err
	}

	err = service.BlogRepository.Update(ctx, &blog)
	if err != nil {
		span.SetStatus(codes.Error, "SetVote failed")
		return nil, err
	}

	span.SetStatus(codes.Ok, "SetVote successful")
	return &blog, nil
}

func (service *BlogService) GetBlogsByTopic(ctx context.Context, topicType model.BlogTopicType) ([]model.Blog, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "GetBlogsByTopic")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"type\": "+string(topicType)+" }"))

	var blogs []model.Blog
	blogs, _ = service.BlogRepository.FindAllByTopic(ctx, topicType)

	span.SetStatus(codes.Ok, "GetBlogsByTopic successful")
	return blogs, nil
}
