package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func LoadSMTPConfig() SMTPConfig {
	return SMTPConfig{
		Host:     strings.TrimSpace(os.Getenv("SMTP_HOST")),
		Port:     strings.TrimSpace(os.Getenv("SMTP_PORT")),
		Username: strings.TrimSpace(os.Getenv("SMTP_USER")),
		Password: strings.TrimSpace(os.Getenv("SMTP_PASSWORD")),
		From:     strings.TrimSpace(os.Getenv("SMTP_FROM")),
	}
}

func (c SMTPConfig) Enabled() bool {
	return c.Host != "" && c.Port != "" && c.From != ""
}

func (c SMTPConfig) Send(to, subject, body string) error {
	if !c.Enabled() {
		return fmt.Errorf("smtp is not configured")
	}

	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", c.From),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	var auth smtp.Auth
	if c.Username != "" && c.Password != "" {
		auth = smtp.PlainAuth("", c.Username, c.Password, c.Host)
	}

	return smtp.SendMail(addr, auth, c.From, []string{to}, []byte(msg))
}
