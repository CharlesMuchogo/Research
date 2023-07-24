package utils

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
)

func GetClient(credentialsFile string) (*sheets.Service, error) {
	// Read the content of the credentials file
	credentialsJSON, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %v", err)
	}

	// Create credentials from the JSON content
	credentials, err := google.CredentialsFromJSON(context.Background(), credentialsJSON, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("unable to create Google Sheets credentials: %v", err)
	}

	// Create the HTTP client with the credentials
	httpClient := oauth2.NewClient(context.Background(), credentials.TokenSource)

	// Create the Google Sheets API client
	sheetsClient, err := sheets.New(httpClient)
	if err != nil {
		return nil, fmt.Errorf("unable to create Google Sheets API client: %v", err)
	}

	return sheetsClient, nil
}

// writeDataToSpreadsheet writes data to a specified range in a Google Spreadsheet.
func WriteDataToSpreadsheet(client *sheets.Service, spreadsheetID, sheetRange string, values [][]interface{}) error {
	writeRequest := &sheets.ValueRange{
		Values: values,
	}

	_, err := client.Spreadsheets.Values.Append(spreadsheetID, sheetRange, writeRequest).
		ValueInputOption("RAW").
		Context(context.Background()).
		Do()

	return err
}
