package repository

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type CommentRepository struct {
	Collection *mongo.Collection
}

func NewCommentRepository(client *mongo.Client) *CommentRepository {
	database := client.Database("soa")
	collection := database.Collection("comments")
	return &CommentRepository{
		Collection: collection,
	}
}

func (repository *CommentRepository) FindById(ctx context.Context, id int) (model.Comment, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "FindById")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	var comment model.Comment
	err := repository.Collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&comment)
	if err != nil {
		span.SetStatus(codes.Error, "FindById failed")
		return model.Comment{}, err
	}

	span.SetStatus(codes.Ok, "FindById successful")
	return comment, nil
}

func (repository *CommentRepository) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()

	reqData, jsonError := json.Marshal(comment)
	if jsonError != nil {
		span.RecordError(jsonError)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return nil, jsonError
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	comment.Id = repository.NextId(ctx)
	_, err := repository.Collection.InsertOne(context.Background(), comment)
	if err != nil {
		span.SetStatus(codes.Error, "Create failed")
		return nil, err
	}

	span.SetStatus(codes.Ok, "Create successful")
	return comment, nil
}

func (repository *CommentRepository) Update(ctx context.Context, commentUpdate *dto.CommentUpdateDto) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()

	reqData, jsonError := json.Marshal(commentUpdate)
	if jsonError != nil {
		span.RecordError(jsonError)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return jsonError
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	filter := bson.M{"id": commentUpdate.ID}
	update := bson.M{"$set": bson.M{"text": commentUpdate.Text}}
	_, err := repository.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		span.SetStatus(codes.Error, "Update failed")
		return err
	}

	span.SetStatus(codes.Ok, "Update successful")
	return nil
}

func (repository *CommentRepository) Delete(ctx context.Context, id int64) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	filter := bson.M{"id": id}
	_, err := repository.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		span.SetStatus(codes.Error, "Delete failed")
		return err
	}

	span.SetStatus(codes.Ok, "Delete successful")
	return nil
}

func (repository *CommentRepository) GetAll(ctx context.Context) ([]model.Comment, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "GetAll")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{}"))

	var comments = make([]model.Comment, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		span.SetStatus(codes.Error, "GetAll failed")
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var comment model.Comment
		err := cur.Decode(&comment)
		if err != nil {
			span.SetStatus(codes.Error, "GetAll failed")
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := cur.Err(); err != nil {
		span.SetStatus(codes.Error, "GetAll failed")
		return nil, err
	}

	span.SetStatus(codes.Ok, "GetAll successful")
	return comments, nil
}

func (repository *CommentRepository) GetAllByBlog(ctx context.Context, id int64) ([]model.Comment, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "GetAllByBlog")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	var comments = make([]model.Comment, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{"blogid": id})
	if err != nil {
		span.SetStatus(codes.Error, "GetAllByBlog failed")
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var comment model.Comment
		err := cur.Decode(&comment)
		if err != nil {
			span.SetStatus(codes.Error, "GetAllByBlog failed")
			return nil, err
		}
		comments = append(comments, comment)
	}

	span.SetStatus(codes.Ok, "GetAllByBlog successful")
	return comments, nil
}

func (repository *CommentRepository) NextId(ctx context.Context) int {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "NextId")
	defer span.End()

	comments, _ := repository.GetAll(ctx)

	maxId := 0
	for _, comment := range comments {
		if comment.Id > maxId {
			maxId = comment.Id
		}
	}

	span.SetStatus(codes.Ok, "NextId successful")
	return maxId + 1
}
