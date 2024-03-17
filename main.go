package main

import (
	"BlogApplication/controller"
	"BlogApplication/model"
	"BlogApplication/repository"
	"BlogApplication/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {

	dsn := "user=postgres password=super dbname=soa-blog host=localhost port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		print(err)
		return nil
	}

	database.AutoMigrate(&model.Blog{})
	database.AutoMigrate(&model.Comment{})
	database.AutoMigrate(&model.Vote{})

	err = database.AutoMigrate(&model.Blog{}, &model.Comment{}, &model.Vote{})
	if err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}

	// newStudent := model.Student{
	// 	Person:     model.Person{Firstname: "John", Lastname: "Doe"},
	// 	Index:      "123456",
	// 	Major:      "Computer Science",
	// 	RandomData: model.RandomData{Years: 22},
	// }

	// Kada upisemo studenta, GORM ce automatski prvo da kreira Osobu i upise u
	// tabelu, a zatim Studenta, i to ce uraditi unutar iste transakcije.
	// database.Create(&newStudent)

	return database
}

func startServer(blogController *controller.BlogController, commentController *controller.CommentController) {
	router := mux.NewRouter().StrictSlash(true)

	// Blog routes
	router.HandleFunc("/blogs/type/{type}", blogController.FindAllWithType).Methods("GET")
	router.HandleFunc("/blogs", blogController.Create).Methods("POST")
	router.HandleFunc("/blogs/author/{id}", blogController.FindAllByAuthor).Methods("GET")
	router.HandleFunc("/blogs/published", blogController.FindAllPublished).Methods("GET")
	router.HandleFunc("/blogs/{id}", blogController.FindById).Methods("GET")
	router.HandleFunc("/blogs/{id}", blogController.Update).Methods("PUT")
	router.HandleFunc("/blogs/{id}", blogController.Delete).Methods("DELETE")

	// Blog vote route
	router.HandleFunc("/blogs/votes", blogController.Vote).Methods("POST")

	// Comment routes
	router.HandleFunc("/comments", commentController.Create).Methods("POST")
	router.HandleFunc("/comments/{id}", commentController.Update).Methods("PUT")
	router.HandleFunc("/comments/{id}", commentController.Delete).Methods("DELETE")
	router.HandleFunc("/comments", commentController.GetAll).Methods("GET")
	router.HandleFunc("/blogComments/{id}", commentController.GetAllBlogComments).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	println("Server starting")

	log.Fatal(http.ListenAndServe(":8090", router))

}

func main() {
	database := initDB()
	if database == nil {
		print("FAILED TO CONNECT TO DB")
		return
	}

	blogRepository := &repository.BlogRepository{DatabaseConnection: database}
	blogService := &service.BlogService{BlogRepository: blogRepository}
	blogController := &controller.BlogController{BlogService: blogService}

	commentRepository := &repository.CommentRepository{DatabaseConnection: database}
	commentService := &service.CommentService{CommentRepo: commentRepository}
	commentController := &controller.CommentController{CommentService: commentService}

	// voteRepository := &repository.VoteRepository{DatabaseConnection: database}
	// voteService := &service.VoteService{VoteRepo: voteRepository}
	// voteController := &controller.VoteController{VoteService: voteService}

	startServer(blogController, commentController)

}
