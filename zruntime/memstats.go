package zruntime

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	_ = 1 << iota

	Descriptions // Descriptions adds a brief description of every field.

	General // General lists general memory stats.
	Heap    // Heap lists heap memory stats.
	Stack   // Stack lists stack memory stats.
	GC      // GC lists garbage collector stats.
	//Internal
)

// MemStats is a wrapper around runtime.MemStats with some better printing of
// values.
//
// After creation you can use String() to get the curreny memory stats, or
// Print() to print them to stdout.
//
// A second call to Print() will add a comparison of the stats of the first
// call. You can use Read() to just read the current memory stats without doing
// anything.
type MemStats struct {
	mu sync.Mutex

	runtime.MemStats
	values, prev reflect.Value
	prevRead     time.Time
	opts         int
	out          *strings.Builder
}

// Reset any previously stored memory stats.
func (m *MemStats) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.MemStats = runtime.MemStats{}
	m.values = reflect.Value{}
	m.prev = reflect.Value{}
	m.prevRead = time.Time{}
}

// Read the current memory stats, moving any current value to the previous one.
func (m *MemStats) Read() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.prev.IsValid() {
		m.prev = reflect.ValueOf(struct{}{})
	}
	if m.TotalAlloc > 0 && m.prevRead.IsZero() {
		m.prev = reflect.ValueOf(m.MemStats)
		m.prevRead = time.Now()
	}
	runtime.ReadMemStats(&m.MemStats)
	m.values = reflect.ValueOf(m.MemStats)
}

// ResetAndRead calls Reset() and Read().
func (m *MemStats) ResetAndRead() {
	m.Reset()
	m.Read()
}

// Print the current memory stats.
func (m *MemStats) Print(opts int) {
	m.opts = opts
	fmt.Println(m.String())
}

// Get the memory stats; this is the same as String() except that you can set
// options.
func (m *MemStats) Get(opts int) string {
	m.opts = opts
	return m.String()
}

// String gets the memory stats.
func (m *MemStats) String() string {
	m.Read()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.opts == 0 {
		m.opts = General | Heap | Stack | GC
	}

	m.out = new(strings.Builder)
	m.out.WriteString("\t                ")
	m.out.WriteString(time.Now().Format("15:04:05.0000"))
	if !m.prevRead.IsZero() {
		m.out.WriteString("               ")
		m.out.WriteString(m.prevRead.Format("15:04:05.0000"))
	}
	m.out.WriteRune('\n')

	if m.opts&General != 0 {
		m.value("TotalAlloc")
		m.value("Sys")
		m.value("Mallocs")
		m.value("Frees")
		m.value("Live")
	}
	if m.opts&Heap != 0 {
		m.value("HeapAlloc")
		m.value("HeapSys")
		m.value("HeapIdle")
		m.value("HeapInuse")
		m.value("HeapReleased")
		m.value("HeapObjects")
	}
	if m.opts&Stack != 0 {
		m.value("StackInuse")
		m.value("StackSys")
	}

	// if m.opts&Internal != 0 {
	// 	m.value("Lookups")
	// 	// TODO
	// }

	if m.opts&GC != 0 {
		m.value("NextGC")
		m.value("LastGC")
		m.value("PauseTotalNs")
		m.value("NumGC")
		m.value("NumForcedGC")
		m.value("GCCPUFraction")
	}

	return m.out.String()
}

func (m *MemStats) value(k string) {
	fmt.Fprintf(m.out, "\t%-14s", k)
	v := m.values.FieldByName(k)
	v2 := m.prev.FieldByName(k)

	switch k {
	case "TotalAlloc", "Sys", "HeapAlloc", "HeapSys", "HeapIdle", "HeapInuse",
		"HeapReleased", "StackInuse", "StackSys", "NextGC":
		m.bytes(v, v2)
	case "Lookups", "Mallocs", "Frees", "HeapObjects", "NumGC", "NumForcedGC":
		m.number(v, v2)
	case "PauseTotalNs":
		m.duration(v, v2)
	case "LastGC":
		m.time(v, v2)
	case "GCCPUFraction":
		m.float(v, v2)
	case "Live":
		a := m.values.FieldByName("Mallocs")
		b := m.values.FieldByName("Frees")
		live := reflect.ValueOf(a.Uint() - b.Uint())

		var live2 reflect.Value
		a = m.prev.FieldByName("Mallocs")
		if a.IsValid() {
			b = m.prev.FieldByName("Frees")
			live2 = reflect.ValueOf(a.Uint() - b.Uint())
		}

		m.number(live, live2)
	}

	if m.opts&Descriptions != 0 {
		m.out.WriteString("     ")
		m.out.WriteString(descriptions[k])
	}

	m.out.WriteRune('\n')
}

