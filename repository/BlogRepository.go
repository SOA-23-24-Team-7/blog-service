package repository

import (
	"BlogApplication/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (repository *BlogRepository) Delete(id int64) error {
	_, err := repository.Collection.DeleteOne(context.Background(), bson.M{"id": id})
	if err != nil {
		return err
	}
	return nil
}
func (repository *BlogRepository) UpdateVotes(blogID int, votes *[]model.Vote) error {
	filter := bson.M{"id": blogID}
	update := bson.M{"$set": bson.M{"votes": votes}}
	_, err := repository.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repository *BlogRepository) SetVote(b *model.Blog, userID int64, voteType model.VoteType) error {

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

			err := repository.UpdateVotes(b.Id, &b.Votes)
			if err != nil {
				return fmt.Errorf("error updating votes: %w", err)
			}
			b.UpdateBlogStatus()
		}
		return nil
	}

	b.Votes = append(b.Votes, model.Vote{UserId: userID, VoteType: voteType})
	err := repository.UpdateVotes(b.Id, &b.Votes)
	if err != nil {
		return fmt.Errorf("error adding vote: %w", err)
	}
	b.UpdateBlogStatus()

	return nil
}

func (repository *BlogRepository) Find(id int64) (model.Blog, error) {
	var blog model.Blog
	err := repository.Collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&blog)
	if err != nil {
		return model.Blog{}, err
	}
	return blog, nil
}

func (repository *BlogRepository) FindAllPublished() ([]model.Blog, error) {
	var blogs = make([]model.Blog, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var blog model.Blog
		err := cur.Decode(&blog)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (repository *BlogRepository) FindAllByAuthor(id int64) ([]model.Blog, error) {
	var blogs = make([]model.Blog, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{"authorid": id})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var blog model.Blog
		err := cur.Decode(&blog)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (repository *BlogRepository) FindAllByTopic(topicType model.BlogTopicType) ([]model.Blog, error) {
	var blogs = make([]model.Blog, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{"blogtopic": topicType})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var blog model.Blog
		err := cur.Decode(&blog)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	return blogs, nil
}

func (repository *BlogRepository) Create(blog *model.Blog) error {
	blog.Id = repository.NextId()
	_, err := repository.Collection.InsertOne(context.Background(), blog)
	if err != nil {
		return err
	}
	return nil
}

func (repository *BlogRepository) Update(blog *model.Blog) error {
	filter := bson.M{"id": blog.Id}
	update := bson.M{"$set": blog}
	_, err := repository.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repository *BlogRepository) NextId() int {
	blogs, _ := repository.FindAllPublished()

	maxId := 0
	for _, blog := range blogs {
		if blog.Id > maxId {
			maxId = blog.Id
		}
	}

	return maxId + 1
}
