package database

import (
	"awesomeProject/models"
	"gorm.io/driver/postgres"
	"log"
	"os"

	"gorm.io/gorm"
)

var DbInstance *gorm.DB
var dbError error

func Connect() {
	connectionString := os.Getenv("PRODUCTION_DATABASE")
	DbInstance, dbError = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
	}
	log.Println("Connected to Database!")
}

func Migrate() {
	err := DbInstance.AutoMigrate(GetSchema()...)
	if err != nil {
		log.Printf("Database Migration Failed %s", err.Error())
	}
	log.Println("Database Migration Completed!")
}

func GetSchema() []any {
	return []any{
		&models.User{},
		&models.Results{},
		&models.Clinic{},
		&models.Message{},
		&models.ForgotPassword{},
	}
}
