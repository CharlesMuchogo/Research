package utils

import (
	"awesomeProject/models"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

type EmailData struct {
	UserName       string
	ResetLink      string
	CompanyName    string
	CompanyAddress string
	SupportLink    string
	PrivacyLink    string
}

func SendForgotPasswordEmail(user models.User, tokenString string) {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")
	domain := os.Getenv("APP_RUNNING_DOMAIN")
	to := user.Email

	link := fmt.Sprintf("https://%s/api/reset_password?id=%s", domain, tokenString)

	tmpl, err := template.ParseFiles("templates/forgot_password.html")

	if err != nil {
		fmt.Printf("Error loading template: %v\n", err)
		return
	}

	data := EmailData{
		UserName:       user.FirstName,
		ResetLink:      link,
		CompanyName:    "SmartTest",
		CompanyAddress: "",
		SupportLink:    fmt.Sprintf("%s/support", domain),
		PrivacyLink:    fmt.Sprintf("%s/privacy", domain),
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	msg := []byte(fmt.Sprintf("To: %s\r\n", to) +
		"Subject: Password Reset Request\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body.String())

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, msg)

	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		return
	}

	fmt.Printf("Password reset email sent successfully to %s\n", to)
}
