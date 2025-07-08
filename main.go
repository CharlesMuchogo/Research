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
	database.Migrate()
	fcm.InitializeFirebase()

	router := initRouter()

	//Set up AWS
	err := router.Run(":9000")
	if err != nil {
		log.Fatal(err.Error())
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
	router.GET("/delete_account", controllers.DeleteAccountForm)
	router.GET("/privacy_policy", controllers.PrivacyPolicy)
	router.GET("/terms_and_conditions", controllers.TermsAndConditions)

	//Users
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/forgot_password", controllers.ForgotPassword)
		api.GET("/reset_password", controllers.ResetPassword)
		api.POST("/update_password", controllers.UpdatePassword)
		api.POST("/register", controllers.RegisterUser)

		users := api.Group("/mobile").Use(middlewares.Auth())
		{
			users.GET("/clinics", controllers.GetClinics)
			users.POST("/results", controllers.UploadResults)
			users.POST("/messages", controllers.SendMessages)
			users.GET("/messages", controllers.GetMessages)
			users.DELETE("/results", controllers.DeleteResults)
			users.GET("/results", controllers.GetResults)
			users.POST("/user", controllers.UpdateUserDetails)
			users.GET("/profile", controllers.GetUserProfile)
			users.PUT("/profile", controllers.UpdateUserProfile)
		}

	}

	//Admin
	adminApi := router.Group("/admin")
	{
		adminApi.POST("/login", controllers.AdminLogin)
		admin := adminApi.Group("/api").Use(middlewares.AdminOnly())
		{
			admin.PUT("/results", controllers.UpdateResults)
			admin.POST("/clinics", controllers.CreateClinic)
			admin.GET("/results", controllers.GetAllResults)
			admin.GET("/users", controllers.GetUserDetails)
		}
	}
	return router
}
