package restapi

import (
	"fmt"
	"github.com/myles-mcdonnell/sendmailserviceproxy"
	"github.com/myles-mcdonnell/sendmailserviceproxy/logging"
	"github.com/myles-mcdonnell/sendmailserviceproxy/restapi/operations"
	"github.com/myles-mcdonnell/sendmailserviceproxy/routes"
	"gopkg.in/myles-mcdonnell/loglight.v3"
	"os"
	"strconv"
)

type config struct {
	mgDomain            string
	mgApiKey            string
	sgApiKey            string
	logOutputFormat     string
	logOutputPretty     bool
	logDebug            bool
	useMockEmailService bool
	pollMessageAddress  string
	apiKey              string
}

// InitApi initialises the server
func InitApi(api *operations.SendmailserviceproxyAPI) {

	conf := config{
		mgDomain:            getMandatoryEnvStr("SMSP_MG_DOMAIN"),
		mgApiKey:            getMandatoryEnvStr("SMSP_MG_API_KEY"),
		sgApiKey:            getMandatoryEnvStr("SMSP_SG_API_KEY"),
		logOutputFormat:     getOptionalEnvStr("SMSP_LOG_OUTPUT_FORMAT", "DEBUG"),
		logDebug:            getEnvBool("SMSP_LOG_OUTPUT_DEBUG"),
		useMockEmailService: getEnvBool("SMSP_MOCK_EMAIL_SERVICE"),
		pollMessageAddress:  getMandatoryEnvStr("SMSP_POLL_MESSAGE_ADDRESS"),
		apiKey:              getMandatoryEnvStr("SMSP_API_KEY"),
	}

	logging.Initialise(loglight.NewLogger(
		conf.logDebug,
		logging.BuildFormatter(conf.logOutputFormat),
	))

	api.Logger = func(msg string, args ...interface{}) {
		logging.LogInfo(logging.Api, nil, args)
	}

	if conf.useMockEmailService {
		logging.LogInfo(logging.Api, nil, "Using Mock Email Service")
		routes.BindRoutes(api, []sendmailserviceproxy.EmailService{
			sendmailserviceproxy.MockEmailService{}},
			conf.pollMessageAddress)
	} else {
		routes.BindRoutes(api, []sendmailserviceproxy.EmailService{
			sendmailserviceproxy.NewMailGunEmailService(conf.mgDomain, conf.mgApiKey),
			sendmailserviceproxy.NewSendGridEmailService(conf.sgApiKey),
		},
			conf.pollMessageAddress)
	}
}

func getEnvBool(key string) bool {
	str := os.Getenv(key)
	bool, _ := strconv.ParseBool(str)

	return bool
}

func getMandatoryEnvStr(key string) string {
	str := os.Getenv(key)
	if str == "" {
		fmt.Println("need envvar: " + key)
	}

	return str
}

func getOptionalEnvStr(key string, def string) string {
	str := os.Getenv(key)
	if str == "" {
		return def
	}

	return str
}
