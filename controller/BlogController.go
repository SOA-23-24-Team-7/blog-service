package controller

import (
	"BlogApplication/model"
	"BlogApplication/service"
	"encoding/json"
	"net/http"
)

type BlogController struct {
	BlogService *service.BlogService
}

// func (controller *BlogController) Get(writer http.ResponseWriter, req *http.Request) {
// 	id := mux.Vars(req)["id"]
// 	log.Printf("Student sa id-em %s", id)
// 	student, err := handler.StudentService.FindStudent(id)
// 	writer.Header().Set("Content-Type", "application/json")
// 	if err != nil {
// 		writer.WriteHeader(http.StatusNotFound)
// 		return
// 	}
// 	writer.WriteHeader(http.StatusOK)
// 	json.NewEncoder(writer).Encode(student)
// }

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