func (m *MemStats) duration(v, p reflect.Value) {
	d := time.Duration(v.Uint())
	fmt.Fprintf(m.out, "%15s", d)

	if p.IsValid() {
		diff := d - time.Duration(p.Uint())
		if diff == 0 {
			fmt.Fprintf(m.out, " %15s %s", "-", m.perc(0, 0))
			return
		}

		// TODO: %
		fmt.Fprintf(m.out, " %+15s %11s", diff, "-")
	}
}

func (m *MemStats) time(v, p reflect.Value) {
	t := time.Unix(0, int64(v.Uint())).In(time.Local)
	fmt.Fprintf(m.out, "%15s", t.Format("15:04:05.0000"))

	if p.IsValid() {
		ns := int64(p.Uint())
		if ns == 0 {
			fmt.Fprintf(m.out, " %15s %11s", "-", "-")
			return
		}

		t2 := time.Unix(0, ns).In(time.Local)
		diff := t.Sub(t2)
		fmt.Fprintf(m.out, " %15s %11s", t2.Format("15:04:05.0000"), diff)
	}
}

func (m *MemStats) bytes(v, p reflect.Value) {
	fmt.Fprintf(m.out, "%14.2fk", float64(v.Uint())/1024)

	if p.IsValid() {
		// TODO: overflow?
		diff := int64(v.Uint()) - int64(p.Uint())
		if diff == 0 {
			fmt.Fprintf(m.out, " %15s %s", "-", m.perc(0, 0))
			return
		}
		fmt.Fprintf(m.out, " %+14.2fk %s", float64(diff)/1024, m.perc(diff, v.Uint()))
	}
}

func (m *MemStats) float(v, p reflect.Value) {
	fmt.Fprintf(m.out, "%15.10f", v.Float())

	if p.IsValid() {
		fmt.Fprintf(m.out, " %+15.10f %11s", v.Float()-p.Float(), "-")
	}
}

func (m *MemStats) number(v, p reflect.Value) {
	fmt.Fprintf(m.out, "%15d", v.Uint())

	if p.IsValid() {
		diff := int64(v.Uint()) - int64(p.Uint())
		if diff == 0 {
			fmt.Fprintf(m.out, " %15s %s", "-", m.perc(0, 0))
			return
		}
		fmt.Fprintf(m.out, " %+15d %s", diff, m.perc(diff, v.Uint()))
	}
}

func (m *MemStats) perc(diff int64, orig uint64) string {
	if diff == 0 {
		return fmt.Sprintf("%+11s", "-")
	}
	return fmt.Sprintf("%+10.1f%%", float64(diff)/float64(orig)*100)
}

var descriptions = map[string]string{
	// General
	"TotalAlloc": "Total amount of memory that has been allocated",
	"Sys":        "OS memory in current use; sum of *Sys",
	"Mallocs":    "Total amount of heap objects allocated",
	"Frees":      "Total amount of heap objects freed",
	"Live":       "Live objects: Allocs - Frees",

	// Heap
	"HeapAlloc":    "All reachable objects and unreachable not yet freed by GC",
	"HeapSys":      "Heap memory obtained from the OS",
	"HeapIdle":     "Cached memory; may be returned to OS",
	"HeapInuse":    "Heap memory currently in use",
	"HeapReleased": "Memory returned to the OS",
	"HeapObjects":  "Number of objects on the heap",

	// Stack
	"StackInuse": "Stack memory currently in use",
	"StackSys":   "Stack memory obtained from OS",

	// Internal
	"Lookups":     "Number of pointer lookups",
	"MSpanInuse":  "",
	"MSpanSys":    "",
	"MCacheInuse": "",
	"MCacheSys":   "",
	"BuckHashSys": "",
	"GCSys":       "",
	"OtherSys":    "",

	// Garbage collector
	"NextGC":        "Target heap size of next GC cycle",
	"LastGC":        "Last garbage collection finished",
	"PauseTotalNs":  "Total amount of time the world was stopped",
	"NumGC":         "Number of completed GC cycles",
	"NumForcedGC":   "Total number of forced GC cycled (runtime.GC())",
	"GCCPUFraction": "Fraction of CPU times used by GC (between 0 and 1)",

	// PauseNs is a circular buffer of recent GC stop-the-world
	// pause times in nanoseconds.
	//
	// The most recent pause is at PauseNs[(NumGC+255)%256]. In
	// general, PauseNs[N%256] records the time paused in the most
	// recent N%256th GC cycle. There may be multiple pauses per
	// GC cycle; this is the sum of all pauses during a cycle.
	//PauseNs [256]uint64

	// PauseEnd is a circular buffer of recent GC pause end times,
	// as nanoseconds since 1970 (the UNIX epoch).
	//
	// This buffer is filled the same way as PauseNs. There may be
	// multiple pauses per GC cycle; this records the end of the
	// last pause in a cycle.
	//PauseEnd [256]uint64
}
