package ztime

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"zgo.at/zstd/ztest"
)

func TestNew(t *testing.T) {
	tz, err := time.LoadLocation("Asia/Makassar")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		in   string
		want time.Time
	}{
		{"2020-06-18 14:15:16.999999999", time.Date(2020, 6, 18, 14, 15, 16, 999999999, time.UTC)},
		{"2020-06-18 14:15:16.0", time.Date(2020, 6, 18, 14, 15, 16, 0, time.UTC)},
		{"2020-06-18 14:15:16", time.Date(2020, 6, 18, 14, 15, 16, 0, time.UTC)},
		{"2020-06-18 14:15", time.Date(2020, 6, 18, 14, 15, 0, 0, time.UTC)},
		{"2020-06-18 14", time.Date(2020, 6, 18, 14, 0, 0, 0, time.UTC)},
		{"2020-06-18", time.Date(2020, 6, 18, 0, 0, 0, 0, time.UTC)},
		{"2020-06", time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC)},
		{"2020", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},

		{"2020-06-18 14:15:16.999999999 WITA", time.Date(2020, 6, 18, 14, 15, 16, 999999999, tz)},
		{"2020-06-18 14:15:16.0 WITA", time.Date(2020, 6, 18, 14, 15, 16, 0, tz)},
		{"2020-06-18 14:15:16 WITA", time.Date(2020, 6, 18, 14, 15, 16, 0, tz)},
		{"2020-06-18 14:15 WITA", time.Date(2020, 6, 18, 14, 15, 0, 0, tz)},
		{"2020-06-18 14 WITA", time.Date(2020, 6, 18, 14, 0, 0, 0, tz)},
		{"2020-06-18 WITA", time.Date(2020, 6, 18, 0, 0, 0, 0, tz)},
		{"2020-06 WITA", time.Date(2020, 6, 1, 0, 0, 0, 0, tz)},
		{"2020 WITA", time.Date(2020, 1, 1, 0, 0, 0, 0, tz)},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			have := New(tt.in)
			if !have.Equal(tt.want) {
				t.Errorf("\nhave: %s\nwant: %s", have, tt.want)
			}
		})
	}
}

func TestStartOf(t *testing.T) {
	var (
		periods = []Period{Second, Minute, QuarterHour, HalfHour, Hour, Day, WeekMonday, WeekSunday, Month, Quarter, HalfYear, Year}
		f       = "2006-01-02 15:04:05.999999999"
		tt      = Time{time.Date(2020, 6, 18, 14, 49, 20, 666, time.UTC)}
		h       = new(strings.Builder)
	)
	h.WriteString("       StartOf: " + tt.Format(f) + "\n")
	for _, p := range periods {
		pad := strings.Repeat(" ", 14-len(p.String()))
		fmt.Fprintln(h, p.String(), pad,
			tt.StartOf(p).Format(f),
			"", tt.StartOf(p).AddTime(-1).StartOf(p).Format(f),
			"", tt.StartOf(p).AddTime(-1).StartOf(p).AddTime(-1).StartOf(p).Format(f),
			"", tt.StartOf(p).AddTime(-1).StartOf(p).AddTime(-1).StartOf(p).AddTime(-1).StartOf(p).Format(f))
	}

	have := h.String()
	want := `
		       StartOf: 2020-06-18 14:49:20.000000666
		second          2020-06-18 14:49:20  2020-06-18 14:49:19  2020-06-18 14:49:18  2020-06-18 14:49:17
		minute          2020-06-18 14:49:00  2020-06-18 14:48:00  2020-06-18 14:47:00  2020-06-18 14:46:00
		quarter hour    2020-06-18 14:45:00  2020-06-18 14:30:00  2020-06-18 14:15:00  2020-06-18 14:00:00
		half hour       2020-06-18 14:30:00  2020-06-18 14:00:00  2020-06-18 13:30:00  2020-06-18 13:00:00
		hour            2020-06-18 14:00:00  2020-06-18 13:00:00  2020-06-18 12:00:00  2020-06-18 11:00:00
		day             2020-06-18 00:00:00  2020-06-17 00:00:00  2020-06-16 00:00:00  2020-06-15 00:00:00
		week (Monday)   2020-06-15 00:00:00  2020-06-08 00:00:00  2020-06-01 00:00:00  2020-05-25 00:00:00
		week (Sunday)   2020-06-14 00:00:00  2020-06-07 00:00:00  2020-05-31 00:00:00  2020-05-24 00:00:00
		month           2020-06-01 00:00:00  2020-05-01 00:00:00  2020-04-01 00:00:00  2020-03-01 00:00:00
		quarter         2020-04-01 00:00:00  2020-01-01 00:00:00  2019-10-01 00:00:00  2019-07-01 00:00:00
		half year       2020-01-01 00:00:00  2019-07-01 00:00:00  2019-01-01 00:00:00  2018-07-01 00:00:00
		year            2020-01-01 00:00:00  2019-01-01 00:00:00  2018-01-01 00:00:00  2017-01-01 00:00:00`
	if d := ztest.Diff(have, want, ztest.DiffNormalizeWhitespace); d != "" {
		t.Error(d)
	}
}

