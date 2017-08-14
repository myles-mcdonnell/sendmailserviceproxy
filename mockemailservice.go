package sendmailserviceproxy

import (
	"errors"
	"github.com/myles-mcdonnell/sendmailserviceproxy/logging"
	"time"
)

// MockEmailService breaks the dependency on real email service implementations for the purpose of testing
type MockEmailService struct{}

// Send will result in no error, error 500 or a timeout depending on the subject of the message.  timeout, fail or anything else for no error.
func (mockEmailService MockEmailService) Send(message EmailMessage) error {

	logging.LogDebug(logging.MockEmailServiceInvoked, nil, message)

	if message.Subject == "timeout" {
		for true {
			time.Sleep(time.Second * 1)
		}
	} else if message.Subject == "fail" {
		return errors.New("500")
	}

	return nil
}
