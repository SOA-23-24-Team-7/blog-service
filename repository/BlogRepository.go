package repository

import (
	"BlogApplication/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type BlogRepository struct {
	Collection *mongo.Collection
}

func NewBlogRepository(client *mongo.Client) *BlogRepository {
	database := client.Database("soa")
	collection := database.Collection("blogs")
	return &BlogRepository{
		Collection: collection,
	}
}

func (repository *BlogRepository) Delete(ctx context.Context, id int64) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Delete")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	_, err := repository.Collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		span.SetStatus(codes.Error, "Delete failed")
		return err
	}

	span.SetStatus(codes.Ok, "Delete successful")
	return nil
}
func (repository *BlogRepository) UpdateVotes(ctx context.Context, blogID int, votes *[]model.Vote) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "UpdateVotes")
	defer span.End()

	reqData, jsonError := json.Marshal(votes)
	if jsonError != nil {
		span.RecordError(jsonError)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return jsonError
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	filter := bson.M{"id": blogID}
	update := bson.M{"$set": bson.M{"votes": votes}}
	_, err := repository.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		span.SetStatus(codes.Error, "UpdateVotes failed")
		return err
	}

	span.SetStatus(codes.Ok, "UpdateVotes successful")
	return nil
}

func (repository *BlogRepository) SetVote(ctx context.Context, b *model.Blog, userID int64, voteType model.VoteType) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "SetVote")
	defer span.End()

	var existingVote *model.Vote
	for _, vote := range b.Votes {
		if vote.UserId == userID {
			existingVote = &vote
			break
		}
	}

	if existingVote != nil {

		if existingVote.VoteType != voteType {
			existingVote.VoteType = voteType

			err := repository.UpdateVotes(ctx, b.Id, &b.Votes)
			if err != nil {
				span.SetStatus(codes.Error, "SetVote failed")
				return fmt.Errorf("error updating votes: %w", err)
			}
			b.UpdateBlogStatus()
		}
		return nil
	}

	b.Votes = append(b.Votes, model.Vote{UserId: userID, VoteType: voteType})
	err := repository.UpdateVotes(ctx, b.Id, &b.Votes)
	if err != nil {
		span.SetStatus(codes.Error, "SetVote failed")
		return fmt.Errorf("error adding vote: %w", err)
	}
	b.UpdateBlogStatus()

	span.SetStatus(codes.Ok, "SetVote successful")
	return nil
}

func (repository *BlogRepository) Find(ctx context.Context, id int64) (model.Blog, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Find")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	var blog model.Blog
	err := repository.Collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&blog)
	if err != nil {
		span.SetStatus(codes.Error, "Find failed")
		return model.Blog{}, err
	}

	span.SetStatus(codes.Ok, "Find successful")
	return blog, nil
}

func (repository *BlogRepository) FindAllPublished(ctx context.Context) ([]model.Blog, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "FindAllPublished")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{}"))

	var blogs = make([]model.Blog, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		span.SetStatus(codes.Error, "FindAllPublished failed")
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var blog model.Blog
		err := cur.Decode(&blog)
		if err != nil {
			span.SetStatus(codes.Error, "FindAllPublished failed")
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	span.SetStatus(codes.Ok, "FindAllPublished successful")
	return blogs, nil
}

func (repository *BlogRepository) FindAllByAuthor(ctx context.Context, id int64) ([]model.Blog, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "FindAllByAuthor")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	var blogs = make([]model.Blog, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{"authorid": id})
	if err != nil {
		span.SetStatus(codes.Error, "FindAllByAuthor failed")
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var blog model.Blog
		err := cur.Decode(&blog)
		if err != nil {
			span.SetStatus(codes.Error, "FindAllByAuthor failed")
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	span.SetStatus(codes.Ok, "FindAllByAuthor successful")
	return blogs, nil
}

func (repository *BlogRepository) FindAllByTopic(ctx context.Context, topicType model.BlogTopicType) ([]model.Blog, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "FindAllByTopic")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"topic\": "+string(topicType)+" }"))

	var blogs = make([]model.Blog, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{"blogtopic": topicType})
	if err != nil {
		span.SetStatus(codes.Error, "FindAllByTopic failed")
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var blog model.Blog
		err := cur.Decode(&blog)
		if err != nil {
			span.SetStatus(codes.Error, "FindAllByTopic failed")
			return nil, err
		}
		blogs = append(blogs, blog)
	}

	span.SetStatus(codes.Ok, "FindAllByTopic successful")
	return blogs, nil
}

func (repository *BlogRepository) Create(ctx context.Context, blog *model.Blog) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()

	reqData, jsonError := json.Marshal(blog)
	if jsonError != nil {
		span.RecordError(jsonError)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return jsonError
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	blog.Id = repository.NextId(ctx)
	blog.Date = time.Now()
	_, err := repository.Collection.InsertOne(context.Background(), blog)
	if err != nil {
		span.SetStatus(codes.Error, "Create failed")
		return err
	}
	span.SetStatus(codes.Ok, "Create successful")
	return nil
}

func (repository *BlogRepository) Update(ctx context.Context, blog *model.Blog) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Update")
	defer span.End()

	reqData, jsonError := json.Marshal(blog)
	if jsonError != nil {
		span.RecordError(jsonError)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return jsonError
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	filter := bson.M{"id": blog.Id}
	update := bson.M{"$set": blog}
	_, err := repository.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		span.SetStatus(codes.Error, "Update failed")
		return err
	}

	span.SetStatus(codes.Ok, "Update successful")
	return nil
}

func (repository *BlogRepository) NextId(ctx context.Context) int {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "NextId")
	defer span.End()

	blogs, _ := repository.FindAllPublished(ctx)

	maxId := 0
	for _, blog := range blogs {
		if blog.Id > maxId {
			maxId = blog.Id
		}
	}

	span.SetStatus(codes.Ok, "NextId successful")
	return maxId + 1
}
