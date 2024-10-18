package controllers

import (
	"awesomeProject/database"
	"awesomeProject/models"
	"awesomeProject/models/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CreateClinic(context *gin.Context) {
	var clinicDTO dto.CreateClinicDTO
	if err := context.ShouldBindJSON(&clinicDTO); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		context.Abort()
		return
	}

	clinic := clinicDTO.ToClinic()

	record := database.Instance.Create(&clinic)

	if record.Error != nil {

		if strings.Contains(record.Error.Error(), "clinics_name_key") {
			context.JSON(http.StatusBadRequest, gin.H{"message": "This clinic already exists"})
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"message": record.Error.Error()})
		}
		context.Abort()
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Clinic added successfully", "clinic": clinic})
}

func GetClinics(context *gin.Context) {
	var clinics []models.Clinic

	record := database.Instance.Find(&clinics)

	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong. Try again"})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Clinics fetched successfully", "clinics": clinics})
}
