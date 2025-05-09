package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/fcm"
	"awesomeProject/models"
	"awesomeProject/models/dto"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var userRequest dto.UserDTO

	if err := context.ShouldBindJSON(&userRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	user := userRequest.ToUser()

	if err := user.HashPassword(userRequest.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	user.CreatedAt = time.Now()
	user.DeviceId = userRequest.DeviceId

	record := database.Instance.Create(&user)

	if record.Error != nil {
		if strings.Contains(record.Error.Error(), "users_email_key") {
			context.JSON(http.StatusBadRequest, gin.H{"message": "An account with this email exists"})
		} else if strings.Contains(record.Error.Error(), "users_phone") {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "An account with this phone number exists"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"message": record.Error.Error()})
		}
		context.Abort()
		return
	}

	userToken, err := auth.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Signup success!, please login with your credentials"})
		context.Abort()
		return
	}
	go fcm.RegisterTopic(user.Email, user.DeviceId)
	context.JSON(http.StatusOK, gin.H{"message": "Signup success", "user": user, "token": userToken})
}

func UpdateUserDetails(context *gin.Context) {
	var user models.User

	var existingUser models.User

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	if err := database.Instance.Where("email = ?", user.Email).Find(&existingUser).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update details"})
		context.Abort()
		return
	}

	if existingUser.Email == "" {
		context.JSON(http.StatusNotFound, gin.H{"message": "Profile not found"})
		context.Abort()
		return
	}

	user.ID = existingUser.ID
	user.Email = existingUser.Email
	user.Phone = existingUser.Phone
	user.Password = existingUser.Password
	user.FirstName = existingUser.FirstName
	user.LastName = existingUser.LastName
	user.ProfilePhoto = existingUser.ProfilePhoto
	user.Country = existingUser.Country
	user.CreatedAt = existingUser.CreatedAt

	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update user details"})
		context.Abort()
		return
	}

	go fcm.RegisterTopic(user.Email, user.DeviceId)
	context.JSON(http.StatusOK, gin.H{"message": "Details updated successfully", "user": user})
}
