package sendmailserviceproxy

import (
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

//SendGridEmailService
type SendGridEmailService struct {
	apiKey string
}

//NewSendGridEmailService ctor
func NewSendGridEmailService(apiKey string) EmailService {
	return &SendGridEmailService{
		apiKey: apiKey,
	}
}

//Send message using SendGrid
func (service *SendGridEmailService) Send(message EmailMessage) error {
	from := mail.NewEmail(message.From, message.From)
	to := mail.NewEmail(message.To, message.To)
	msg := mail.NewSingleEmail(from, message.Subject, to, message.PlainBody, message.PlainBody)

	client := sendgrid.NewSendClient(service.apiKey)

	_, err := client.Send(msg)
	return err
}
