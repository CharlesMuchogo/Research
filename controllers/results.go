package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func UploadResults(context *gin.Context) {

	userTestResultsPhoto, _ := context.FormFile("user_photo")
	partnerTestResultsPhoto, _ := context.FormFile("partner_photo")

	tokenString := context.GetHeader("Authorization")

	claims, err := auth.GetUserDetailsFromToken(tokenString)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	phone := claims.Phone
	userId := claims.ID

	var userImageLink string
	if userTestResultsPhoto != nil {
		userImageLink, err = utils.SavePhoto(context, userTestResultsPhoto, phone)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading test image"})
			return
		}
	}

	var partnerImageLink string
	if partnerTestResultsPhoto != nil {
		partnerImageLink, err = utils.SavePhoto(context, partnerTestResultsPhoto, phone)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading partner test image"})
			return
		}
	}

	nairobiLocation, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	now := time.Now().In(nairobiLocation)
	formattedDateTime := now.Format("02/01/2006 15:04")

	results := models.Results{
		Results:        context.PostForm("results"),
		PartnerResults: context.PostForm("partner_results"),
		Image:          userImageLink,
		PartnerImage:   partnerImageLink,
		CareOption:     context.PostForm("care_option"),
		Date:           formattedDateTime,
		Status:         "Pending",
		UserId:         userId,
	}

	record := database.Instance.Create(&results)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Test submitted successfully", "data": results})
}

func GetResults(context *gin.Context) {
	var results []models.Results
	fetchAll := context.Query("all")

	tokenString := context.GetHeader("Authorization")

	user, _ := auth.GetUserDetailsFromToken(tokenString)

	if fetchAllBool, err := strconv.ParseBool(fetchAll); err == nil && fetchAllBool {
		if err := database.Instance.Preload("User").Find(&results).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
			return
		}
	} else {
		if err := database.Instance.Preload("User").Where("user_id = ?", user.ID).Find(&results).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
			return
		}
	}
	context.JSON(http.StatusOK, gin.H{"message": "Results fetched successfully", "results": results})
}
