package service

import (
	"BlogApplication/model"
	"BlogApplication/repository"
)

type ReportService struct {
	ReportRepository *repository.ReportRepository
}

func (service *ReportService) FindAllByBlog(id int64) (*[]model.Report, error) {
	reports, _ := service.ReportRepository.FindAllByBlog(id)
	return &reports, nil
}

func (service *ReportService) Create(report *model.Report) error {
	err := service.ReportRepository.Create(report)
	if err != nil {
		return err
	}
	return nil
}
