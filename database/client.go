package database

import (
	"awesomeProject/models"
	"fmt"
	"gorm.io/driver/postgres"
	"log"
	"os"

	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

func Connect(connectionString string) {
	Instance, dbError = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
	}
	log.Println("Connected to Database!")
}

func Migrate() {
	Instance.AutoMigrate(&models.User{})
	Instance.AutoMigrate(&models.Results{})
	Instance.AutoMigrate(&models.Clinic{})
	log.Println("Database Migration Completed!")
}

func GetPostgresConnectionString() string {

	environment := os.Getenv("ENVIRONMENT")

	if environment == "production" {
		return os.Getenv("PRODUCTION_DATABASE")
	}

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}
