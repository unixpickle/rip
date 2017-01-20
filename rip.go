// +build !js

package rip

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

// A RIP listens for the next interrupt signal.
//
// When an interrupt is received, a message is printed out
// to the user indicating that the interrupt was caught.
// The interrupt handler is then deregistered, meaning any
// further interrupts will terminate the program.
type RIP struct {
	cancelLock sync.Mutex
	cancel     chan struct{}

	killChan chan struct{}
}

// NewRIP starts a RIP.
func NewRIP() *RIP {
	cancel := make(chan struct{})
	kill := make(chan struct{})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case <-c:
			fmt.Println("\nCaught interrupt. Ctrl+C again to terminate.")
		case <-cancel:
		}
		close(kill)
		signal.Stop(c)
	}()

	return &RIP{cancel: cancel, killChan: kill}
}

// Done returns true if an interrupt was received or if
// Close was called.
func (r *RIP) Done() bool {
	select {
	case <-r.killChan:
		return true
	default:
		return false
	}
}

// Chan returns a channel which is closed when the RIP
// receives an interrupt or when Close is called.
func (r *RIP) Chan() chan struct{} {
	return r.killChan
}

// Close stops listening for interrupts.
func (r *RIP) Close() error {
	var err error
	r.cancelLock.Lock()
	if r.cancel != nil {
		close(r.cancel)
		r.cancel = nil
	} else {
		err = errors.New("already closed")
	}
	r.cancelLock.Unlock()
	return err
}
