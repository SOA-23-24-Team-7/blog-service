package controller

import (
	"BlogApplication/dto"
	"BlogApplication/model"
	"BlogApplication/service"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
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
	err = json.Unmarshal([]byte(decodedBody), &blog)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Your existing code for creating a new blog
	err = controller.BlogService.Create(&blog)
	if err != nil {
		fmt.Println("Error while creating a new blog:", err)
		writer.WriteHeader(http.StatusExpectationFailed)
		return
	}

	// Set the response status code and content type
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

func (controller *BlogController) Vote(writer http.ResponseWriter, req *http.Request) {

	var voteRequest dto.VoteRequest
	err := json.NewDecoder(req.Body).Decode(&voteRequest)
	if err != nil {
		http.Error(writer, "Error parsing request body", http.StatusBadRequest)
		return
	}

	blog, err := controller.BlogService.SetVote(voteRequest.BlogId, voteRequest.UserId, voteRequest.VoteType)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error upvoting blog: %v", err), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	err = json.NewEncoder(writer).Encode(blog)
	if err != nil {
		http.Error(writer, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
func (c *BlogController) FindAllWithType(w http.ResponseWriter, r *http.Request) {

	topicTypeStr, ok := mux.Vars(r)["type"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required parameter: id")
		return
	}
	topicType, err := model.ParseBlogTopicType(topicTypeStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid blog topic type: %v", err)
		return
	}
	println("TOPIC JE" + topicType)
	blogs, err := c.BlogService.GetBlogsByTopic(topicType)
	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error fetching blogs: %v", err)
		return
	}

	println(blogs)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(blogs)
	if err != nil {
		log.Printf("Error encoding blog data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error processing blog data")
		return
	}
}
