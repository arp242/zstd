// Package zexec runs external commands.
package zexec

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
)

type Output struct {
	Text   string // Read text, excluding the end delimiter.
	Stderr bool   // Text was sent on stderr.

	// Errors reading the process stdout/stderr; sent only once, and will stop
	// reading after an error.
	Error error
}

func (o Output) String() string {
	if o.Stderr {
		return "error: " + o.Text
	}
	return o.Text
}

// Readlines calls cmd.Start sends every line of output on the returned channel.
//
// Example usage:
//
//	cmd := exec.Command("long-running-process")
//	ch, err := zexec.Readlines(cmd)
//
//	for {
//		line, ok := <-ch
//		if !ok {
//			break
//		}
//		if line.Error != nil {
//			fmt.Fprintln(os.Stderr, "error reading:", line.Error)
//			break
//		}
//
//		fmt.Println(line)
//	}
func Readlines(cmd *exec.Cmd) (<-chan Output, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	ch := make(chan Output)
	go read(ch, stdout, false)
	go read(ch, stderr, true)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			ch <- Output{Error: err}
		}
		close(ch)
	}()

	return ch, nil
}

func read(ch chan<- Output, r io.Reader, stderr bool) {
	data := make([]byte, 0, 1024)
	for {
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			ch <- Output{Error: err}
			break
		}

		buf = buf[:n]
		for {
			i := bytes.IndexByte(buf, '\n')
			if i < 0 {
				if len(buf) > 0 {
					data = append(data, buf...)
				}
				break
			}

			line := buf[:i]
			if len(data) > 0 {
				line = append(data, line...)
				data = data[:0]
			}

			ch <- Output{Text: string(line), Stderr: stderr}
			buf = buf[i+1:]
		}
	}
}
