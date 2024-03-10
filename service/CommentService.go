package service

import (
	"BlogApplication/model"
	"BlogApplication/repository"
	"fmt"

	"gorm.io/gorm"
)

type CommentService struct {
	CommentRepo *repository.CommentRepository
}

func (service *CommentService) FindCommentById(id int) (*model.Comment, error) {
	comment, err := service.CommentRepo.FindById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("comment with id %d not found", id)
		}
		return nil, err
	}
	return &comment, nil
}

func (service *CommentService) Create(comment *model.Comment) error {
	err := service.CommentRepo.Create(comment)
	if err != nil {
		return fmt.Errorf("error creating comment: %w", err)
	}
	return nil
}

func (service *CommentService) Update(comment *model.Comment) error {
	err := service.CommentRepo.Update(comment)
	if err != nil {
		return fmt.Errorf("error updating comment: %w", err)
	}
	return nil
}

func (service *CommentService) Delete(id int) error {
	err := service.CommentRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}
	return nil
}

func (service *CommentService) GetAllComments() ([]model.Comment, error) {
	comments, err := service.CommentRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error fetching all comments: %w", err)
	}
	return comments, nil
}
