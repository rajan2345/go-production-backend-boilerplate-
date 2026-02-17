package email

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/pkg/errors"
	"github.com/rajan2345/go-boilerplate/internal/config"
	"github.com/resend/resend-go/v2"
	"github.com/rs/zerolog"
)

// this file contains the actual logic which will send the emails
// the backgroud job created , its task is to trigger the function which will send the email it was just a backgroud task processor,
// but to send the email we need a email client and service
// for this purpose we will be using the "Resend"

// Resend a modern email provider ,
// -- register -> get a key -> put it into environment -> and use it
// Declare a new client for sending the email

type Client struct {
	client *resend.Client
	logger *zerolog.Logger
}

func NewClient(cfg *config.Config, logger *zerolog.Logger) *Client {
	return &Client{
		client: resend.NewClient(cfg.Integration.ResendAPIKey),
		logger: logger,
	}
}

// we need a function  which will use this initialized client sdk and send the email
func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {
	tmplPath := fmt.Sprintf("%s/%s.html", "templates/emails", templateName) // whenever email is sent this will be read alongwith email

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return errors.Wrapf(err, "failed to execute email template %s", templateName)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return errors.Wrapf(err, "failed to execute email template %s", templateName)
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s> ", "Boilerplate", "onboarding@resend.dev"),
		To:      []string{to},
		Subject: subject,
		Html:    body.String(),
	}
	_, err = c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
