package routes

import (
	"github.com/myles-mcdonnell/sendmailserviceproxy"
	"github.com/myles-mcdonnell/sendmailserviceproxy/logging"
	"gopkg.in/myles-mcdonnell/loglight.v3"
	"testing"
	"time"
)

type mockService struct {
	SendFunc sendmailserviceproxy.CallFunc
}

func NewMockService(sendFunc sendmailserviceproxy.CallFunc) mockService {
	return mockService{SendFunc: sendFunc}
}

func (service mockService) Send(msg sendmailserviceproxy.EmailMessage) error {
	return service.SendFunc(msg)
}

func init() {
	logging.Initialise(loglight.NewLogger(
		true,
		logging.BuildFormatter("DEBUG"),
	))
}

func TestWithTwoServicesBothOk(t *testing.T) {

	var emailHandler = NewEmailHandler(
		[]sendmailserviceproxy.EmailService{
			NewMockService(tWrap{T: t}.Ok),
			NewMockService(tWrap{T: t}.Ok),
		},
		"",
	)

	if emailHandler.Send(sendmailserviceproxy.EmailMessage{}) != nil {
		t.Fail()
	}
}

func TestWithTwoServicesFirstTimeout(t *testing.T) {

	var emailHandler = NewEmailHandler(
		[]sendmailserviceproxy.EmailService{
			NewMockService(tWrap{T: t}.ThreeSeconds),
			NewMockService(tWrap{T: t}.Ok),
		},
		"",
	)

	// This should open the first circuit
	if emailHandler.Send(sendmailserviceproxy.EmailMessage{}) == nil {
		t.Log("First call did not error")
		t.Fail()
	}

	select {
	case isClosed := <-emailHandler.circuits[0].IsClosedChangeChannel():
		if isClosed {
			t.Log("circuit expected to open but closed")
			t.Fail()
		}
	case <-time.After(time.Second * 3):
		t.Log("Timeout waiting for circuit to open")
		t.Fail()
	}

	if emailHandler.Send(sendmailserviceproxy.EmailMessage{}) != nil {
		t.Log("Second call did error")
		t.Fail()
	}
}

type tWrap struct {
	T *testing.T
}

func (t tWrap) Ok(msg sendmailserviceproxy.EmailMessage) error {
	t.T.Log("OK")
	return nil
}

func (t tWrap) ThreeSeconds(msg sendmailserviceproxy.EmailMessage) error {
	t.T.Log("ThreeSeconds")
	time.Sleep(time.Second * 10)
	return nil
}
