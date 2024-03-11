package controller

import (
	"BlogApplication/model"
	"BlogApplication/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BlogController struct {
	BlogService *service.BlogService
}

func (controller *BlogController) FindById(writer http.ResponseWriter, req *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	blog, err := controller.BlogService.Find(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blog)
}

func (controller *BlogController) FindAllPublished(writer http.ResponseWriter, req *http.Request) {
	blogs, err := controller.BlogService.FindAllPublished()
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blogs)
}

func (controller *BlogController) FindAllByAuthor(writer http.ResponseWriter, req *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	blogs, err := controller.BlogService.FindAllByAuthor(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(blogs)
}

func (controller *BlogController) Create(writer http.ResponseWriter, req *http.Request) {
	var blog model.Blog
	err := json.NewDecoder(req.Body).Decode(&blog)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = controller.BlogService.Create(&blog)
	if err != nil {
		println("Error while creating a new blog")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}

func (controller *BlogController) Update(writer http.ResponseWriter, req *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	var blog model.Blog
	err := json.NewDecoder(req.Body).Decode(&blog)
	if err != nil {
		println("Error while parsing json")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	err = controller.BlogService.Update(id, &blog)
	if err != nil {
		println("Error while creating a new blog")
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}

func (controller *BlogController) Delete(writer http.ResponseWriter, req *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)

	err := controller.BlogService.Delete(id)
	if err != nil {
		fmt.Fprintf(writer, "Error deleting comment")
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
