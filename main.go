package main

import (
	"awesomeProject/controllers"
	"awesomeProject/database"
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
	database.Migrate()

	router := initRouter()
	router.Run(":9000")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	assetsDir := os.Getenv("PHOTO_DIRECTORY")

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	router.Static("/images", assetsDir)
	router.LoadHTMLGlob("templates/*")
	api := router.Group("/api")
	{
		api.POST("/login", controllers.GenerateToken)
		api.POST("/forgot_password", controllers.ForgotPasword)
		api.GET("/reset_password", controllers.ResetPassword)
		api.GET("/delete_account", controllers.DeleteAccountForm)
		api.GET("/privacy_policy", controllers.PrivacyPolicy)
		api.POST("/register", controllers.RegisterUser)
		secured := api.Group("/mobile").Use(middlewares.Auth())
		{
			secured.POST("/upload", controllers.Upload)
			secured.GET("/results", controllers.GetResults)
			secured.POST("/user", controllers.UpdateUserDetails)
			secured.POST("/check_authentication_status", controllers.CheckAuthenticationStatus)
		}
	}
	return router
}
