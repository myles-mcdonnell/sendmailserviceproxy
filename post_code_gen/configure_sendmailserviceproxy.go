package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	graceful "github.com/tylerb/graceful"

	"github.com/myles-mcdonnell/sendmailserviceproxy/logging"
	"github.com/myles-mcdonnell/sendmailserviceproxy/restapi/operations"
	"github.com/rs/cors"
	"os"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target .. --name taskpilot_api --spec ../swagger.json

func configureFlags(api *operations.SendmailserviceproxyAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

var apiKey string

func configureAPI(api *operations.SendmailserviceproxyAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	apiKey = os.Getenv("SMSP_API_KEY")
	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	InitApi(api)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {

	r := &logging.RequestKey{}
	handler = r.Handler(handler)

	//CORS enabled for https://swagger.io/swagger-ui/
	c := cors.New(cors.Options{AllowedOrigins: []string{"*"}})
	handler = c.Handler(handler)

	return authenticate(handler)
}

func authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//This is obviously not a viable piece of code in a a production system
		if r.URL.Path != "/healthcheck" && r.URL.Path != "/swagger.json" && r.Header.Get("X-API-KEY") != apiKey {
			logging.LogInfo(logging.Unauthorised, r, "X-API-KEY : "+r.Header.Get("X-API-KEY"))
			http.Error(w, "Unauthorised", 401)
			return
		}

		h.ServeHTTP(w, r)
	})
}
