package zsync

import (
	"bytes"
	"io"
	"sync"
)

// Buffer is a wrapper around bytes.Buffer which protects every operation with a
// lock, ensuring it can be read/write in a thread-safe manner.
type Buffer struct {
	buf *bytes.Buffer
	mu  *sync.Mutex
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{
		buf: bytes.NewBuffer(buf),
		mu:  new(sync.Mutex),
	}
}

func NewBufferString(s string) *Buffer {
	return &Buffer{
		buf: bytes.NewBufferString(s),
		mu:  new(sync.Mutex),
	}
}

func (b *Buffer) l() func() {
	b.mu.Lock()
	return b.mu.Unlock
}

func (b Buffer) Bytes() []byte {
	defer b.l()()
	return b.buf.Bytes()
}
func (b *Buffer) Cap() int {
	defer b.l()()
	return b.buf.Cap()
}
func (b *Buffer) Grow(n int) {
	defer b.l()()
	b.buf.Grow(n)
}
func (b *Buffer) Len() int {
	defer b.l()()
	return b.buf.Len()
}
func (b *Buffer) Next(n int) []byte {
	defer b.l()()
	return b.buf.Next(n)
}
func (b *Buffer) Read(p []byte) (int, error) {
	defer b.l()()
	return b.buf.Read(p)
}
func (b *Buffer) ReadByte() (byte, error) {
	defer b.l()()
	return b.buf.ReadByte()
}
func (b *Buffer) ReadBytes(delim byte) (line []byte, err error) {
	defer b.l()()
	return b.buf.ReadBytes(delim)
}
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	defer b.l()()
	return b.buf.ReadFrom(r)
}
func (b *Buffer) ReadRune() (r rune, size int, err error) {
	defer b.l()()
	return b.buf.ReadRune()
}
func (b *Buffer) ReadString(delim byte) (line string, err error) {
	defer b.l()()
	return b.buf.ReadString(delim)
}
func (b *Buffer) Reset() {
	defer b.l()()
	b.buf.Reset()
}
func (b Buffer) String() string {
	defer b.l()()
	return b.buf.String()
}
func (b *Buffer) Truncate(n int) {
	defer b.l()()
	b.buf.Truncate(n)
}
func (b *Buffer) UnreadByte() error {
	defer b.l()()
	return b.buf.UnreadByte()
}
func (b *Buffer) UnreadRune() error {
	defer b.l()()
	return b.buf.UnreadRune()
}
func (b *Buffer) Write(p []byte) (int, error) {
	defer b.l()()
	return b.buf.Write(p)
}
func (b *Buffer) WriteByte(c byte) error {
	defer b.l()()
	return b.buf.WriteByte(c)
}
func (b *Buffer) WriteRune(r rune) (n int, err error) {
	defer b.l()()
	return b.buf.WriteRune(r)
}
func (b *Buffer) WriteString(s string) (n int, err error) {
	defer b.l()()
	return b.buf.WriteString(s)
}
func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	defer b.l()()
	return b.buf.WriteTo(w)
}
