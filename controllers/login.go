package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPaswordRequest struct {
	Email string `json:"email"`
}

func GenerateToken(context *gin.Context) {
	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	// check if email exists and password is correct
	record := database.Instance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "invalid credential"})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "invalid credentials"})
		context.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login success", "user": user, "token": tokenString})
}

func ForgotPasword(context *gin.Context) {
	var request ForgotPaswordRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}
	// check if email exists
	record := database.Instance.Where("email = ?", request.Email).First(&user)

	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": record.Error.Error()})
		context.Abort()
		return
	}
	tokenString, err := auth.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		context.Abort()
		return
	}
	go utils.SendMail(user, tokenString)
	context.JSON(http.StatusOK, gin.H{"message": "Please check your email for reset instructions", "token": tokenString})
}

func ResetPassword(context *gin.Context) {
	context.HTML(200, "resetPassword.html", nil)
}

func DeleteAccountForm(context *gin.Context) {
	context.HTML(200, "DeleteAccount.html", nil)
}
func PrivacyPolicy(context *gin.Context) {
	context.HTML(200, "privacyPolicy.html", nil)
}
