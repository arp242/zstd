package ztime

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Takes records how long the execution of f takes.
func Takes(f func()) time.Duration {
	s := time.Now()
	f()
	return time.Since(s)
}

// TimeFunc prints how long it took for this function to end to stderr.
//
// You usually want to use this from defer:
//
//   defer ztime.TimeFunc()()
func TimeFunc() func() {
	s := time.Now()
	return func() { fmt.Fprintln(os.Stderr, time.Since(s)) }
}

// Sleep for d duration, or until the context times out.
func Sleep(ctx context.Context, d time.Duration) {
	if ctx == nil {
		time.Sleep(d)
		return
	}
	t := time.NewTimer(d)
	select {
	case <-t.C:
	case <-ctx.Done():
	}
}
