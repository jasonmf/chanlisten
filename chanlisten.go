/* chanlisten implements net.Listener with a channel of net.Conn. */
package chanlisten

import (
	"io"
	"net"
)

type fakeAddr string

func (f fakeAddr) Network() string {
	return string(f)
}
func (f fakeAddr) String() string {
	return string(f)
}

// ChanListener implements net.Listener internally as a chan net.Conn and a
// channel for communicating that the channel is no longer available for
// Accept or Queue.
type ChanListener struct {
	c    chan net.Conn
	stop chan bool
}

// Create a new ChanListener with the specified buffer length. 0 makes it
// synchronous.
func New(length int) *ChanListener {
	return &ChanListener{
		c:    make(chan net.Conn, length),
		stop: make(chan bool),
	}
}

// Accept a net.Conn. Blocks if none are ready, returns io.EOF if closed.
func (c *ChanListener) Accept() (net.Conn, error) {
	select {
	case conn := <-c.c:
		return conn, nil
	case <-c.stop:
		return nil, io.EOF
	}
	panic("Accept should never hit this")
}

// Closes the Listener, prevening new net.Conn from being added or Accepted.
// The returned error is always nil.
func (c *ChanListener) Close() error {
	close(c.stop)
	return nil
}

// Return the address of the listener whic his always "channel".
func (c *ChanListener) Addr() net.Addr {
	return fakeAddr("channel")
}

// Send a net.Conn to Accept. If the listener is closed, returns io.EOF.
func (c *ChanListener) Queue(conn net.Conn) error {
	select {
	case c.c <- conn:
		return nil
	case <-c.stop:
		return io.EOF
	}
	panic("Queue should never hit this")
}
