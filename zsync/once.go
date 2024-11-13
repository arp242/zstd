package zsync

import (
	"sync"
)

// Once is an object that will perform exactly one action per key.
//
// This is similar to sync.Once, but allows grouping per-key and has a return
// value informing whether the function was already run.
//
// This implementation is a bit slower than sync.Once; the benchmark regresses
// from ~1.6ns/op to ~52ns/op on my system.
type Once struct {
	m    sync.Mutex
	done map[string]struct{}
}

// Do calls the function f for the given key only on the first invocation.
//
// In other words, given:
//
//	var once Once
//
// If once.Do("x", f) is called multiple times, only the first call will invoke
// f, even if f has a different value in each invocation. A new key or instance
// of Once is required for each function to execute.
//
// The return value tells you if f is run; it's true on the first caller, and
// false on all subsequent calls.
//
// It may be necessary to use a function literal to capture the arguments to a
// function to be invoked by Do:
//
//	once.Do("x", func() { config.init(filename) })
//
// Because no call to Do returns until the one call to f returns, if f causes Do
// to be called, it will deadlock.
//
// If f panics, Do considers it to have returned; future calls of Do return
// without calling f.
func (o *Once) Do(key string, f func()) bool {
	o.m.Lock()
	defer o.m.Unlock()

	if o.done == nil {
		o.done = make(map[string]struct{})
	}
	_, ok := o.done[key]
	if ok {
		return false
	}

	defer func() { o.done[key] = struct{}{} }()
	f()
	return true
}

// Did reports if something has been run for the given key.
func (o *Once) Did(key string) bool {
	o.m.Lock()
	_, ok := o.done[key]
	o.m.Unlock()
	return ok
}

// Forget about a key, causing the next invocation to Do() to run again.
func (o *Once) Forget(key string) {
	o.m.Lock()
	delete(o.done, key)
	o.m.Unlock()
}
