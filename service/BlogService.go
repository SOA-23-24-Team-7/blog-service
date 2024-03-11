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
	blog, err := service.BlogRepository.Find(id)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("menu item with id %d not found", id))
	}

	return &blog, nil
}

func (service *BlogService) FindAllPublished() (*[]model.Blog, error) {
	blogs, _ := service.BlogRepository.FindAllPublished()
	return &blogs, nil
}

func (service *BlogService) FindAllByAuthor(id int64) (*[]model.Blog, error) {
	blogs, _ := service.BlogRepository.FindAllByAuthor(id)
	return &blogs, nil
}

func (service *BlogService) Create(blog *model.Blog) error {
	blog.DownvoteCount = 0
	blog.UpvoteCount = 0
	blog.VoteCount = 0
	blog.Status = "draft"
	blog.Visibility = "public"
	err := blog.Validate()
	if err != nil {
		return err
	}
	err = service.BlogRepository.Create(blog)
	if err != nil {
		return err
	}
	return nil
}

func (service *BlogService) Update(id int64, blog *model.Blog) error {
	oldBlog, err := service.BlogRepository.Find(id)
	if err != nil {

		return fmt.Errorf(fmt.Sprintf("menu item with id %d not found", id))
	}
	oldBlog.Title = blog.Title
	oldBlog.Description = blog.Description
	err = oldBlog.Validate()
	if err != nil {
		return err
	}
	err = service.BlogRepository.Update(&oldBlog)
	if err != nil {
		return err
	}
	return nil
}

func (service *BlogService) Delete(id int64) error {
	err := service.BlogRepository.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting blog: %w", err)
	}
	return nil
}
