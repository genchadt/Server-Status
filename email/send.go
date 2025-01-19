package email

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"serverstatus/config"

	"github.com/wneessen/go-mail"
)

// SendEmail sends an email with the given subject and body to the recipient at the given
// address. It uses the SMTP server and credentials specified in the function body.
// Note that the auth information should be replaced with a secure storage mechanism.
func SendEmail(subject, body, to string) error {
	// Get credentials from environment variables or secure storage
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Convert SMTP port to integer
	smtpPortInt, err := strconv.Atoi(cfg.SmtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	// Use go-mail library with TLS
	m := mail.NewMsg()
	if err := m.FromFormat(cfg.FromName, cfg.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}
	if err := m.To(to); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	m.Subject(subject)
	m.SetBodyString(mail.TypeTextHTML, body)

	// Establish a TLS connection
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.SmtpHost,
	}

	c, err := mail.NewClient(cfg.SmtpHost,
		mail.WithPort(smtpPortInt),
		mail.WithUsername(cfg.SmtpUser),
		mail.WithPassword(cfg.SmtpPass),
		mail.WithTLSConfig(tlsConfig),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %v", err)
	}

	if err := c.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
