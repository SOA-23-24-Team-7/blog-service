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

	err = database.AutoMigrate(&model.Blog{})
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

func startServer(blogController *controller.BlogController) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/blogs", blogController.Create).Methods("POST")

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

	startServer(blogController)
}
