package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/fcm"
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	DeviceId string `json:"deviceId"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func Login(context *gin.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	// check if email exists and password is correct
	record := database.DbInstance.Where("phone = ? OR email = ?", request.Email, request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid credentials"})
		context.Abort()
		return
	}
	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid credentials"})
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

	go func() {
		defer wg.Done()
		go fcm.RegisterTopic(user.Email, request.DeviceId)
	}()

	wg.Wait()

}

func AdminLogin(context *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

	var request TokenRequest
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	// check if email exists and password is correct
	record := database.DbInstance.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid credentials"})
		context.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "invalid credentials"})
		context.Abort()
		return
	}

	if user.Role != "admin" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized access"})
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

	go func() {
		defer wg.Done()
		go fcm.RegisterTopic(user.Email, request.DeviceId)
	}()
	wg.Wait()
}

func ForgotPassword(context *gin.Context) {
	var request ForgotPasswordRequest
	var user models.User
	var wg sync.WaitGroup
	wg.Add(1)
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}
	// check if email exists
	record := database.DbInstance.Where("email = ?", request.Email).First(&user)

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

	context.JSON(http.StatusOK, gin.H{"message": "Please check your email for reset instructions", "token": tokenString})

	go func() {
		defer wg.Done()
		go utils.SendForgotPasswordEmail(user, tokenString)
	}()
	wg.Wait()
}

func ResetPassword(context *gin.Context) {
	context.HTML(200, "resetPassword.html", nil)
}

func UpdatePassword(context *gin.Context) {
	tokenString := context.GetHeader("Authorization")

	fmt.Printf("token is %s \n", tokenString)
	context.HTML(http.StatusOK, "privacyPolicy.html", nil)
}

func DeleteAccountForm(context *gin.Context) {
	context.HTML(http.StatusOK, "DeleteAccount.html", nil)
}

func PrivacyPolicy(context *gin.Context) {
	context.HTML(http.StatusOK, "privacyPolicy.html", nil)
}
func TermsAndConditions(context *gin.Context) {
	context.HTML(http.StatusOK, "termsAndConditions.html", nil)
}
func NotFound(context *gin.Context) {
	context.HTML(http.StatusNotFound, "404.html", nil)
}
