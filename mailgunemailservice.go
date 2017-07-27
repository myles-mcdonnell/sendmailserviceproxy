package sendmailserviceproxy

import (
	"github.com/mailgun/mailgun-go"
)

//MailGunEmailService
type MailGunEmailService struct {
	mailGun mailgun.Mailgun
}

//NewMailGunEmailService ctor
func NewMailGunEmailService(domain string, apikey string) EmailService {

	return &MailGunEmailService{
		mailGun: mailgun.NewMailgun(
			domain,
			apikey,
			apikey)}
}

//Send message using mail gun
func (service MailGunEmailService) Send(message EmailMessage) error {

	msg := service.mailGun.NewMessage(
		message.From,
		message.Subject,
		message.PlainBody)

	msg.SetHtml(message.PlainBody)
	msg.AddRecipient(message.To)

	_, _, err := service.mailGun.Send(msg)

	return err
}
