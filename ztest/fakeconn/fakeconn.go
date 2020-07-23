// Package fakeconn provides a "fake" net.Conn implementation.
package fakeconn

import (
	"bytes"
	"errors"
	"net"
	"time"
)

// Conn is a fake net.Conn implementations. Everything that is written to it
// with Write() is available in the Written buffer, and Read() reads from the
// data in the ReadFrom buffer.
type Conn struct {
	Written  *bytes.Buffer
	ReadFrom *bytes.Buffer
	closed   *bool
}

// New instance factory.
func New() Conn {
	var b bool
	return Conn{
		Written:  bytes.NewBuffer([]byte{}),
		ReadFrom: bytes.NewBuffer([]byte{}),
		closed:   &b,
	}
}

// Write data to the Written buffer.
func (c Conn) Write(b []byte) (n int, err error) {
	if *c.closed {
		return 0, errors.New("write to closed connection")
	}

	c.Written.Write(b)
	return len(b), nil
}

// Read data from the ReadFrom buffer.
func (c Conn) Read(b []byte) (n int, err error) {
	if *c.closed {
		return 0, errors.New("read from closed connection")
	}
	return c.ReadFrom.Read(b)
}

// Close clears the buffers and prevents further Read() and Write() operations.
func (c Conn) Close() error {
	*c.closed = true
	c.Written.Reset()
	c.ReadFrom.Reset()
	return nil
}

// LocalAddr does nothing.
func (c Conn) LocalAddr() net.Addr { return &net.TCPAddr{} }

// RemoteAddr does nothing.
func (c Conn) RemoteAddr() net.Addr { return &net.TCPAddr{} }

// SetDeadline does nothing.
func (c Conn) SetDeadline(t time.Time) error { return nil }

// SetReadDeadline does nothing.
func (c Conn) SetReadDeadline(t time.Time) error { return nil }

// SetWriteDeadline does nothing.
func (c Conn) SetWriteDeadline(t time.Time) error { return nil }
