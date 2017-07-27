package functional_tests

import (
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/myles-mcdonnell/sendmailserviceproxy/client"
	"github.com/myles-mcdonnell/sendmailserviceproxy/client/email"
	"github.com/myles-mcdonnell/sendmailserviceproxy/models"
	"os"
	"testing"
)

var authInfoWriter = runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {

	r.SetHeaderParam("X-API-KEY", os.Getenv("SMSP_API_KEY"))

	return nil
})

func TestSendEmailOk(t *testing.T) {

	httptrans := httptransport.New(os.Getenv("SMSP_HOST_PORT"), "/", []string{os.Getenv("SMSP_PROTOCOL")})
	apiclient := client.New(httptrans, nil)

	message := &models.Email{
		Toaddress:     SPtr("myles.mcdonnell@keyshift.co"),
		Fromaddress:   SPtr("nosuchaddress@keyshift.co"),
		Subject:       "test",
		Plaintextbody: "test_body",
	}

	params := email.NewPostEmailParams().WithEmail(message)

	_, err := apiclient.Email.PostEmail(params, authInfoWriter)

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestSendEmailTimeout(t *testing.T) {

	httptrans := httptransport.New(os.Getenv("SMSP_HOST_PORT"), "/", []string{"http"})
	apiclient := client.New(httptrans, nil)

	message := &models.Email{
		Toaddress:     SPtr("myles.mcdonnell@keyshift.co"),
		Fromaddress:   SPtr("nosuchaddress@keyshift.co"),
		Subject:       "timeout",
		Plaintextbody: "test_body",
	}

	params := email.NewPostEmailParams().WithEmail(message)

	_, err := apiclient.Email.PostEmail(params, authInfoWriter)

	if err == nil {
		t.Fail()
	}
}

func TestSendEmailFail(t *testing.T) {

	httptrans := httptransport.New(os.Getenv("SMSP_HOST_PORT"), "/", []string{"http"})
	apiclient := client.New(httptrans, nil)

	message := &models.Email{
		Toaddress:     SPtr("myles.mcdonnell@keyshift.co"),
		Fromaddress:   SPtr("nosuchaddress@keyshift.co"),
		Subject:       "should fail as all circuits open",
		Plaintextbody: "test_body",
	}

	params := email.NewPostEmailParams().WithEmail(message)

	_, err := apiclient.Email.PostEmail(params, authInfoWriter)

	if err == nil {
		t.Fail()
	}
}

func SPtr(s string) *string { return &s }
