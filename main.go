package main

import (
	"awesomeProject/controllers"
	"awesomeProject/database"
	"awesomeProject/fcm"
	"awesomeProject/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	connectionString := database.GetPostgresConnectionString()

	// Initialize Database
	database.Connect(connectionString)
	//database.Migrate()
	fcm.InitializeFirebase()

	router := initRouter()
	err := router.Run(":9000")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func initRouter() *gin.Engine {
	router := gin.Default()
	assetsDir := os.Getenv("PHOTO_DIRECTORY")

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	router.Static("/images", assetsDir)
	router.LoadHTMLGlob("templates/*")
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/forgot_password", controllers.ForgotPassword)
		api.GET("/reset_password", controllers.ResetPassword)
		api.GET("/delete_account", controllers.DeleteAccountForm)
		api.GET("/privacy_policy", controllers.PrivacyPolicy)
		api.POST("/register", controllers.RegisterUser)
		secured := api.Group("/mobile").Use(middlewares.Auth())
		{
			secured.POST("/clinics", controllers.CreateClinic)
			secured.GET("/clinics", controllers.GetClinics)
			secured.POST("/results", controllers.UploadResults)
			secured.PUT("/results", controllers.UpdateResults)
			secured.GET("/results", controllers.GetResults)
			secured.POST("/user", controllers.UpdateUserDetails)
			secured.GET("/users", controllers.GetUserDetails)
			secured.POST("/check_authentication_status", controllers.CheckAuthenticationStatus)
		}
	}
	return router
}
