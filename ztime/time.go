package ztime

import (
	"context"
	"time"
)

// Takes records how long the execution of f takes.
func Takes(f func()) time.Duration {
	s := time.Now()
	f()
	return time.Since(s)
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
