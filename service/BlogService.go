package service

import (
	"BlogApplication/model"
	"BlogApplication/repository"
	"fmt"
)

type BlogService struct {
	BlogRepository *repository.BlogRepository
}

func (service *BlogService) Find(id int64) (*model.Blog, error) {
	student, err := service.BlogRepository.Find(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %d not found", id))
	}

	return &student, nil
}

func (service *BlogService) Create(blog *model.Blog) error {
	err := service.BlogRepository.Create(blog)
	if err != nil {
		return err
	}
	return nil
}
