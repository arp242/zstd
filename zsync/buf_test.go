package zsync

import (
	"fmt"
	"io"
	"testing"
	"time"
)

var _ io.ReadWriter = &Buffer{}

func TestBuffer(t *testing.T) {
	buf := NewBuffer(nil)

	go buf.Write([]byte("one "))
	go fmt.Println(buf.String())
	go buf.Write([]byte("two "))
	go fmt.Println(buf.String())
	go buf.Write([]byte("three "))
	go fmt.Println(buf.String())
	go buf.Write([]byte("four "))
	go fmt.Println(buf.String())
	go buf.Write([]byte("five "))
	go fmt.Println(buf.String())
	go buf.Write([]byte("six "))
	go fmt.Println(buf.String())

	time.Sleep(50 * time.Millisecond)
	fmt.Println(buf.String())
}
