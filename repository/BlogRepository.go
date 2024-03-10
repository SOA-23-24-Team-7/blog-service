package repository

import (
	"BlogApplication/model"

	"gorm.io/gorm"
)

type BlogRepository struct {
	DatabaseConnection *gorm.DB
}

func (repository *BlogRepository) Find(id int64) (model.Blog, error) {
	blog := model.Blog{}
	dbResult := repository.DatabaseConnection.Preload("Comments").Preload("Votes").First(&blog, id)
	if dbResult.Error != nil {
		return blog, dbResult.Error
	}
	return blog, nil
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
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}

func (repository *BlogRepository) Delete(id int64) error {
	dbResult := repository.DatabaseConnection.Delete(&model.Blog{}, id)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	println("Rows affected: ", dbResult.RowsAffected)
	return nil
}
func (repository *BlogRepository) FindAll() ([]model.Blog, error) {
	var blogs []model.Blog
	dbResult := repository.DatabaseConnection.Preload("Comments").Preload("Votes").Find(&blogs)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return blogs, nil
}
