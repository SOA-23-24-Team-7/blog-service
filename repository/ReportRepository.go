package repository

import (
	"BlogApplication/model"
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

func (repository *ReportRepository) FindAllByBlog(ctx context.Context, blogID int64) ([]model.Report, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "FindAllByBlog")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(blogID)+" }"))

	var reports = make([]model.Report, 0)
	filter := bson.M{"blogid": blogID}
	cur, err := repository.Collection.Find(context.Background(), filter)
	if err != nil {
		span.SetStatus(codes.Error, "FindAllByBlog failed")
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var report model.Report
		err := cur.Decode(&report)
		if err != nil {
			span.SetStatus(codes.Error, "FindAllByBlog failed")
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := cur.Err(); err != nil {
		span.SetStatus(codes.Error, "FindAllByBlog failed")
		return nil, err
	}

	span.SetStatus(codes.Ok, "FindAllByBlog successful")
	return reports, nil
}

func (repository *ReportRepository) Create(ctx context.Context, report *model.Report) error {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()

	reqData, jsonError := json.Marshal(report)
	if jsonError != nil {
		span.RecordError(jsonError)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return jsonError
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	report.Id = repository.NextId(ctx)
	_, err := repository.Collection.InsertOne(context.Background(), report)
	if err != nil {
		span.SetStatus(codes.Error, "Create failed")
		return err
	}

	span.SetStatus(codes.Ok, "Create successful")
	return nil
}

func (repository *ReportRepository) GetAll(ctx context.Context) ([]model.Report, error) {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "GetAll")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{}"))

	var reports []model.Report
	cur, err := repository.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		span.SetStatus(codes.Error, "GetAll failed")
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var report model.Report
		err := cur.Decode(&report)
		if err != nil {
			span.SetStatus(codes.Error, "GetAll failed")
			return nil, err
		}
		reports = append(reports, report)
	}
	if err := cur.Err(); err != nil {
		span.SetStatus(codes.Error, "GetAll failed")
		return nil, err
	}

	span.SetStatus(codes.Ok, "GetAll successful")
	return reports, nil
}

func (repository *ReportRepository) NextId(ctx context.Context) int {
	tracer := otel.Tracer("repository")
	ctx, span := tracer.Start(ctx, "NextId")
	defer span.End()

	reports, _ := repository.GetAll(ctx)

	maxId := 0
	for _, report := range reports {
		if report.Id > maxId {
			maxId = report.Id
		}
	}

	span.SetStatus(codes.Ok, "NextId successful")
	return maxId + 1
}
