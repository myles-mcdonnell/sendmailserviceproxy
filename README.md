# Send Mail Service Proxy

[![Go Report Card](https://goreportcard.com/badge/github.com/myles-mcdonnell/sendmailserviceproxy)](https://goreportcard.com/report/github.com/myles-mcdonnell/sendmailserviceproxy)
[![GoDoc](https://godoc.org/github.com/myles-mcdonnell/sendmailserviceproxy?status.svg)](http://godoc.org/github.com/myles-mcdonnell/sendmailserviceproxy)


This is a JSON over HTTP service that proxies for 2 email services and will fail over between them as necessary. It fails fast and does not attempt to retry failed operations.

A bespoke implementation of [Circuit Breaker](https://martinfowler.com/bliki/CircuitBreaker.html) is used to detect downstream failure (reported errors and timeouts).  Real messages are never sent along circuits known to be open, rather the circuit is polled with a dummy recipient until such time as it closes.

For demonstration purposes there is a two node instance of the system deployed on [podspace.io](http://podspace.io) here [https://sendmailserviceproxy.keyshift.co](https://sendmailserviceproxy.keyshift.co/healthcheck)

To send an email execute the following

`curl -i -H "X-API-KEY: {insert key here}" -H "Content-Type: application/json" -X POST -d '{"fromaddress":"an@email.com","plaintextbody":"test_body","subject":"test","toaddress":"another@email.com"}' https://sendmailserviceproxy.keyshift.co/email`


## Tech. Stack

* The service is written in Go using [Go-Swagger](https://github.com/go-swagger/go-swagger) to generate the server skeleton from a [Swagger Spec](spec).
* json-refs (a Nodejs package) is used to aggregate the spec to a single JSON file from the many YML files that define the models, routes etc
* Code is hosted on GitHub
* Go dependencies are managed with [Glide](http://github.com/masterminds/glide)
* Unit and [functional tests](functionaltests) both use the Go testing framework.
* The build output is an immutable docker image which is published to [Dockerhub](https://hub.docker.com/r/mylesmcdonnell/sendmailserviceproxy/).
* [12 Factor App](https://12factor.net/) principles are applied.  All configuration is read from the environment, none exists in source control or the build output.
* The service is currently hosted on [Podspace.io](https://www.podspace.io/) (which is Kubernetes and OpenShift aaS) running 2 load balanced nodes.
* Logging uses http://github.com/myles-mcdonnell/loglight to create structured log output to stdout which could then be aggregated across multiple nodes using a number techniques, e.g. https://byteshuffle.net/2016/12/02/microservices-log-aggregation/

## Scalability

Maximum throughput and scalability are design goals for this service.

With regard to throughput no mutually exclusive locks are taken when mail is sent.  Running the unit tests with race detection enabled will evidence a race condition; this is by design.  The effect is latency around the behaviour according to the state of the circuit, this is preferable to the reduction in throughput incurred were this race eliminated through synchronization.

This service can be scaled out over *n* nodes.  The only state within the service is that of the circuits.  It is desirable to hold the state in isolation as circuits may be open from some nodes and closed from others according to network conditions etc.

## Security

* An API Key set as header X-API-KEY is required to send email.
* The healthcheck endpoint (https://sendmailserviceproxy.keyshift.co/healthcheck) is unsecured.

## API Documentation

Browse to [http://petstore.swagger.io/](http://petstore.swagger.io/) (or pull and run the [Swagger UI](https://swagger.io/swagger-ui/) locally if you prefer) then enter [https://sendmailserviceproxy.keyshift.co/swagger.json](https://sendmailserviceproxy.keyshift.co/swagger.json) and click explore to see the interactive documentation


## Build Instructions

* Create a new go workspace and clone repo
* Install recent version of NodeJs if not already present
* Run `chmod +x ./scripts/`
* Run [scripts/install_build_tools.sh](install_build_tools.sh)
* Run [scripts/swagger_code_gen.sh](swagger_code_gen.sh) to generate the server skeleton and [client](client) package.
* Run `glide install` to pull Go dependencies (see here for Glide info http://github.com/masterminds/glide) (this will terminate in error due to version conflict, this can be ignored)
* Run [scripts/run_unit_tests.sh](scripts/run_unit_tests.sh)
* Set the following environment variables
    * SMSP_MG_DOMAIN - MailGun Domain
    * SMSP_MG_API_KEY - MailGun API Key
    * SMSP_SG_API_KEY - SendGrid API Key
    * SMSP_LOG_OUTPUT_FORMAT - JSON_ONLINE | JSON_PRETTY | DEBUG (optional - default DEBUG)
    * SMSP_LOG_OUTPUT_DEBUG - true/false (optional - default false)
    * SMSP_MOCK_EMAIL_SERVICE - true/false (optional) if true then the MailGun and SendGrid service will not be used.  This is used for [functional testing](functional_tests).
    * SMSP_POLL_MESSAGE_ADDRESS - address used for dummy messages sent along open circuits
    * SMSP_API_KEY - the API key used by clients to send email
* Run [scripts/build_and_run_server.sh](scripts/build_and_run_server.sh)

To run the function tests:

* Set the follwoing environment variables:
    * SMSP_HOST_PORT = e.g. localhost:8080
    * SMSP_PROTOCOL = http | https
    * SMSP_API_KEY = same value used for server environment
* Run `go test ./functional_tests`


## AOB

* The API caters only for single recipient, plain text emails with no attachments sent individually. In production environment this would of course be extended to cater for the full suite of SMTP functionality.
* The implementation of the API-KEY and authentication middleware is not suitable for a production system and was done this way simply to save time.
* The circuit breaker implementation opens a circuit on a first failure.  There are many different ways to determine an open circuit, such as consecutive failure, failure rate etc.  I choose not to spend time working on this for the demo and went for a naive implementation.
* To load test the system I would probably use [http://locust.io](http://locust.io)
* The system is hardcoded to use SendMail and MailGun services.  However the [email handler](routes/email.go) is coded to work with *n* email services so it's easy to see how the service could be extended to enable multiple underlying email services and with Go plugins these could be injected dynamically at startup.
* There could be a more comprehensive set of functional and unit tests.  What has been implemented proves the model but I would certainly invest more time here for a production system.
* In a production scenario I would use a build platform, such as CirciCI, to automate the build pipeline including unit and functional testing and the publishing of docker images on successful build.  A continuous delivery/deployment pipeline could then be implemented onwards if that where a requirement.


