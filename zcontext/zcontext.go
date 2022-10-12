// Package zcontext implements context functions.
package zcontext

import (
	"context"
	"time"
)

type ntContext struct{ p context.Context }

func (ntContext) Deadline() (time.Time, bool) { return time.Time{}, false }
func (ntContext) Done() <-chan struct{}       { return nil }
func (ntContext) Err() error                  { return nil }
func (c ntContext) Value(key any) any         { return c.p.Value(key) }

// WithoutTimeout returns a new context without any cancellations from
// WithTimeout() or WithDeadline(), but preserves any values.
func WithoutTimeout(ctx context.Context) context.Context { return ntContext{ctx} }
