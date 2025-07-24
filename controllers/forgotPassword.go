package controllers

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func ForgotPassword(context *gin.Context) {
	var request ForgotPasswordRequest
	var user models.User

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	// check if email exists
	if err := database.DbInstance.Where("email = ?", request.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "User not found. Please register"})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tokenString := uuid.New().String()

	forgotPasswordRequest := models.ForgotPassword{
		Token:  tokenString,
		UserId: user.ID,
	}

	if err := database.DbInstance.Save(&forgotPasswordRequest).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	go utils.SendForgotPasswordEmail(user, tokenString)
	context.JSON(http.StatusOK, gin.H{"message": "Please check your email for reset instructions"})
}

func UpdatePassword(context *gin.Context) {
	var request models.ForgotPasswordRequest
	var user models.User
	var forgotPasswordRecord models.ForgotPassword

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := database.DbInstance.Unscoped().
		Where("token = ? AND deleted_at IS NULL", request.Token).
		First(&forgotPasswordRecord).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "You have used an invalid link"})
			return
		}
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if time.Since(forgotPasswordRecord.CreatedAt) > 10*time.Minute {
		context.JSON(http.StatusBadRequest, gin.H{"message": "You have used an expired link"})
		return
	}

	if err := database.DbInstance.Where("id = ?", forgotPasswordRecord.UserId).First(&user).Error; err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := user.HashPassword(request.Password); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := database.DbInstance.Model(&models.User{}).
		Where("email = ?", user.Email).
		Select("Password").
		Updates(user).Error; err != nil {
		fmt.Printf("Error updating user %s", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := database.DbInstance.Model(&models.ForgotPassword{}).
		Where("token = ?", request.Token).
		Update("deleted_at", time.Now()).
		Error; err != nil {
		fmt.Printf("Error updating deleted_at: %s\n", err.Error())
	}

	context.JSON(http.StatusOK, gin.H{"message": "Password reset successfully", "user": user})
}
