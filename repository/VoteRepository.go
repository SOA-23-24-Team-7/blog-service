package repository

import (
	"BlogApplication/model"

	"gorm.io/gorm"
)

type VoteRepository struct {
	DatabaseConnection *gorm.DB
}

func (repo *VoteRepository) FindById(id int) (model.Vote, error) {
	vote := model.Vote{}
	dbResult := repo.DatabaseConnection.First(&vote, id)
	if dbResult.Error != nil {
		return vote, dbResult.Error
	}
	return vote, nil
}

func (repo *VoteRepository) Create(vote *model.Vote) error {
	dbResult := repo.DatabaseConnection.Create(vote)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *VoteRepository) Update(vote *model.Vote) error {
	dbResult := repo.DatabaseConnection.Save(vote)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *VoteRepository) Delete(id int) error {
	dbResult := repo.DatabaseConnection.Delete(&model.Vote{}, id)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}

func (repo *VoteRepository) GetAll() ([]model.Vote, error) {
	var votes []model.Vote
	dbResult := repo.DatabaseConnection.Find(&votes)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return votes, nil
}
