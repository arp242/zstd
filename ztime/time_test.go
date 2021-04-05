package ztime

import (
	"context"
	"testing"
	"time"
)

func TestTakes(t *testing.T) {
	tt := Takes(func() { time.Sleep(50 * time.Millisecond) })
	if tt < 50*time.Millisecond || tt > 52*time.Millisecond {
		t.Error(tt)
	}

	func() {
		defer TimeFunc()()
		time.Sleep(50 * time.Millisecond)
	}()
}

func TestSleep(t *testing.T) {
	ctx := context.Background()
	Sleep(ctx, 20*time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()
	Sleep(ctx, 20*time.Hour)
}
