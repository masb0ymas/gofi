package lib

import (
	"bytes"
	"errors"
	"fmt"
	"gofi/src/config"
	"log"
	"net/smtp"
	"text/template"

	"github.com/masb0ymas/go-utils/pkg"
)

type SMTPConfig struct {
	Host     string
	Port     string
	From     string
	Username string
	Password string
}

func NewSMTPConfig() *SMTPConfig {
	return &SMTPConfig{
		Host:     config.Env("SMTP_HOST", "smtp.gmail.com"),
		Port:     config.Env("SMTP_PORT", "587"),
		From:     config.Env("SMTP_FROM", "your_email_from"),
		Username: config.Env("SMTP_USERNAME", "your-smtp-username"),
		Password: config.Env("SMTP_PASSWORD", "your-smtp-password"),
	}
}

type SendEmailParams struct {
	Subject          string      `json:"subject"`
	To               string      `json:"to"`
	Data             interface{} `json:"data"`
	FilenameTemplate string      `json:"filename_template"`
}

// SendEmail sends an email
// @example
//
//	SendEmail(SendEmailParams{
//		Subject:          "Your Subject Email",
//		To:               "your_audience@mail.com",
//		Data:             Object,
//		FilenameTemplate: "template.html",
//	})
func SendEmail(value SendEmailParams) error {
	config := NewSMTPConfig()

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	from := config.From
	send_to := []string{value.To}

	// Load the HTML template
	template := fmt.Sprintf("public/email-template/%s", value.FilenameTemplate)
	buf, err := ParseTemplate(template, value.Data)
	if err != nil {
		return errors.New("error loading template: " + err.Error())
	}

	// Email content
	new_subject := "Subject: " + value.Subject + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := buf
	message := []byte(new_subject + mime + body)

	// Authentication
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// Sending email
	err = smtp.SendMail(addr, auth, from, send_to, message)
	if err != nil {
		return errors.New("error sending email: " + err.Error())
	}

	// Print success message
	msg := pkg.Println("SMTP", fmt.Sprintf("Sent to: %s", value.To))
	log.Println(msg, "Email sent successfully")

	return nil
}

// ParseTemplate parses a template file
func ParseTemplate(filename string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		return "", errors.New("error parsing template: " + err.Error())
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", errors.New("error executing template: " + err.Error())
	}

	return buf.String(), nil
}
