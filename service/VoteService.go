package service

import (
	"BlogApplication/model"
	"BlogApplication/repository"
	"fmt"

	"gorm.io/gorm"
)

type VoteService struct {
	VoteRepo *repository.VoteRepository
}

func (service *VoteService) FindVoteById(id int) (*model.Vote, error) {
	vote, err := service.VoteRepo.FindById(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("vote with id %d not found", id)
		}
		return nil, err
	}
	return &vote, nil
}

func (service *VoteService) Create(vote *model.Vote) error {
	err := service.VoteRepo.Create(vote)
	if err != nil {
		return fmt.Errorf("error creating vote: %w", err)
	}
	return nil
}

func (service *VoteService) Update(vote *model.Vote) error {
	err := service.VoteRepo.Update(vote)
	if err != nil {
		return fmt.Errorf("error updating vote: %w", err)
	}
	return nil
}

func (service *VoteService) Delete(id int) error {
	err := service.VoteRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting vote: %w", err)
	}
	return nil
}

func (service *VoteService) GetAllVotes() ([]model.Vote, error) {
	votes, err := service.VoteRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error fetching all votes: %w", err)
	}
	return votes, nil
}