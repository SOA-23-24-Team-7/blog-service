package repository

import (
	"BlogApplication/model"

	"gorm.io/gorm"
)

type ReportRepository struct {
	DatabaseConnection *gorm.DB
}

func (repository *ReportRepository) FindAllByBlog(id int64) ([]model.Report, error) {
	var reports []model.Report
	dbResult := repository.DatabaseConnection.Where("blog_id = ?", id).Find(&reports)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return reports, nil
}

func (repository *ReportRepository) Create(report *model.Report) error {
	dbResult := repository.DatabaseConnection.Create(report)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
