package routes

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/myles-mcdonnell/sendmailserviceproxy"
	"github.com/myles-mcdonnell/sendmailserviceproxy/logging"
	"github.com/myles-mcdonnell/sendmailserviceproxy/models"
	"github.com/myles-mcdonnell/sendmailserviceproxy/restapi/operations/email"
)

// Encapsulates n email services and wraps them in circuit breakers
type EmailHandler struct {
	circuits []*sendmailserviceproxy.Circuit
}

// NewEmailHandler; ctor for EmailHandler
func NewEmailHandler(emailServices []sendmailserviceproxy.EmailService, pollMessageAddress string) *EmailHandler {

	pollMessage := sendmailserviceproxy.EmailMessage{To: pollMessageAddress, From: pollMessageAddress, Subject: "poll", PlainBody: "poll"}
	circuits := make([]*sendmailserviceproxy.Circuit, len(emailServices))

	for index, service := range emailServices {
		circuits[index] = sendmailserviceproxy.NewCircuit(service.Send, pollMessage)
	}

	return &EmailHandler{
		circuits: circuits,
	}
}

// Send message; This method will attempt to send an email by iterating over it's circuits.
// The result from the first closed circuit will be returned to the caller.
// It is not possible to determine if a failed or timed out call may have resulted in the downstream system
// sending the message therefore the decision to retry or not is deferred to the upstream system.
func (emailHandler *EmailHandler) Send(message sendmailserviceproxy.EmailMessage) error {

	var lastError error

	for _, circuit := range emailHandler.circuits {

		var status sendmailserviceproxy.CallResult
		lastError, status = circuit.Call(message)

		if lastError == nil {
			logging.LogDebug(logging.PostEmailDebug, nil, "emailHandler.Send OK")
			return nil
		}

		logging.LogDebug(logging.PostEmailDebug, nil, "emailHandler.Send Error: "+lastError.Error())

		//Try another circuit if failure due to circuit being open
		if status != sendmailserviceproxy.CIRCUIT_OPEN {
			logging.LogDebug(logging.PostEmailDebug, nil, "emailHandler.Send CIRCUIT_OPEN")
			return lastError
		}
	}

	return lastError
}

// This method processes the API request to send a message.
func (emailHandler EmailHandler) EmailPost(params email.PostEmailParams, data interface{}) middleware.Responder {

	defer func() {
		logging.LogInfo(logging.PostEmailEnd, params.HTTPRequest, nil)
	}()

	logging.LogInfo(logging.PostEmailStart, params.HTTPRequest, nil)

	msg := sendmailserviceproxy.EmailMessage{
		To:        *params.Email.Toaddress,
		From:      *params.Email.Fromaddress,
		Subject:   params.Email.Subject,
		PlainBody: params.Email.Plaintextbody,
	}

	logging.LogDebug(logging.PostEmailDebug, params.HTTPRequest, msg)

	err := emailHandler.Send(msg)

	if err == nil {
		return email.NewPostEmailOK()
	}

	logging.LogError(logging.PostEmailError, params.HTTPRequest, err.Error())
	return email.NewPostEmailDefault(500).WithStatusCode(500).WithPayload(&models.ErrorMessage{Message: SPtr(err.Error())})
}

func SPtr(s string) *string { return &s }
