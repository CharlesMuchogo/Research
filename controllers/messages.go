package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
	"net/http"
)

func SendMessages(context *gin.Context) {

	var message models.Message

	token := context.GetHeader("Authorization")

	if err := context.ShouldBindJSON(&message); err != nil {

		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error() /*"Invalid message"*/})
		context.Abort()
		return
	}

	userFromToken, err := auth.GetUserDetailsFromToken(token)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user details"})
		context.Abort()
		return
	}

	message.UserId = userFromToken.ID

	fmt.Printf("message to save is %v", message)

	if err := database.DbInstance.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "timestamp"}},
		DoUpdates: clause.AssignmentColumns([]string{"message", "sender", "user_id"}),
	}).Create(&message).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Message sent successfully", "data": message})

}

func GetMessages(context *gin.Context) {
	var messages []models.Message
	token := context.GetHeader("Authorization")

	user, _ := auth.GetUserDetailsFromToken(token)

	if err := database.DbInstance.Where("user_id = ?", user.ID).Find(&messages).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error() /*"Something went wrong, try again"*/})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Results fetched successfully", "data": messages})

}
