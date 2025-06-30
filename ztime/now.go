package ztime

import (
	"context"
	"time"
)

var ctxkey = &struct{ n string }{"now"}

// Now returns the time on the context, if any.
//
// Returns just time.Now().UTC() if there is no time on the context.
func Now(ctx context.Context) time.Time {
	t, ok := ctx.Value(ctxkey).(time.Time)
	if ok {
		return t.UTC()
	}
	return time.Now().UTC()
}

// WithNow returns a context with the time set.
func WithNow(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, ctxkey, t)
}
