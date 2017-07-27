package sendmailserviceproxy

type Configuration struct {
	MailGunDomain, MailGunApiKey, MailgunPublicApiKey string
	SendGridApiKey                                    string
}

type EmailMessage struct {
	Subject   string
	PlainBody string
	From      string
	To        string
}

type EmailService interface {
	Send(message EmailMessage) error
}
