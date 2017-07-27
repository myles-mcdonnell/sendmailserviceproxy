package sendmailserviceproxy

import (
	"github.com/pkg/errors"
	"testing"
	"time"
)

type mockFunc struct {
	failAfterReqs   int
	resumeAfterReqs int
	reqCount        int
	failFunc        CallFunc
}

func (mockFunc *mockFunc) Call(message EmailMessage) error {
	//this only works for sequential access of course which is fine here as all of these tests do just that
	mockFunc.reqCount++

	if mockFunc.reqCount > mockFunc.failAfterReqs && mockFunc.reqCount <= mockFunc.resumeAfterReqs {
		return mockFunc.failFunc(message)
	}

	return nil
}

func TestClosedErrorClosed(t *testing.T) {

	testClosedOpenClosed(t, func(message EmailMessage) error { return errors.New("kaboom!") })
}

func TestClosedTimeoutClosed(t *testing.T) {

	testClosedOpenClosed(t, func(message EmailMessage) error {
		time.Sleep(time.Millisecond * 300)
		return nil
	})
}

func testClosedOpenClosed(t *testing.T, callFunc CallFunc) {

	var mockFunc = &mockFunc{
		failAfterReqs:   2,
		resumeAfterReqs: 4,
		failFunc:        callFunc,
	}

	var circuit = NewCircuit(
		mockFunc.Call,
		EmailMessage{},
	).WithPollInterval(time.Millisecond * 250).WithTimeout(time.Millisecond * 250)

	done := make(chan int)
	go func() {
		hasOpened := false
		for {
			select {
			case isClosed := <-circuit.IsClosedChangeChannel():
				if !isClosed {
					hasOpened = true
				}
				if isClosed && hasOpened {
					done <- 0
					break
				}
			}
		}
	}()

	go func() {
		for i := 0; i < 3; i++ {
			circuit.Call(EmailMessage{})
		}
	}()

	select {
	case <-done:
	case <-time.After(time.Second * 3):
		t.Fail()
	}
}
