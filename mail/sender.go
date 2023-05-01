package mail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

const (
	smtpAuthAddress   = "email-smtp.eu-central-1.amazonaws.com"
	smtpServerAddress = "email-smtp.eu-central-1.amazonaws.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
}

type AwsSesSender struct {
	name             string
	fromEmailAddress string
	smtpUsername     string
	smtpPassword     string
}

func NewAwsSesSender(name string, fromEmailAddress string, smtpUsername string, smtpPassword string) EmailSender {
	return &AwsSesSender{
		name:             name,
		fromEmailAddress: fromEmailAddress,
		smtpUsername:     smtpUsername,
		smtpPassword:     smtpPassword,
	}
}

func (sender *AwsSesSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, attachFile := range attachFiles {
		_, err := e.AttachFile(attachFile)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", attachFile, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", sender.smtpUsername, sender.smtpPassword, smtpAuthAddress)
	return e.Send(smtpServerAddress, smtpAuth)
}
