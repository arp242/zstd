package fakeconn

import (
	"net"
	"reflect"
	"testing"
)

var _ net.Conn = Conn{}

func TestWrite(t *testing.T) {
	conn := New()
	n, err := conn.Write([]byte("Hello"))
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}
	if n != 5 {
		t.Errorf("n != 5: %v", n)
	}
	if !reflect.DeepEqual(conn.Written.Bytes(), []byte("Hello")) {
		t.Errorf("Written wrong: %v", conn.Written)
	}

	n, err = conn.Write([]byte(", world"))
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}
	if n != 7 {
		t.Errorf("n != 7: %v", n)
	}
	if !reflect.DeepEqual(conn.Written.Bytes(), []byte("Hello, world")) {
		t.Errorf("Written wrong: %v", conn.Written)
	}
}

func TestRead(t *testing.T) {
	conn := New()
	conn.ReadFrom.WriteString("read me")

	b := make([]byte, 7)
	n, err := conn.Read(b)
	if err != nil {
		t.Errorf("err != nil: %v", err)
	}
	if n != 7 {
		t.Errorf("n != 7: %v", n)
	}
	if !reflect.DeepEqual(b, []byte("read me")) {
		t.Errorf("Read wrong: %v", b)
	}
}

func TestSame(t *testing.T) {
	conn := New()
	conn.ReadFrom = conn.Written

	_, err := conn.Write([]byte("Hello"))
	if err != nil {
		t.Fatal(err)
	}

	b := make([]byte, 5)
	_, err = conn.Read(b)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(b, []byte("Hello")) {
		t.Errorf("Read wrong: %v", b)
	}
}

func TestClose(t *testing.T) {
	conn := New()
	err := conn.Close()
	if err != nil {
		t.Fatal(err)
	}

	_, err = conn.Write([]byte("a"))
	if err == nil {
		t.Errorf("err is nil for Write")
	}

	_, err = conn.Read(nil)
	if err == nil {
		t.Errorf("err is nil for Read")
	}
}
