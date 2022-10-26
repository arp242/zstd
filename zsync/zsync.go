// Package zsync adds functions for synchronization.
package zsync

import (
	"context"
	"sync"
	"sync/atomic"
)

// WithLock locks the passed mutex, runs the function, and unlocks.
//
//	WithLock(mu, func() {
//	    // .. stuff ..
//	})
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

func (i *AtomicInt) Set(value int32)   { atomic.StoreInt32((*int32)(i), value) }
func (i *AtomicInt) Add(n int32) int32 { return atomic.AddInt32((*int32)(i), n) }
func (i *AtomicInt) Value() int32      { return atomic.LoadInt32((*int32)(i)) }

// AtomicInt64 uses sync/atomic to store and read the value of an int64.
type AtomicInt64 int64

// NewAtomicInt creates an new AtomicInt.
func NewAtomicInt64(value int64) *AtomicInt64 {
	var i AtomicInt64
	i.Set(value)
	return &i
}

func (i *AtomicInt64) Set(value int64)   { atomic.StoreInt64((*int64)(i), value) }
func (i *AtomicInt64) Add(n int64) int64 { return atomic.AddInt64((*int64)(i), n) }
func (i *AtomicInt64) Value() int64      { return atomic.LoadInt64((*int64)(i)) }

// AtMost runs at most a certain number of goroutines in parallel.
type AtMost struct {
	ch chan struct{}
	wg sync.WaitGroup
}

// NewAtMost creates a new AtMost instance.
func NewAtMost(max int) AtMost {
	return AtMost{ch: make(chan struct{}, max)}
}

// Wait for all jobs to finish.
func (a *AtMost) Wait() {
	a.wg.Wait()
}

// Run a function. Blocks if the job queue is full.
func (a *AtMost) Run(f func()) {
	a.ch <- struct{}{}
	a.wg.Add(1)
	go func() {
		defer func() {
			<-a.ch
			a.wg.Done()
		}()
		f()
	}()
}
