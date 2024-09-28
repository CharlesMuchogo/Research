package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/models/dto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var userRequest dto.UserDTO

	if err := context.ShouldBindJSON(&userRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	if err := userRequest.HashPassword(userRequest.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	user := models.User{
		FirstName:      userRequest.FirstName,
		LastName:       userRequest.LastName,
		Phone:          userRequest.Phone,
		Email:          userRequest.Email,
		Password:       userRequest.Password,
		ProfilePhoto:   userRequest.ProfilePhoto,
		Age:            userRequest.Age,
		EducationLevel: userRequest.EducationLevel,
		TestedBefore:   userRequest.TestedBefore,
		Gender:         userRequest.Gender,
	}

	record := database.Instance.Create(&user)

	if record.Error != nil {
		if strings.Contains(record.Error.Error(), "users_email_key") {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Email has already been used"})
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

	if err := database.Instance.Save(&user).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update user details"})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Details updated successfully", "user": user})
}
