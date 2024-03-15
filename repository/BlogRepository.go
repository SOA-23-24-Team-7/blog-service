package repository

import (
	"BlogApplication/model"
	"fmt"

	"gorm.io/gorm"
)

type BlogRepository struct {
	DatabaseConnection *gorm.DB
}

func (repository *BlogRepository) Delete(id int64) error {
	dbResult := repository.DatabaseConnection.Delete(&model.Blog{}, id)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
func (repository *BlogRepository) UpdateVotes(blogID int, votes *[]model.Vote) error {
	err := repository.DatabaseConnection.Model(&model.Blog{}).
		Where("id = ?", blogID).
		Update("votes", votes).
		Error
	if err != nil {
		return fmt.Errorf("error updating votes: %w", err)
	}
	return nil
}

func (repository *BlogRepository) SetVote(b *model.Blog, userID int64, voteType model.VoteType) error {

	var existingVote *model.Vote
	for _, vote := range b.Votes {
		if vote.UserId == userID {
			existingVote = &vote
			break
		}
	}

	if existingVote != nil {

		if existingVote.VoteType != voteType {
			existingVote.VoteType = voteType

			err := repository.UpdateVotes(b.Id, &b.Votes)
			if err != nil {
				return fmt.Errorf("error updating votes: %w", err)
			}
			b.UpdateBlogStatus()
		}
		return nil
	}

	b.Votes = append(b.Votes, model.Vote{UserId: userID, VoteType: voteType})
	err := repository.UpdateVotes(b.Id, &b.Votes)
	if err != nil {
		return fmt.Errorf("error adding vote: %w", err)
	}
	b.UpdateBlogStatus()

	return nil
}

func (repository *BlogRepository) Find(id int64) (model.Blog, error) {
	blog := model.Blog{}
	dbResult := repository.DatabaseConnection. /*Preload("Comments").*/ Preload("Votes").First(&blog, id)
	println(blog.Title)
	if dbResult.Error != nil {
		return blog, dbResult.Error
	}
	return blog, nil
}

func (repository *BlogRepository) FindAllPublished() ([]model.Blog, error) {
	var blogs []model.Blog
	dbResult := repository.DatabaseConnection.Where("status = ?", "published").Find(&blogs)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return blogs, nil
}

func (repository *BlogRepository) FindAllByAuthor(id int64) ([]model.Blog, error) {
	var blogs []model.Blog
	dbResult := repository.DatabaseConnection.Where("author_id = ?", id).Find(&blogs)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return blogs, nil
}

func (repository *BlogRepository) Create(blog *model.Blog) error {
	dbResult := repository.DatabaseConnection.Create(blog)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repository *BlogRepository) Update(blog *model.Blog) error {
	dbResult := repository.DatabaseConnection.Save(blog)
	println(blog)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
