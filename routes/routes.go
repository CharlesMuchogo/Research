package routes

import (
	"awesomeProject/controllers"
	"awesomeProject/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func InitRouter() *gin.Engine {
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
	router.NoRoute(controllers.NotFound)
	router.GET("/delete_account", controllers.DeleteAccountForm)
	router.GET("/privacy_policy", controllers.PrivacyPolicy)
	router.GET("/terms_and_conditions", controllers.TermsAndConditions)

	//Users
	api := router.Group("/api")
	{
		api.POST("/login", controllers.Login)
		api.POST("/google-login", controllers.GoogleLogin)
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
