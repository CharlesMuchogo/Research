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

	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	token := context.GetHeader("Authorization")

	userFromToken, _ := auth.GetUserDetailsFromToken(token)

	user.ID = userFromToken.ID
	user.Email = userFromToken.Email
	user.Phone = userFromToken.Phone
	user.Password = userFromToken.Password
	user.FirstName = userFromToken.FirstName
	user.LastName = userFromToken.LastName
	user.ProfilePhoto = userFromToken.ProfilePhoto
	user.CreatedAt = time.Now()

	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update user details"})
		context.Abort()
		return
	}

	go fcm.RegisterTopic(user.Email, user.DeviceId)
	context.JSON(http.StatusOK, gin.H{"message": "Details updated successfully", "user": user})
}
