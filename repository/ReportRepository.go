package repository

import (
	"BlogApplication/model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepository struct {
	Collection *mongo.Collection
}

func NewReportRepository(client *mongo.Client) *ReportRepository {
	database := client.Database("soa")
	collection := database.Collection("reports")
	return &ReportRepository{
		Collection: collection,
	}
}

func (repository *ReportRepository) FindAllByBlog(blogID int64) ([]model.Report, error) {
	var reports  = make([]model.Report, 0)
	filter := bson.M{"blogid": blogID}
	cur, err := repository.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var report model.Report
		err := cur.Decode(&report)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

func (repository *ReportRepository) Create(report *model.Report) error {
	report.Id = repository.NextId()
	_, err := repository.Collection.InsertOne(context.Background(), report)
	if err != nil {
		return err
	}
	return nil
}

func (repository *ReportRepository) GetAll() ([]model.Report, error) {
	var reports []model.Report
	cur, err := repository.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var report model.Report
		err := cur.Decode(&report)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

func (repository *ReportRepository) NextId() int {
	reports, _ := repository.GetAll()

	maxId := 0
	for _, report := range reports {
		if report.Id > maxId {
			maxId = report.Id
		}
	}

	return maxId + 1
}
