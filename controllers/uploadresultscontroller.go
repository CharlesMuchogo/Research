package controllers

import (
	"awesomeProject/auth"
	"awesomeProject/models"
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func Upload(context *gin.Context) {
	var results models.Results
	if err := context.ShouldBindJSON(&results); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	tokenString := context.GetHeader("Authorization")

	claims, err := auth.GetUserDetailsFromToken(tokenString)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	email := claims.Email
	firstName := claims.FirstName
	lastName := claims.LastName
	phone := claims.Phone

	spreadsheetID := "12nLNSb0n9kdpoETkPQjTEx-MzmbuEPdq5YCrnXqT2JU"
	credentialsFile := "./credentials.json"
	client, err := utils.GetClient(credentialsFile)
	if err != nil {
		log.Fatalf("Error getting Google Sheets client: %v", err)
	}

	now := time.Now()
	formattedDateTime := now.Format("02/01/2006 15:04")

	sheetRange := "Sheet1!A1:J5"
	values := [][]interface{}{
		{firstName, lastName, phone, email, results.Results, results.PartnerResults, results.Image, results.PartnerImage, formattedDateTime},
	}

	err = utils.WriteDataToSpreadsheet(client, spreadsheetID, sheetRange, values)
	if err != nil {
		log.Fatalf("Error writing data to spreadsheet: %v", err)
	}

	fmt.Println("Data written to the spreadsheet successfully!")
	context.JSON(http.StatusOK, gin.H{"message": "Test submitted successfully"})
}
