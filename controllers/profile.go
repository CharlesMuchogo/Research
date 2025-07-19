package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func UpdateUserProfile(context *gin.Context) {

	token := context.GetHeader("Authorization")

	userFromToken, err := auth.GetUserDetailsFromToken(token)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user details"})
		context.Abort()
		return
	}

	profilePhoto, _ := context.FormFile("image")

	randomUUID := uuid.New().String()

	var userImageLink string

	if profilePhoto != nil {
		userImageLink, err = utils.SavePhoto(profilePhoto, randomUUID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating profile"})
			return
		}
	}

	testedBeforeStr := context.PostForm("tested_before")
	testedBefore := testedBeforeStr == "true"

	if err := database.DbInstance.Model(&models.User{}).
		Where("email = ?", userFromToken.Email).
		Updates(models.User{
			TestedBefore:   testedBefore,
			FirstName:      context.PostForm("first_name"),
			LastName:       context.PostForm("last_name"),
			ProfilePhoto:   userImageLink,
			EducationLevel: context.PostForm("education_level"),
			Phone:          context.PostForm("phone"),
		}).Error; err != nil {
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

	context.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": updatedUser})
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
