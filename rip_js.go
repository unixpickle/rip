// +build js

package rip

import (
	"errors"
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

	go func() {
		<-cancel
		close(kill)
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
