package utils

import (
	"awesomeProject/models"
	"fmt"
	"net/smtp"
)

func SendMail(user models.User, tokenString string) {
	from := "xxxxxxxx" // Update with your email
	password := "xxxxxxxxx" // Update with your email password
	to := user.Email

	link := "http://192.168.0.77:9000/api/reset_password?token=" + tokenString
	sendBody := fmt.Sprintf(`<html>
	<head>
		<title>Password Reset</title>
	</head>
	<body>
		<p>Hello %v,</p>
		<p>Someone requested to change your password.</p>
		<p>Click <a href="%v">here</a> to reset your password.</p>
	</body>
	</html>`, user.FirstName, link)

	msg := []byte(fmt.Sprintf("To: %s\r\n", to) +
		"Subject: Password Reset\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" + // Specify charset
		"\r\n" +
		sendBody)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, msg)

	if err != nil {
		fmt.Println("Error sending email: ", err)
		return
	}

	fmt.Println("Email sent successfully.")

}
