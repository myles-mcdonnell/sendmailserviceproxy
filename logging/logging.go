package logging

import (
	"fmt"
	"gopkg.in/myles-mcdonnell/loglight.v3"
	"gopkg.in/satori/go.uuid.v1"
	"net/http"
	"os"
	"time"
)

var log *loglight.Logger
var hostAddress string

const requestKeyHeaderName string = "RequestKey"

type LogEventType string

const (
	GetHealthcheckStart     LogEventType = "Get Healthcheck - Start"
	GetHealthcheckEnd       LogEventType = "Get Healthcheck - End"
	PostEmailStart          LogEventType = "Post Email - Start"
	PostEmailEnd            LogEventType = "Post Email - End"
	PostEmailDebug          LogEventType = "Post Email - Debug"
	PostEmailError          LogEventType = "Post Email - Error"
	Panic                   LogEventType = "Panic"
	Api                     LogEventType = "Api"
	Unauthorised            LogEventType = "Unauthorised Request"
	MockEmailServiceInvoked LogEventType = "MockEmailServiceInvoked"
)

func Initialise(logger *loglight.Logger) {

	log = logger
	hostAddress, _ = os.Hostname()
}

var jsonOneLineFormatter = loglight.NewJsonLogFormatter(false)

func BuildFormatter(key string) func(data loglight.LogEntry) string {

	if key == "JSON_ONELINE" {
		return jsonOneLineFormatter.Format
	}

	if key == "JSON_PRETTY" {
		return loglight.NewJsonLogFormatter(true).Format
	}

	return DebugFormatter
}

type debugOutput struct {
	LogLevel loglight.LogLevel
	Data     interface{}
}

func DebugFormatter(logEntry loglight.LogEntry) string {

	logEvent, _ := logEntry.Data.(*LogEvent)

	return fmt.Sprintf("%s : %s", logEntry.LogLevel, loglight.GetJson(logEvent.Additional, false))
}

type LogEvent struct {
	TimeUtc     time.Time
	ServiceKey  string
	Title       LogEventType
	Additional  interface{}
	RequestKey  string
	HostAddress string
}

func LogDebug(logType LogEventType, request *http.Request, additional interface{}) {
	log.Debug(newLogEvent(logType, request, additional))
}

func LogInfo(logType LogEventType, request *http.Request, additional interface{}) {
	log.Info(newLogEvent(logType, request, additional))
}

func LogError(logType LogEventType, request *http.Request, additional interface{}) {
	log.Error(newLogEvent(logType, request, additional))
}

func newLogEvent(logType LogEventType, request *http.Request, additional interface{}) *LogEvent {

	var requestKey string = ""

	if request != nil {
		requestKey = request.Header.Get(requestKeyHeaderName)
	}

	return &LogEvent{
		ServiceKey:  "SMSP_SVR",
		TimeUtc:     time.Now().UTC(),
		Title:       logType,
		RequestKey:  requestKey,
		HostAddress: hostAddress,
		Additional:  additional,
	}
}

type RequestKey struct {
	handler http.Handler
}

func (key *RequestKey) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(requestKeyHeaderName) == "" {
		r.Header.Set(requestKeyHeaderName, uuid.NewV4().String())
	}
}

// Handler apply the CORS specification on the request, and add relevant CORS headers
// as necessary.
func (c *RequestKey) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c.ServeHTTP(w, r)
		h.ServeHTTP(w, r)
	})
}
