// Package syncutil adds functions for synchronization.
package syncutil

import (
	"context"
	"sync"
	"sync/atomic"
)

// WithLock locks the passed mutex, runs the function, and unlocks.
//
//   WithLock(mu, func() {
//       // .. stuff ..
//   })
//
// This is convenient especially in cases where you don't want to defer the
// Unlock(), but also want to ensure the Unlock() is always called, regardless
// of runtime errors.
func WithLock(mu *sync.Mutex, f func()) {
	mu.Lock()
	defer mu.Unlock()
	f()
}

// Wait for a sync.WaitGroup with support for timeout/cancellations from
// context.
func Wait(ctx context.Context, wg *sync.WaitGroup) error {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ch:
		return nil
	}
}

// AtomicInt uses sync/atomic to store and read the value of an int32.
type AtomicInt int32

// NewAtomicInt creates an new AtomicInt.
func NewAtomicInt(value int32) *AtomicInt {
	var i AtomicInt
	i.Set(value)
	return &i
}

func (i *AtomicInt) Set(value int32) { atomic.StoreInt32((*int32)(i), value) }
func (i *AtomicInt) Value() int32    { return atomic.LoadInt32((*int32)(i)) }
