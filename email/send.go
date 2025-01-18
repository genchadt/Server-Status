package email

import (
	"fmt"
	"net/smtp"
	"os"
	// "github.com/wneessen/go-mail"
)

// SendEmail sends an email with the given subject and body to the recipient at the given
// address. It uses the SMTP server and credentials specified in the function body.
// Note that the auth information should be replaced with a secure storage mechanism.
func SendEmail(subject, body, to string) error {
	// Get credentials from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP_HOST and SMTP_PORT environment variables must be set")
	}

	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	if smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP_USER and SMTP_PASS environment variables must be set")
	}

	fromEmail := os.Getenv("FROM_EMAIL")

	if fromEmail == "" {
		return fmt.Errorf("FROM_EMAIL environment variable must be set")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	from := fmt.Sprintf("Lightsail Web Updates <%s>", fromEmail)

	msg := "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	msg += fmt.Sprintf("From: %s\r\n", from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n\r\n", subject)
	msg += body

	smtpAddr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	// SMTP info goes here
	err := smtp.SendMail(smtpAddr, auth, fromEmail, []string{to}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}
