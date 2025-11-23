package services

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"gofi/internal/config"

	"github.com/resend/resend-go/v3"
)

type EmailService struct {
	Config config.ConfigResend
}

type SendEmailParams struct {
	Subject      string
	To           string
	Data         interface{}
	HtmlTemplate string
}

func (s EmailService) SendEmail(value SendEmailParams) (string, error) {
	client := resend.NewClient(s.Config.ApiKey)

	// Load the HTML template
	htmlStr, err := ParseTemplate(value.HtmlTemplate, value.Data)
	if err != nil {
		return "", fmt.Errorf("error loading template: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("GoFi <%s>", s.Config.FromEmail),
		To:      []string{value.To},
		Html:    htmlStr,
		Subject: value.Subject,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Printf("Error sending email: %s\n", err.Error())
		return "", err
	}
	fmt.Printf("Email sent with ID: %s\n", sent.Id)
	return sent.Id, nil
}

// Parse Template from file path
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
