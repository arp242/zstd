// Package zcontext implements context functions.
package zcontext

import (
	"context"
)

// WithoutTimeout returns a new context without any cancellations from
// WithTimeout() or WithDeadline(), but preserves any values.
//
// Deprecated: use context.WithoutCancel()
//
//go:fix inline
func WithoutTimeout(ctx context.Context) context.Context { return context.WithoutCancel(ctx) }
