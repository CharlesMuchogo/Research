package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": record.Error.Error()})
		context.Abort()
		return
	}
	userToken, err := auth.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Failed to Log you in, please login with your credentials"})
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Signup success", "user": user, "token": userToken})
}
