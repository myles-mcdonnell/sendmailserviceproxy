package sendmailserviceproxy

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// CallResult CIRCUIT_OPEN | TIMEOUT | CALL_COMPLETE
type CallResult int

const (
	CIRCUIT_OPEN CallResult = iota
	TIMEOUT
	CALL_COMPLETE
)

//CallFunc is the function being wrapped by a circuit breaker
type CallFunc func(message EmailMessage) error

type Circuit struct {
	isClosed              bool
	callFunc              CallFunc
	timeout               time.Duration
	pollInterval          time.Duration
	pollMessage           EmailMessage
	isClosedChangeChannel chan bool
	stateLock             sync.Mutex
}

//NewCircuit ctor
func NewCircuit(callfunc CallFunc, pollMessage EmailMessage) *Circuit {
	return &Circuit{
		isClosed:              true,
		timeout:               2 * time.Second,
		pollInterval:          5 * time.Second,
		pollMessage:           pollMessage,
		callFunc:              callfunc,
		isClosedChangeChannel: make(chan bool),
		stateLock:             sync.Mutex{},
	}
}

//IsClosed is true if the circuit is not open
func (circuit *Circuit) IsClosed() bool {
	return circuit.isClosed
}

// IsClosedChangeChannel: the state of the circuit is cent on this channel each time it changes
func (circuit *Circuit) IsClosedChangeChannel() chan bool {
	return circuit.isClosedChangeChannel
}

//WithTimeout enables the timeout to be changed
func (circuit *Circuit) WithTimeout(timeout time.Duration) *Circuit {
	circuit.timeout = timeout
	return circuit
}

//WithPollInterval enables the poll interval to be changed
func (circuit *Circuit) WithPollInterval(pollInterval time.Duration) *Circuit {
	circuit.pollInterval = pollInterval
	return circuit
}

//Call is the outer call of the underlying CallFunc and is where the circuit break login exists
func (circuit *Circuit) Call(message EmailMessage) (error, CallResult) {

	if !circuit.isClosed {
		fmt.Println("circuit.Call : cirxcuit open - fail fast")
		return errors.New("Circuit is open"), CIRCUIT_OPEN
	}

	c := make(chan error)
	go func() { c <- circuit.callFunc(message) }()

	var err error
	var callResult CallResult
	select {
	case err = <-c:
		callResult = CALL_COMPLETE
	case <-time.After(circuit.timeout):
		err, callResult = errors.New("Call timeout"), TIMEOUT
	}

	if err != nil {
		go func() { circuit.openCircuit() }()
	}

	return err, callResult

}

func (circuit *Circuit) openCircuit() {
	defer circuit.stateLock.Unlock()

	circuit.stateLock.Lock()

	if !circuit.isClosed {
		return
	}

	fmt.Println("circuit.Open")
	circuit.isClosed = false
	circuit.isClosedChangeChannel <- circuit.isClosed
	ticker := time.NewTicker(circuit.pollInterval)

	for range ticker.C {

		c := make(chan error)
		go func() {

			c <- circuit.callFunc(circuit.pollMessage)
		}()

		var err error
		select {
		case err = <-c:
		case <-time.After(circuit.timeout):
			err = errors.New("Call timeout")
		}

		if err == nil {
			break
		}
	}

	ticker.Stop()
	fmt.Println("circuit.Closed")
	circuit.isClosed = true
	circuit.isClosedChangeChannel <- circuit.isClosed
}
