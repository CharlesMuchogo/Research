package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/fcm"
	"awesomeProject/models"
	"awesomeProject/models/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

func RegisterUser(context *gin.Context) {
	var userRequest dto.UserDTO

	var wg sync.WaitGroup
	wg.Add(1)

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

	record := database.DbInstance.Create(&user)

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

	context.JSON(http.StatusOK, gin.H{"message": "Signup success", "user": user, "token": userToken})

	go func() {
		defer wg.Done()
		go fcm.RegisterTopic(user.Email, user.DeviceId)
	}()
	wg.Wait()
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

	if err := database.DbInstance.Model(&models.User{}).
		Where("email = ?", userFromToken.Email).
		Select("TestedBefore", "Age", "EducationLevel", "Gender").
		Updates(user).Error; err != nil {
		fmt.Printf("Error updating user %s", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update details"})
		return
	}

	var updatedUser models.User
	if err := database.DbInstance.
		Where("email = ?", userFromToken.Email).
		First(&updatedUser).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch updated user"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Details updated successfully", "user": updatedUser})
}
