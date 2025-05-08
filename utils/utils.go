package utils

import (
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	from := "bomaflorent23@gmail.com"
	password := "fnxf jdma nadm mypm "

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail(
		"smtp.gmail.com:587", // Replace with your SMTP server
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from,
		[]string{to},
		[]byte(msg),
	)
	return err
}
