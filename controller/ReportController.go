package controller

import (
	"BlogApplication/model"
	"BlogApplication/service"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ReportController struct {
	ReportService *service.ReportService
}

func (controller *ReportController) FindAllByBlog(writer http.ResponseWriter, req *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(req)["id"], 10, 64)
	reports, err := controller.ReportService.FindAllByBlog(id)
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(reports)
}

func (controller *ReportController) Create(writer http.ResponseWriter, req *http.Request) {
	var report model.Report

	// Read the request body
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	// Decode HTML entities in the request body
	decodedBody := html.UnescapeString(string(requestBody))

	// Print the decoded request body
	fmt.Println("Request Body:", decodedBody)

	// Decode the request body into the 'blog' struct
	err = json.Unmarshal([]byte(decodedBody), &report)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Your existing code for creating a new blog
	err = controller.ReportService.Create(&report)
	if err != nil {
		fmt.Println("Error while creating a new report:", err)
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	// Set the response status code and content type
	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
}