func TestEndOf(t *testing.T) {
	var (
		periods = []Period{Second, Minute, QuarterHour, HalfHour, Hour, Day, WeekMonday, WeekSunday, Month, Quarter, HalfYear, Year}
		f       = "2006-01-02 15:04:05.999999999"
		tt      = Time{time.Date(2020, 6, 18, 14, 49, 20, 666, time.UTC)}
		h       = new(strings.Builder)
	)
	h.WriteString("\n         EndOf: " + tt.Format(f) + "\n")
	for _, p := range periods {
		pad := strings.Repeat(" ", 14-len(p.String()))
		fmt.Fprintln(h, p.String(), pad,
			tt.EndOf(p).Format(f),
			"", tt.EndOf(p).AddTime(1).EndOf(p).Format(f),
			"", tt.EndOf(p).AddTime(1).EndOf(p).AddTime(1).EndOf(p).Format(f),
			"", tt.EndOf(p).AddTime(1).EndOf(p).AddTime(1).EndOf(p).AddTime(1).EndOf(p).Format(f))
	}

	have := h.String()
	want := `
		         EndOf: 2020-06-18 14:49:20.000000666
		second          2020-06-18 14:49:20.999999999  2020-06-18 14:49:21.999999999  2020-06-18 14:49:22.999999999  2020-06-18 14:49:23.999999999
		minute          2020-06-18 14:49:59.999999999  2020-06-18 14:50:59.999999999  2020-06-18 14:51:59.999999999  2020-06-18 14:52:59.999999999
		quarter hour    2020-06-18 14:59:59.999999999  2020-06-18 15:14:59.999999999  2020-06-18 15:29:59.999999999  2020-06-18 15:44:59.999999999
		half hour       2020-06-18 14:59:59.999999999  2020-06-18 15:29:59.999999999  2020-06-18 15:59:59.999999999  2020-06-18 16:29:59.999999999
		hour            2020-06-18 14:59:59.999999999  2020-06-18 15:59:59.999999999  2020-06-18 16:59:59.999999999  2020-06-18 17:59:59.999999999
		day             2020-06-18 23:59:59.999999999  2020-06-19 23:59:59.999999999  2020-06-20 23:59:59.999999999  2020-06-21 23:59:59.999999999
		week (Monday)   2020-06-21 23:59:59.999999999  2020-06-28 23:59:59.999999999  2020-07-05 23:59:59.999999999  2020-07-12 23:59:59.999999999
		week (Sunday)   2020-06-20 23:59:59.999999999  2020-06-27 23:59:59.999999999  2020-07-04 23:59:59.999999999  2020-07-11 23:59:59.999999999
		month           2020-06-30 23:59:59.999999999  2020-07-31 23:59:59.999999999  2020-08-31 23:59:59.999999999  2020-09-30 23:59:59.999999999
		quarter         2020-06-30 23:59:59.999999999  2020-09-30 23:59:59.999999999  2020-12-31 23:59:59.999999999  2021-03-31 23:59:59.999999999
		half year       2020-06-30 23:59:59.999999999  2020-12-31 23:59:59.999999999  2021-06-30 23:59:59.999999999  2021-12-31 23:59:59.999999999
		year            2020-12-31 23:59:59.999999999  2021-12-31 23:59:59.999999999  2022-12-31 23:59:59.999999999  2023-12-31 23:59:59.999999999`
	if d := ztest.Diff(have, want, ztest.DiffNormalizeWhitespace); d != "" {
		t.Error(d)
	}
}

