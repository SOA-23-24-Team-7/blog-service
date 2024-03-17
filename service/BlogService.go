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

func (service *BlogService) FindAllPublished() ([]model.Blog, error) {
	blogs, _ := service.BlogRepository.FindAllPublished()
	return blogs, nil
}

func (service *BlogService) FindAllByAuthor(id int64) (*[]model.Blog, error) {
	blogs, _ := service.BlogRepository.FindAllByAuthor(id)
	return &blogs, nil
}

func (service *BlogService) Create(blog *model.Blog) error {
	blog.DownvoteCount = 0
	blog.UpvoteCount = 0
	blog.VoteCount = 0
	blog.Status = "published"
	blog.Visibility = "public"
	blog.Votes = []model.Vote{}
	blog.Comments = []model.Comment{}
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

func (service *BlogService) Block(id int64) error {
	oldBlog, err := service.BlogRepository.Find(id)
	if err != nil {

		return fmt.Errorf(fmt.Sprintf("menu item with id %d not found", id))
	}
	oldBlog.Visibility = "private"
	err = oldBlog.Validate()
	if err != nil {
		println("lol")
		return err
	}
	err = service.BlogRepository.Update(&oldBlog)
	if err != nil {
		println("ok")
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

func (service *BlogService) SetVote(blogID int64, userID int64, voteType model.VoteType) (*model.Blog, error) {
	blog, err := service.BlogRepository.Find(blogID)
	if err != nil {
		return nil, err
	}

	err = blog.SetVote(userID, voteType)
	if err != nil {
		return nil, err
	}

	err = service.BlogRepository.Update(&blog)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func (service *BlogService) GetBlogsByTopic(topicType model.BlogTopicType) ([]model.Blog, error) {
	var blogs []model.Blog
	err := service.BlogRepository.DatabaseConnection.Where("blog_topic = ?", topicType).
		Find(&blogs).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching blogs with topic type %s: %w", topicType, err)
	}

	return blogs, nil
}
