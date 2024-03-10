package repository

import (
	"BlogApplication/model"

	"gorm.io/gorm"
)

type BlogRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *BlogRepository) FindById(id int64) (model.Blog, error) {
	blog := model.Blog{}
	dbResult := repo.DatabaseConnection.Preload("Comments").Preload("Votes").First(&blog, id)
	if dbResult.Error != nil {
		return blog, dbResult.Error
	}
	return blog, nil
}

func (repo *BlogRepository) CreateBlog(blog *model.Blog) error {
	dbResult := repo.DatabaseConnection.Create(blog)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *BlogRepository) UpdateBlog(blog *model.Blog) error {
	dbResult := repo.DatabaseConnection.Save(blog)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repo *BlogRepository) DeleteBlog(id int64) error {
	dbResult := repo.DatabaseConnection.Delete(&model.Blog{}, id)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
func (repo *BlogRepository) GetAll() ([]model.Blog, error) {
	var blogs []model.Blog
	dbResult := repo.DatabaseConnection.Preload("Comments").Preload("Votes").Find(&blogs)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return blogs, nil
}
