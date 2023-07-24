package controllers

import (
	"awesomeProject/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Upload(context *gin.Context) {

	spreadsheetID := "12nLNSb0n9kdpoETkPQjTEx-MzmbuEPdq5YCrnXqT2JU"
	credentialsFile := "./credentials.json"
	client, err := utils.GetClient(credentialsFile)
	if err != nil {
		log.Fatalf("Error getting Google Sheets client: %v", err)
	}
	sheetRange := "Sheet1!A1:G5"
	values := [][]interface{}{
		{"Baby Boo", 23, "Kenya"},
		{"chaos", 22, "Germany"},
	}

	err = utils.WriteDataToSpreadsheet(client, spreadsheetID, sheetRange, values)
	if err != nil {
		log.Fatalf("Error writing data to spreadsheet: %v", err)
	}

	fmt.Println("Data written to the spreadsheet successfully!")
	context.JSON(http.StatusOK, gin.H{"message": "Test submitted successfully"})
}
