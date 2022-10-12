package ztime

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// Durations is a list of time.Durations.
//
// This is useful if you want to record a list of durations; for example to
// keep track of performance metrics.
//
// All operations are thread-safe.
type Durations struct {
	mu      *sync.Mutex
	list    []time.Duration
	maxSize int
	data    map[int]any
	off     int

	modified      bool
	sum, min, max time.Duration
}

// NewDurations creates a new Durations.
//
// The maximum size of the list is set to max items. After this, the oldest
// entries will be removed. Set to not have an upper limit.
func NewDurations(max int) Durations {
	return Durations{
		modified: true,
		mu:       new(sync.Mutex),
		data:     make(map[int]any),
		maxSize:  max,
	}
}

func (d Durations) String() string {
	return fmt.Sprintf("%d durations, from %s to %s", d.Len(), d.Min(), d.Max())
}

// List returns a copy of all durations.
func (d Durations) List() []time.Duration {
	d.mu.Lock()
	cpy := make([]time.Duration, len(d.list))
	copy(cpy, d.list)
	d.mu.Unlock()
	return cpy
}

// Grow the list of durations for another n durations.
//
// After Grow(n), at least n items can be appended without another allocation.
//
// Grow will panic if n is negative or if the list can't grow (i.e. larger than
// MaxInt items).
func (d *Durations) Grow(n int) {
	if n < 0 {
		panic("ztime.Durations.Grow: negative count")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// Enough space, don't need to do anything.
	if cap(d.list) > len(d.list)+n {
		return
	}

	cpy := make([]time.Duration, len(d.list), len(d.list)+n)
	copy(cpy, d.list)
	d.list = cpy
}

// Append durations to the list.
func (d *Durations) Append(add ...time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.append(add...)
}

// appendWithData appends an item to this list and sets meta-data for it.
//
// TODO: unexported for now; the idea here is that you can add some meta-data to
// specific duration entries (i.e. which parameters, query, whatnot were used
// for this specific duration), but I can't really think of a convenient API to
// access this right now.
//
// Grouping by data is probably a common use case, so storing it a bit different
// (map[string][]int, as data → indexes) might be a good idea(?)
func (d *Durations) appendWithData(add time.Duration, data any) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.append(add)
	d.data[d.off+len(d.list)-1] = data
}

// Len returns the number of durations in this list.
func (d Durations) Len() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.list)
}

// Sum returns the sum of all durations in this list.
func (d Durations) Sum() time.Duration {
	if d.modified {
		d.mu.Lock()
		defer d.mu.Unlock()

		var sum time.Duration
		for _, l := range d.list {
			sum += l
		}
		d.sum = sum
		d.modified = false
	}
	return d.sum
}

// Min returns the minimum value in this list.
func (d Durations) Min() time.Duration {
	if d.Len() == 0 {
		return 0
	}
	if d.modified {
		d.mu.Lock()
		defer d.mu.Unlock()

		m := time.Duration(math.MaxInt64)
		for _, l := range d.list {
			if l < m {
				m = l
			}
		}
		d.min = m
		d.modified = false
	}
	return d.min
}

// Max returns the maximum value in this list.
func (d Durations) Max() time.Duration {
	if d.modified {
		d.mu.Lock()
		defer d.mu.Unlock()

		var m time.Duration
		for _, l := range d.list {
			if l > m {
				m = l
			}
		}

		d.max = m
		d.modified = false
	}
	return d.max
}

// Mean returns the mean average of the durations in this list.
func (d Durations) Mean() time.Duration {
	if d.Len() == 0 {
		return 0
	}
	return d.Sum() / time.Duration(d.Len())
}

// Median returns the median average of the durations in this list.
func (d Durations) Median() time.Duration {
	if d.Len() == 0 {
		return 0
	}

	cpy := d.copyAndSort()
	return cpy[len(cpy)/2]
}

// Top gets the top percent items in this list.
//
// e.g. with a list of 200 durations, Top(10) returns the top 10% longest
// durations (20 of them). Use negative numbers to get the bottom percentile
// (the fastest 10%).
//
// The return value is a copy.
//
// Numbers higher than 100 or lower than -100 will panic.
func (d Durations) Top(percent int) Durations {
	if percent > 100 || percent < -100 {
		panic("ztime.Durations.Top: percent out of bounds")
	}

	// This is a bit funky, but we want to:
	//
	// 1. Keep this Durations unmodified.
	// 2. Retain the original order (rather than a sorted order) in the return
	//    value.
	// 3. Copy the data as well.
	//
	// So get the "cutoff" point and then create a new list in a loop.
	//
	// None of this wins any prizes for performance; it takes ~75ms to get the
	// top 20% of list with 1 million items on my laptop, so it's not *too* bad.
	// Almost all of that is actually in the sorting; we can improve on that by
	// using a tree rather than a list, but that would make Append() slower, and
	// Top() is a rare enough operation that I'd rather have a faster Append()
	// and slower Top(), and in reality most lists will be far fewer than 1M
	// items.

	var (
		cpy    = d.copyAndSort()
		num    = int(float64(len(cpy)) * (float64(percent) / 100))
		newd   = NewDurations(d.maxSize)
		cutoff time.Duration
	)
	if percent > 0 {
		cutoff = cpy[len(cpy)-num-1]
	} else {
		num = -num
		cutoff = cpy[num]
	}
	newd.Grow(num)

	d.mu.Lock()
	defer d.mu.Unlock()
	d.modified = true
	for i, l := range d.list {
		if len(newd.list) >= num {
			break
		}
		if (percent < 0 && l > cutoff) || (percent > 0 && l <= cutoff) {
			continue
		}

		newd.list = append(newd.list, l)
		if data := d.data[d.off+i]; data != nil {
			newd.data[len(newd.list)-1] = data
		}
	}
	return newd
}

// Distrubute the list of Durations in n blocks.
//
// For example with Distribute(5) it returns 5 set of durations, from the fastest
// 20% to the slowest 20%.
func (d Durations) Distrubute(n int) []Durations {
	var (
		cpy   = d.copyAndSort()
		bins  = make([]Durations, n)
		num   = int(float64(len(cpy)) * (1.0 / float64(n)))
		min   = d.Min()
		bSize = (d.Max() - min) / time.Duration(n)
	)
	for i := range bins {
		bins[i] = NewDurations(d.maxSize)
		bins[i].Grow(num)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	b := 0
	for i, l := range cpy {
		if l-min > bSize*(time.Duration(b+1)) && b < n-1 {
			b++
		}

		bins[b].list = append(bins[b].list, l)
		if data := d.data[d.off+i]; data != nil {
			bins[b].data[len(bins[b].list)-1] = data
		}
	}

	return bins
}

// Not locked!
func (d *Durations) append(add ...time.Duration) {
	d.list = append(d.list, add...)

	if d.maxSize > 0 && len(d.list) > d.maxSize {
		l := len(d.list) - d.maxSize
		for rm := range d.list[:l] {
			delete(d.data, d.off+l+rm)
		}
		d.off += l

		d.list = d.list[l:]
	}
}

func (d Durations) copyAndSort() []time.Duration {
	d.mu.Lock()
	cpy := make([]time.Duration, len(d.list))
	copy(cpy, d.list)
	d.mu.Unlock()

	sort.Slice(cpy, func(i, j int) bool { return cpy[i] < cpy[j] })
	return cpy
}
