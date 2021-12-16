package zsync

import (
	"context"
	"reflect"
	"sort"
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
	n := atom.Add(2)
	if n != 668 {
		t.Errorf("wrong value: %v", n)
	}
	if v := atom.Value(); v != 668 {
		t.Errorf("wrong value: %v", v)
	}

	// For go test -race to ensure there are no data races.
	for i := 0; i < 10; i++ {
		go func(ii int) { atom.Set(int32(ii)) }(i)
		go func(ii int) { atom.Value() }(i)
	}
}

func TestAtomicInt64(t *testing.T) {
	atom := NewAtomicInt64(42)
	if v := atom.Value(); v != 42 {
		t.Errorf("wrong value: %v", v)
	}

	atom.Set(666)
	if v := atom.Value(); v != 666 {
		t.Errorf("wrong value: %v", v)
	}

	n := atom.Add(2)
	if n != 668 {
		t.Errorf("wrong value: %v", n)
	}
	if v := atom.Value(); v != 668 {
		t.Errorf("wrong value: %v", v)
	}

	// For go test -race to ensure there are no data races.
	for i := 0; i < 10; i++ {
		go func(ii int) { atom.Set(int64(ii)) }(i)
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

func TestAtMost(t *testing.T) {
	r := NewAtMost(3)

	var (
		list1, list2 []int
		mu           sync.Mutex
	)
	for i := 1; i <= 20; i++ {
		list1 = append(list1, i)
		func(i int) {
			r.Run(func() {
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				list2 = append(list2, i)
				mu.Unlock()
			})
		}(i)
	}
	r.Wait()

	sort.Ints(list1)
	sort.Ints(list2)

	if !reflect.DeepEqual(list1, list2) {
		t.Error()
	}
}
