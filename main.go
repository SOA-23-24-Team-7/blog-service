package main

import (
	//"BlogApplication/model"
	//"BlogApplication/repository"
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

	// database.AutoMigrate(&model.Person{})

	// err = database.AutoMigrate(&model.Person{}, &model.Student{})
	// if err != nil {
	// 	log.Fatalf("Error migrating models: %v", err)
	// }

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

func startServer() {
	router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc("/students/{id}", handler.Get).Methods("GET")
	// router.HandleFunc("/students", handler.Create).Methods("POST")

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
	// repo := &repo.StudentRepository{DatabaseConnection: database}
	// service := &service.StudentService{StudentRepo: repo}
	// handler := &handler.StudentHandler{StudentService: service}
	startServer()
}
