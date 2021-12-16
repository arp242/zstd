package ztime

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"zgo.at/zstd/zruntime"
)

func TestDurationsGrow(t *testing.T) {
	d := NewDurations(0)

	test := func(wantCap, wantLen int, vals ...time.Duration) {
		t.Helper()
		if cap(d.list) != wantCap || len(d.list) != wantLen {
			t.Errorf("wrong size:\nhave: cap=%d; len=%d\nwant: cap=%d; len=%d",
				cap(d.list), len(d.list), wantCap, wantLen)
		}
		if fmt.Sprintf("%s", d.list) != fmt.Sprintf("%s", vals) {
			t.Errorf("wrong values:\nhave: %s\nwant: %s", d.list, vals)
		}
	}

	test(0, 0)
	d.Grow(8)
	test(8, 0)

	d.Append(123)
	d.Append(456)
	test(8, 2, 123, 456)

	d.Grow(2)
	test(8, 2, 123, 456)

	d.Grow(8)
	test(10, 2, 123, 456)

	d.Append(999)
	test(10, 3, 123, 456, 999)

	d.Append(1, 2, 3, 4, 5, 6, 7, 8)
	test(20, 11, 123, 456, 999, 1, 2, 3, 4, 5, 6, 7, 8)
}

func cmpDurations(t *testing.T, have, want Durations) {
	t.Helper()

	if !reflect.DeepEqual(have.list, want.list) {
		t.Errorf("list wrong\nhave: %v\nwant: %v", have.list, want.list)
	}
	if len(have.data) > 0 && len(want.data) > 0 && !reflect.DeepEqual(have.data, want.data) {
		t.Errorf("data wrong\nhave: %v\nwant: %v", have.data, want.data)
	}
}

func TestDurationsAppend(t *testing.T) {
	d := NewDurations(4)

	d.Append(1)
	d.appendWithData(2, "two")
	d.Append(3)
	d.appendWithData(4, "four")
	cmpDurations(t, d, Durations{
		list: []time.Duration{1, 2, 3, 4},
		data: map[int]interface{}{
			1: "two",
			3: "four",
		},
	})

	d.appendWithData(5, "five")
	d.appendWithData(6, "six")
	cmpDurations(t, d, Durations{
		list: []time.Duration{3, 4, 5, 6},
		data: map[int]interface{}{
			3: "four",
			4: "five",
			5: "six",
		},
	})
}

func TestDurationsTop(t *testing.T) {
	d := NewDurations(0)
	d.Append(1)
	d.appendWithData(4, "four")
	d.Append(3)
	d.appendWithData(2, "two")
	d.Append(999)
	d.Append(42)
	d.Append(666)
	d.Append(123)

	cmpDurations(t, d.Top(50), Durations{
		list: []time.Duration{999, 42, 666, 123},
	})
	cmpDurations(t, d.Top(20), Durations{
		list: []time.Duration{999},
	})

	cmpDurations(t, d.Top(-50), Durations{
		list: []time.Duration{1, 4, 3, 2},
		data: map[int]interface{}{1: "four", 3: "two"},
	})
	cmpDurations(t, d.Top(-20), Durations{list: []time.Duration{1}})
}

func TestDurationsDistrubute(t *testing.T) {
	d := NewDurations(0)
	d.Append(20)
	d.Append(20)
	d.Append(40)
	d.Append(40)
	d.Append(60)
	d.Append(60)
	d.Append(80)
	d.Append(80)

	{
		h := d.Distrubute(4)
		cmpDurations(t, h[0], Durations{list: []time.Duration{20, 20}})
		cmpDurations(t, h[1], Durations{list: []time.Duration{40, 40}})
		cmpDurations(t, h[2], Durations{list: []time.Duration{60, 60}})
		cmpDurations(t, h[3], Durations{list: []time.Duration{80, 80}})
	}

	{
		h := d.Distrubute(2)
		cmpDurations(t, h[0], Durations{list: []time.Duration{20, 20, 40, 40}})
		cmpDurations(t, h[1], Durations{list: []time.Duration{60, 60, 80, 80}})
	}
}

func TestDurationsStats(t *testing.T) {
	return
	n := 1_000_000
	d := NewDurations(0)

	fill := Takes(func() {
		for i := 0; i <= n; i++ {
			d.Append(time.Millisecond * time.Duration(rand.Int31n(2000)))
		}
	})

	metrics := Takes(func() {
		fmt.Println("Metrics (only for queries without error):")
		fmt.Printf("  Sum:    %6s ms\n", d.Sum())
		fmt.Printf("  Min:    %6s ms\n", d.Min())
		fmt.Printf("  Max:    %6s ms\n", d.Max())
		fmt.Printf("  Median: %6s ms\n", d.Median())
		fmt.Printf("  Mean:   %6s ms\n", d.Mean())
	})

	fmt.Println()
	// 1 million is ~8M
	// 24bytes per entry? Hm
	fmt.Println("size: ", zruntime.SizeOf(d)/1024, "K")
	fmt.Println("fill: ", fill.Round(time.Millisecond))
	fmt.Println("print:", metrics.Round(time.Millisecond))
}

func BenchmarkTop(b *testing.B) {
	d := NewDurations(0)
	d.Grow(1_000_000)
	for i := 0; i < 1_000_000; i++ {
		d.Append(time.Duration(rand.Intn(200)))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = d.Top(20)
	}
}
