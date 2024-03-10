package repository

import (
	"BlogApplication/model"

	"gorm.io/gorm"
)

type CommentRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *CommentRepository) FindById(id int) (model.Comment, error) {
	comment := model.Comment{}
	dbResult := repo.DatabaseConnection.First(&comment, id)
	if dbResult.Error != nil {
		return comment, dbResult.Error
	}
	return comment, nil
}

func (repo *CommentRepository) CreateComment(comment *model.Comment) error {
	dbResult := repo.DatabaseConnection.Create(comment)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *CommentRepository) UpdateComment(comment *model.Comment) error {
	dbResult := repo.DatabaseConnection.Save(comment)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *CommentRepository) DeleteComment(id int) error {
	dbResult := repo.DatabaseConnection.Delete(&model.Comment{}, id)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *CommentRepository) GetAll() ([]model.Comment, error) {
	var comments []model.Comment
	dbResult := repo.DatabaseConnection.Find(&comments)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return comments, nil
}
