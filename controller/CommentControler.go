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

type CommentController struct {
	CommentService *service.CommentService
}

func (c *CommentController) FindById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required parameter: id")
		return
	}

	commentId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid comment Id format")
		return
	}

	comment, err := c.CommentService.FindById(commentId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Comment with Id %d not found", commentId)
			return
		}
		log.Printf("Error fetching comment: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func (c *CommentController) Create(w http.ResponseWriter, r *http.Request) {
	var comment model.Comment
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&comment)
	if err != nil {
		log.Printf("Error decoding comment request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid comment data format")
		return
	}

	err = c.CommentService.Create(&comment)
	if err != nil {
		log.Printf("Error creating comment: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error creating comment")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *CommentController) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required parameter: id")
		return
	}

	commentId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid comment Id format")
		return
	}

	var comment model.Comment
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&comment)
	if err != nil {
		log.Printf("Error decoding comment update request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid comment data format")
		return
	}

	comment.Id = commentId

	err = c.CommentService.Update(&comment)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Comment with Id %d not found", commentId)
			return
		}
		log.Printf("Error updating comment: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating comment")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CommentController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required parameter: id")
		return
	}

	commentId, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid comment Id format")
		return
	}

	err = c.CommentService.Delete(commentId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Comment with Id %d not found", commentId)
			return
		}
		log.Printf("Error deleting comment: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error deleting comment")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *CommentController) GetAll(w http.ResponseWriter, r *http.Request) {
	comments, err := c.CommentService.GetAll()
	if err != nil {
		log.Printf("Error fetching comments: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching comments")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
