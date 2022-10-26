// Package zruntime provides utilities to interface with the Go runtime.
package zruntime

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"zgo.at/zstd/zstring"
)

// Test reports if we're running a go test command.
func Test() bool {
	for _, a := range os.Args[1:] {
		if strings.HasPrefix(a, "-test.") {
			return true
		}
	}
	return false
}

// TestVerbose reports if the test was started with the -v flag.
func TestVerbose() bool {
	for _, a := range os.Args[1:] {
		if a == "-test.v=true" {
			return true
		}
	}
	return false
}

// FuncName gets the name of a function.
func FuncName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

// GoroutineID gets the current goroutine ID.
//
// Go doesn't give access to this to discourage certain unwise design patterns,
// but in some cases it can still be useful; for example for some tests or
// debugging.
func GoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))

	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("could not parse %q", b))
	}

	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("could not parse %q: %v", b, err))

	}
	return n
}

// Callers gets a list of callers.
func Callers(filterFun ...string) []runtime.Frame {
	var (
		pc     = make([]uintptr, 50)
		n      = runtime.Callers(2, pc)
		frames = runtime.CallersFrames(pc[:n])
	)

	filterFun = append(filterFun, []string{
		"runtime.goexit",
		"runtime.gopanic",
		"runtime.panicdottypeE",
		"runtime.goPanicIndex",
		"runtime.main",
		"testing.tRunner",
		"zgo.at/zstd/zdebug.Stack",
		"zgo.at/zstd/zdebug.PrintStack",
	}...)

	ret := make([]runtime.Frame, 0, n)
	for f, more := frames.Next(); more; f, more = frames.Next() {
		if zstring.HasPrefixes(f.Function, filterFun...) {
			continue
		}
		ret = append(ret, f)
	}
	if len(ret) == 0 && len(filterFun) > 0 { // Everything was filtered: re-run without filter.
		return Callers()
	}
	return ret
}

// SizeOf gets the memory size of an object in bytes.
//
// This recurses struct fields and pointers, but there are a few limitations:
//
//  1. Space occupied by code and data reachable through variables captured in
//     the closure of a function pointer are not counted. A value of function
//     type is counted only as a pointer.
//
//  2. Unused buckets of a map cannot be inspected by the reflect package. Their
//     size is estimated by assuming unfilled slots contain zeroes of their type.
//
//  3. Unused capacity of the array underlying a slice is estimated by assuming
//     the unused slots contain zeroes of their type. It is possible they contain
//     non zero values from sharing or reslicing, but without explicitly
//     reslicing the reflect package cannot touch them.
//
// This is adapted from: https://github.com/creachadair/misctools/blob/master/sizeof/size.go
//
// Also see: https://github.com/golang/go/issues/34561
func SizeOf(v any) int64 {
	return int64(sizeOf(reflect.ValueOf(v), make(map[uintptr]bool)))

	// Copyright (c) 2016, Michael J. Fromberger
	// All rights reserved.
	//
	// Redistribution and use in source and binary forms, with or without
	// modification, are permitted provided that the following conditions are met:
	//
	// 1. Redistributions of source code must retain the above copyright notice, this
	//    list of conditions and the following disclaimer.
	//
	// 2. Redistributions in binary form must reproduce the above copyright notice,
	//    this list of conditions and the following disclaimer in the documentation
	//    and/or other materials provided with the distribution.
	//
	// 3. Neither the name of the copyright holder nor the names of its contributors
	//    may be used to endorse or promote products derived from this software
	//    without specific prior written permission.
	//
	// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
	// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
	// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
	// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
	// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
	// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
	// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
	// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
	// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
}

func sizeOf(v reflect.Value, seen map[uintptr]bool) uintptr {
	size := v.Type().Size()

	switch v.Kind() {
	case reflect.Ptr:
		p := v.Pointer()
		if !seen[p] && !v.IsNil() {
			seen[p] = true
			return size + sizeOf(v.Elem(), seen)
		}

	case reflect.String:
		return size + uintptr(v.Len())

	case reflect.Slice:
		n := v.Len()
		for i := 0; i < n; i++ {
			size += sizeOf(v.Index(i), seen)
		}

		// Account for the parts of the array not covered by this slice.  Since
		// we can't get the values directly, assume they're zeroes. That may be
		// incorrect, in which case we may underestimate.
		if cap := v.Cap(); cap > n {
			size += v.Type().Size() * uintptr(cap-n)
		}

	case reflect.Map:
		// A map m has len(m) / 6.5 buckets, rounded up to a power of two, and
		// a minimum of one bucket. Each bucket is 16 bytes + 8*(keysize + valsize).
		//
		// We can't tell which keys are in which bucket by reflection, however,
		// so here we count the 16-byte header for each bucket, and then just add
		// in the computed key and value sizes.
		nb := uintptr(math.Pow(2, math.Ceil(math.Log(float64(v.Len())/6.5)/math.Log(2))))
		if nb == 0 {
			nb = 1
		}
		size = 16 * nb
		for _, key := range v.MapKeys() {
			size += sizeOf(key, seen)
			size += sizeOf(v.MapIndex(key), seen)
		}

		// We have nb buckets of 8 slots each, and v.Len() slots are filled.
		// The remaining slots we will assume contain zero key/value pairs.
		zk := v.Type().Key().Size()  // a zero key
		zv := v.Type().Elem().Size() // a zero value
		size += (8*nb - uintptr(v.Len())) * (zk + zv)

	case reflect.Struct:
		// Chase pointer and slice fields and add the size of their members.
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			switch f.Kind() {
			case reflect.Ptr:
				p := f.Pointer()
				if !seen[p] && !f.IsNil() {
					seen[p] = true
					size += sizeOf(f.Elem(), seen)
				}
			case reflect.Slice:
				size += sizeOf(f, seen)
			}
		}
	}

	return size
}
