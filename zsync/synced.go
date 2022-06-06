package zsync

import (
	"fmt"
	"sync"
)

// Synced provides a thread-safe synced variable.
type Synced[T any] struct {
	mu  *sync.RWMutex
	val T
}

// NewSynced creates a new Synced with the initial value set to val.
func NewSynced[T any](val T) Synced[T] {
	return Synced[T]{
		mu:  new(sync.RWMutex),
		val: val,
	}
}

func (s Synced[T]) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return fmt.Sprintf("%s", any(s.val))
}

// Get the value.
func (s *Synced[T]) Get() T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.val
}

// Set the value.
func (s *Synced[T]) Set(to T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.val = to
}
