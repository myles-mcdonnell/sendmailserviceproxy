package routes

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/myles-mcdonnell/sendmailserviceproxy/logging"
	"github.com/myles-mcdonnell/sendmailserviceproxy/restapi/operations/healthcheck"
)

// HealthcheckGet; This method simply returns a 200 code and is for the purpose of automated monitoring
// It could be useful to include circuit breaker state in the response although currently is does not
func HealthcheckGet(params healthcheck.GetHealthcheckParams) middleware.Responder {

	defer func() {
		logging.LogInfo(logging.GetHealthcheckEnd, params.HTTPRequest, nil)
	}()
	logging.LogInfo(logging.GetHealthcheckStart, params.HTTPRequest, nil)

	return healthcheck.NewGetHealthcheckOK()
}
