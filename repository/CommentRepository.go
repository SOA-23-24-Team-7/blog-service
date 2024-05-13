package repository

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (repository *CommentRepository) FindById(id int) (model.Comment, error) {
	var comment model.Comment
	err := repository.Collection.FindOne(context.Background(), bson.M{"id": id}).Decode(&comment)
	if err != nil {
		return model.Comment{}, err
	}
	return comment, nil
}

func (repository *CommentRepository) Create(comment *model.Comment) (*model.Comment, error) {
	comment.Id = repository.NextId()
	_, err := repository.Collection.InsertOne(context.Background(), comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (repository *CommentRepository) Update(commentUpdate *dto.CommentUpdateDto) error {
	filter := bson.M{"id": commentUpdate.ID}
	update := bson.M{"$set": bson.M{"text": commentUpdate.Text}}
	_, err := repository.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (repository *CommentRepository) Delete(id int64) error {
	filter := bson.M{"id": id}
	_, err := repository.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (repository *CommentRepository) GetAll() ([]model.Comment, error) {
	var comments = make([]model.Comment, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var comment model.Comment
		err := cur.Decode(&comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (repository *CommentRepository) GetAllByBlog(id int64) ([]model.Comment, error) {
	var comments = make([]model.Comment, 0)
	cur, err := repository.Collection.Find(context.Background(), bson.M{"blogid": id})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var comment model.Comment
		err := cur.Decode(&comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repository *CommentRepository) NextId() int {
	comments, _ := repository.GetAll()

	maxId := 0
	for _, comment := range comments {
		if comment.Id > maxId {
			maxId = comment.Id
		}
	}

	return maxId + 1
}
