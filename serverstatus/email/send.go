package email

import (
	"fmt"
	"net/smtp"
	// "github.com/wneessen/go-mail"
)

// TODO: Decide on method for securely storing email credentials

func SendEmail(subject, body, to string) error {
	// Auth info goes here
	auth := smtp.PlainAuth("", "your_email@example.com", "your_password", "smtp.example.com")

	from := "Lightsail Web Updates <lightsail-updates@thecollectivegc.com>"
	msg := "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	msg += fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n\r\n", subject)
	msg += body

	// SMTP info goes here
	err := smtp.SendMail("smtp.example.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
