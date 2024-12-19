package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserDetails(context *gin.Context) {
	var users []models.User

	tokenString := context.GetHeader("Authorization")
	user, _ := auth.GetUserDetailsFromToken(tokenString)

	if err := database.Instance.Where("email != ?", user.Email).Order("id DESC").Find(&users).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Users fetched successfully", "users": users})
}
