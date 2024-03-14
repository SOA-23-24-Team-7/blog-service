package repository

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"fmt"

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

func (repo *CommentRepository) Create(comment *model.Comment) (*model.Comment, error) {
	dbResult := repo.DatabaseConnection.Create(comment)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	var createdComment model.Comment
	err := repo.DatabaseConnection.Where("id = ?", comment.Id).First(&createdComment).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching created comment: %w", err)
	}

	return &createdComment, nil
}

func (repo *CommentRepository) Update(commentUpdate *dto.CommentUpdateDto) error {

	var comment model.Comment
	err := repo.DatabaseConnection.Where("id = ?", commentUpdate.ID).First(&comment).Error
	if err != nil {
		return fmt.Errorf("error finding comment with ID %d: %w", commentUpdate.ID, err)
	}

	comment.Text = commentUpdate.Text

	dbResult := repo.DatabaseConnection.Save(&comment)
	if dbResult.Error != nil {
		return fmt.Errorf("error saving updated comment: %w", dbResult.Error)
	}

	return nil
}

func (repo *CommentRepository) Delete(id int) error {
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
