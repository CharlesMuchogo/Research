package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/fcm"
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"sync"
)

func UpdateUserProfile(context *gin.Context) {

	var user models.User
	var existingUser models.User
	token := context.GetHeader("Authorization")

	var wg sync.WaitGroup
	wg.Add(1)

	userFromToken, err := auth.GetUserDetailsFromToken(token)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user details"})
		context.Abort()
		return
	}

	if err := database.DbInstance.Where("email = ?", userFromToken.Email).Find(&existingUser).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user details"})
		context.Abort()
		return
	}

	profilePhoto, _ := context.FormFile("image")

	randomUUID := uuid.New().String()

	userImageLink := existingUser.ProfilePhoto

	if profilePhoto != nil {
		userImageLink, err = utils.SavePhoto(profilePhoto, randomUUID)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating profile"})
			return
		}
	}

	testedBeforeStr := context.PostForm("tested_before")
	testedBefore := testedBeforeStr == "true"

	user.ID = existingUser.ID
	user.Email = existingUser.Email
	user.Phone = context.PostForm("phone")
	user.Country = context.PostForm("country")
	user.Password = existingUser.Password
	user.FirstName = context.PostForm("first_name")
	user.LastName = context.PostForm("last_name")
	user.Gender = context.PostForm("gender")
	user.Age = context.PostForm("age")
	user.Role = existingUser.Role
	user.EducationLevel = existingUser.EducationLevel
	user.SaveResults = existingUser.SaveResults
	user.TestedBefore = testedBefore
	user.ProfilePhoto = userImageLink
	user.DeviceId = context.PostForm("device_id")
	user.CreatedAt = existingUser.CreatedAt

	if err := database.DbInstance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user details"})
		context.Abort()
		return
	}

	go func() {
		defer wg.Done()
		fcm.RegisterTopic(user.Email, user.DeviceId)
	}()

	context.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": user})
	wg.Wait()
}

func GetUserProfile(context *gin.Context) {
	var user models.User

	token := context.GetHeader("Authorization")

	userFromToken, err := auth.GetUserDetailsFromToken(token)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user details"})
		context.Abort()
		return
	}

	if err := database.DbInstance.Where("email = ?", userFromToken.Email).Find(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user details"})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Profile retrieved successfully", "user": user})
}