func TestWeek(t *testing.T) {
	var (
		mon = Time{time.Date(2021, 4, 5, 14, 49, 20, 666, time.UTC)}
		sun = Time{time.Date(2021, 4, 4, 14, 49, 20, 666, time.UTC)}
		f   = "Mon Jan _2"
		h   = new(strings.Builder)
	)

	h.WriteString("Monday:\n")
	for i := 0; i < 7; i++ {
		fmt.Fprintf(h, "%d  %s → %s %s\n", i, mon.Add(i, Day).Format(f),
			mon.Add(i, Day).StartOf(Week(false)).Format(f),
			mon.Add(i, Day).EndOf(Week(false)).Format(f))
	}
	h.WriteString("\nSunday:\n")
	for i := 0; i < 7; i++ {
		fmt.Fprintf(h, "%d  %s → %s %s\n", i, sun.Add(i, Day).Format(f),
			sun.Add(i, Day).StartOf(Week(true)).Format(f),
			sun.Add(i, Day).EndOf(Week(true)).Format(f))
	}

	have := h.String()
	want := `
		Monday:
		0  Mon Apr  5 → Mon Apr  5 Sun Apr 11
		1  Tue Apr  6 → Mon Apr  5 Sun Apr 11
		2  Wed Apr  7 → Mon Apr  5 Sun Apr 11
		3  Thu Apr  8 → Mon Apr  5 Sun Apr 11
		4  Fri Apr  9 → Mon Apr  5 Sun Apr 11
		5  Sat Apr 10 → Mon Apr  5 Sun Apr 11
		6  Sun Apr 11 → Mon Apr  5 Sun Apr 11

		Sunday:
		0  Sun Apr  4 → Sun Apr  4 Sat Apr 10
		1  Mon Apr  5 → Sun Apr  4 Sat Apr 10
		2  Tue Apr  6 → Sun Apr  4 Sat Apr 10
		3  Wed Apr  7 → Sun Apr  4 Sat Apr 10
		4  Thu Apr  8 → Sun Apr  4 Sat Apr 10
		5  Fri Apr  9 → Sun Apr  4 Sat Apr 10
		6  Sat Apr 10 → Sun Apr  4 Sat Apr 10`
	if d := ztest.Diff(have, want, ztest.DiffNormalizeWhitespace); d != "" {
		t.Error(d)
	}
}

func TestAdd(t *testing.T) {
	var (
		f = "2006-01-02 15:04:05"
		h = new(strings.Builder)
	)
	test := func(p Period, tt time.Time) {
		pad := strings.Repeat(" ", 10-len(p.String()))
		fmt.Fprintln(h, p.String(), tt.Format(f), pad,
			"", Add(tt, 1, p).Format(f),
			"", Add(tt, 2, p).Format(f),
			"", Add(tt, 3, p).Format(f),
			"\n"+strings.Repeat(" ", 32),
			Add(tt, -1, p).Format(f),
			"", Add(tt, -2, p).Format(f),
			"", Add(tt, -3, p).Format(f))
	}
	test(Month, time.Date(2020, 1, 31, 14, 49, 20, 666, time.UTC))
	test(Quarter, time.Date(2020, 1, 31, 14, 49, 20, 666, time.UTC))
	test(HalfYear, time.Date(2020, 3, 31, 14, 49, 20, 666, time.UTC))
	test(Year, time.Date(2020, 2, 29, 14, 49, 20, 666, time.UTC))

	have := h.String()
	want := `
		month 2020-01-31 14:49:20        2020-02-29 14:49:20  2020-03-31 14:49:20  2020-04-30 14:49:20
										 2019-12-31 14:49:20  2019-11-30 14:49:20  2019-10-31 14:49:20
		quarter 2020-01-31 14:49:20      2020-04-30 14:49:20  2020-07-31 14:49:20  2020-10-31 14:49:20
										 2019-10-31 14:49:20  2019-07-31 14:49:20  2019-04-30 14:49:20
		half year 2020-03-31 14:49:20    2020-09-30 14:49:20  2021-03-31 14:49:20  2021-09-30 14:49:20
										 2019-09-30 14:49:20  2019-03-31 14:49:20  2018-09-30 14:49:20
		year 2020-02-29 14:49:20         2021-02-28 14:49:20  2022-02-28 14:49:20  2023-02-28 14:49:20
										 2019-02-28 14:49:20  2018-02-28 14:49:20  2017-02-28 14:49:20`
	if d := ztest.Diff(have, want, ztest.DiffNormalizeWhitespace); d != "" {
		t.Error(d)
	}
}
