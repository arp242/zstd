package zsync

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestWait(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := Wait(ctx, &wg)
		if err != context.Canceled {
			t.Errorf("wrong error: %v", err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := Wait(ctx, &wg)
		if err != context.DeadlineExceeded {
			t.Errorf("wrong error: %v", err)
		}
	})

	t.Run("finish", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		wg.Done()
		wg.Done()

		err := Wait(ctx, &wg)
		if err != nil {
			t.Errorf("wrong error: %v", err)
		}
	})
}

func TestAtomicInt(t *testing.T) {
	atom := NewAtomicInt(42)
	if v := atom.Value(); v != 42 {
		t.Errorf("wrong value: %v", v)
	}

	atom.Set(666)
	if v := atom.Value(); v != 666 {
		t.Errorf("wrong value: %v", v)
	}

	// For go test -race to ensure there are no data races.
	for i := 0; i < 10; i++ {
		go func(ii int) { atom.Set(int32(ii)) }(i)
		go func(ii int) { atom.Value() }(i)
	}
}

func TestWithLock(t *testing.T) {
	mu := new(sync.Mutex)

	// Be lazy and rely on the race detector to report failures.
	var s int
	go WithLock(mu, func() { s = 1 })
	mu.Lock()
	_ = s
	mu.Unlock()
	time.Sleep(10 * time.Millisecond)
}
