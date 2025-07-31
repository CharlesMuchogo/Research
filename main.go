package main

import (
	"awesomeProject/database"
	"awesomeProject/fcm"
	"awesomeProject/routes"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Database
	database.Connect()
	database.Migrate()

	//Initialize firebase
	fcm.InitializeFirebase()

	router := routes.InitRouter()

	if err := router.Run(":9000"); err != nil {
		log.Fatal(err.Error())
	}
}
