package zcontext

import (
	"context"
	"testing"
	"time"
)

func TestNoCancel(t *testing.T) {
	type (
		a struct{}
		b struct{}
	)

	ctx := context.WithValue(context.Background(), a{}, "value a")
	ctx, c := context.WithTimeout(ctx, 50*time.Millisecond)
	defer c()
	ctx = context.WithValue(ctx, b{}, "value b")

	ctx2 := WithoutTimeout(ctx)
	if ctx2.Value(a{}).(string) != "value a" {
		t.Fatal()
	}
	if ctx2.Value(b{}).(string) != "value b" {
		t.Fatal()
	}

	c()
	select {
	case <-ctx2.Done():
		t.Fatal()
	case <-time.After(100 * time.Millisecond):
	}
	if ctx2.Err() != nil {
		t.Fatal()
	}
	if ctx.Err() == nil {
		t.Fatal()
	}
}
