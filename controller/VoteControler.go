package controller

import (
	"BlogApplication/model"
	"BlogApplication/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type VoteController struct {
	VoteService *service.VoteService
}

func (c *VoteController) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required parameter: id")
		return
	}

	voteId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid vote Id format")
		return
	}

	vote, err := c.VoteService.FindById(voteId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Vote with Id %d not found", voteId)
			return
		}
		log.Printf("Error fetching vote: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vote)
}

func (c *VoteController) Create(w http.ResponseWriter, r *http.Request) {
	var vote model.Vote
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&vote)
	if err != nil {
		log.Printf("Error decoding vote request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid vote data format")
		return
	}

	err = c.VoteService.Create(&vote)
	if err != nil {
		log.Printf("Error creating vote: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating vote")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *VoteController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required parameter: id")
		return
	}

	voteId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid vote Id format")
		return
	}

	var vote model.Vote
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&vote)
	if err != nil {
		log.Printf("Error decoding vote update request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid vote data format")
		return
	}

	vote.Id = voteId

	err = c.VoteService.Update(&vote)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Vote with Id %d not found", voteId)
			return
		}
		log.Printf("Error updating vote: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating vote")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *VoteController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required parameter: id")
		return
	}

	voteId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid vote Id format")
		return
	}

	err = c.VoteService.Delete(voteId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Vote with Id %d not found", voteId)
			return
		}
		log.Printf("Error deleting vote: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting vote")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *VoteController) GetAll(w http.ResponseWriter, r *http.Request) {
	votes, err := c.VoteService.GetAll()
	if err != nil {
		log.Printf("Error fetching votes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching votes")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votes)
}
