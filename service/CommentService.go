package service

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"BlogApplication/repository"
	"fmt"

	"gorm.io/gorm"
)

type CommentService struct {
	CommentRepo *repository.CommentRepository
}

func (service *CommentService) FindById(id int) (*model.Comment, error) {
	comment, err := service.CommentRepo.FindById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("comment with id %d not found", id)
		}
		return nil, err
	}
	return &comment, nil
}

func (service *CommentService) Create(commentRequest *dto.CommentRequestDTO) (*model.Comment, error) {
	comment := model.Comment{
		AuthorId:  commentRequest.AuthorId,
		BlogId:    commentRequest.BlogId,
		CreatedAt: commentRequest.CreatedAt,
		Text:      commentRequest.Text,
	}

	err := comment.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating comment: %w", err)
	}

	createdComment, err := service.CommentRepo.Create(&comment)
	if err != nil {
		return nil, fmt.Errorf("error creating comment: %w", err)
	}

	return createdComment, nil
}

func (service *CommentService) Update(comment *dto.CommentUpdateDto) error {

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

func (service *CommentService) GetAll() ([]model.Comment, error) {
	comments, err := service.CommentRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error fetching all comments: %w", err)
	}
	return comments, nil
}

func (service *CommentService) GetAllBlogComments(blogID int) ([]model.Comment, error) {
	var comments []model.Comment
	err := service.CommentRepo.DatabaseConnection.Where("blog_id = ?", blogID).Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching comments for blog ID %d: %w", blogID, err)
	}
	return comments, nil
}
