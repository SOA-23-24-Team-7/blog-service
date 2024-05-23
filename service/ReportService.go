package service

import (
	"BlogApplication/model"
	"BlogApplication/repository"
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type ReportService struct {
	ReportRepository *repository.ReportRepository
}

func (service *ReportService) FindAllByBlog(ctx context.Context, id int64) ([]model.Report, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "FindAllByBlog")
	defer span.End()

	span.SetAttributes(attribute.String("request.data", "{ \"id\": "+string(id)+" }"))

	reports, _ := service.ReportRepository.FindAllByBlog(ctx, id)

	span.SetStatus(codes.Ok, "FindAllByBlog successful")
	return reports, nil
}

func (service *ReportService) Create(ctx context.Context, report *model.Report) error {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "Create")
	defer span.End()

	reqData, err := json.Marshal(report)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal request data")
		return err
	}
	span.SetAttributes(attribute.String("request.data", string(reqData)))

	err = service.ReportRepository.Create(ctx, report)
	if err != nil {
		span.SetStatus(codes.Error, "Create failed")
		return err
	}

	span.SetStatus(codes.Ok, "Create successful")
	return nil
}
