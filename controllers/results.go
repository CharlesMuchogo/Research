package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/database"
	"awesomeProject/fcm"
	"awesomeProject/models"
	"awesomeProject/models/dto"
	"awesomeProject/utils"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"time"
)

func UploadResults(context *gin.Context) {

	var wg sync.WaitGroup
	wg.Add(1)

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
		userImageLink, err = utils.SavePhoto(userTestResultsPhoto, phone)
		if err != nil {
			fmt.Println(err.Error())
			context.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading test image"})
			return
		}
	}

	var partnerImageLink string
	if partnerTestResultsPhoto != nil {
		partnerImageLink, err = utils.SavePhoto(partnerTestResultsPhoto, phone)
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

	user, err := database.GetUserById(userId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		context.Abort()
		return
	}

	results := models.Results{
		Results:        context.PostForm("results"),
		PartnerResults: context.PostForm("partner_results"),
		Image:          userImageLink,
		PartnerImage:   partnerImageLink,
		CareOption:     context.PostForm("care_option"),
		Date:           formattedDateTime,
		Status:         "Pending",
		UserId:         userId,
		User:           user,
	}

	record := database.Instance.Create(&results)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		context.Abort()
		return
	}

	go func() {
		defer wg.Done()
		fcm.SendNotification("Test results submission.", "Thank you for submitting your test results. We are reviewing your results, we will be in touch shortly", user.Email, nil)
	}()

	context.JSON(http.StatusOK, gin.H{"message": "Test results submitted successfully. Wait for the approval of your results.", "result": results})
	wg.Wait()
}

func GetResults(context *gin.Context) {
	var results []models.Results

	tokenString := context.GetHeader("Authorization")

	user, _ := auth.GetUserDetailsFromToken(tokenString)

	if err := database.Instance.Preload("User").Where("user_id = ? AND deleted = false", user.ID).Find(&results).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Results fetched successfully", "results": results})
}
func GetAllResults(context *gin.Context) {
	var results []models.Results

	if err := database.Instance.Preload("User").Order("id DESC").Find(&results).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Results fetched successfully", "results": results})
}

func UpdateResults(context *gin.Context) {
	var request dto.ResultDTO
	var results models.Results
	var wg sync.WaitGroup
	wg.Add(1)

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	if err := database.Instance.Preload("User").Where("uuid = ?", request.UUID).Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid results"})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
		return
	}

	if results.Date == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid results"})
		return
	}

	results.Results = request.Results
	results.PartnerResults = request.PartnerResults
	results.Status = request.Status
	results.Reason = request.Reason

	if err := database.Instance.Save(&results).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
		return
	}

	go func() {
		defer wg.Done()
		fcm.SendNotification("Test results update", "Your test results feedback is ready. Please check the test page to view your results", results.User.Email, nil)
	}()

	context.JSON(http.StatusOK, gin.H{"message": "Results updated successfully", "results": results})
	wg.Wait()
}

func DeleteResults(context *gin.Context) {
	var results models.Results

	uuid := context.Query("uuid")

	if err := database.Instance.Preload("User").Where("uuid = ?", uuid).Find(&results).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid test id"})
			return
		}

		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
		return
	}

	if results.Date == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid results"})
		return
	}

	results.Deleted = true

	if err := database.Instance.Save(&results).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong, try again"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Test deleted successfully"})
}
