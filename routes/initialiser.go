package routes

import (
	"github.com/myles-mcdonnell/sendmailserviceproxy"
	"github.com/myles-mcdonnell/sendmailserviceproxy/restapi/operations"
	"github.com/myles-mcdonnell/sendmailserviceproxy/restapi/operations/email"
	"github.com/myles-mcdonnell/sendmailserviceproxy/restapi/operations/healthcheck"
)

// BindRoutes Binds the route handlers to the API
func BindRoutes(api *operations.SendmailserviceproxyAPI, emailServices []sendmailserviceproxy.EmailService, pollMessageAddress string) {

	api.HealthcheckGetHealthcheckHandler = healthcheck.GetHealthcheckHandlerFunc(HealthcheckGet)

	emailHandler := NewEmailHandler(emailServices, pollMessageAddress)
	api.EmailPostEmailHandler = email.PostEmailHandlerFunc(emailHandler.EmailPost)
}
