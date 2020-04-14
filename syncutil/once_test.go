// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syncutil_test

import (
	"testing"
	"time"

	. "zgo.at/utils/syncutil"
)

type one int

func (o *one) Increment() { *o++ }

func run(t *testing.T, once *Once, key string, o *one, want bool, c chan bool) {
	got := once.Do(key, func() { o.Increment() })
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	if want != got {
		t.Errorf("wrong return: %t", got)
	}
	c <- true
}

func TestOnce(t *testing.T) {
	o := new(one)
	o2 := new(one)
	once := new(Once)
	c := make(chan bool)
	const N = 10

	go run(t, once, "x", o, true, c)
	time.Sleep(25 * time.Millisecond)

	for i := 0; i < N; i++ {
		go run(t, once, "x", o, false, c)
	}
	go run(t, once, "y", o2, true, c)
	for i := 0; i < N+2; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}

func TestOncePanic(t *testing.T) {
	var once Once
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("Once.Do did not panic")
			}
		}()
		once.Do("x", func() {
			panic("failed")
		})
	}()

	once.Do("x", func() {
		t.Fatalf("Once.Do called twice")
	})
}

func BenchmarkOnce(b *testing.B) {
	var once Once
	f := func() {}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.Do("x", f)
		}
	})
}
