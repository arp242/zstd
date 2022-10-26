// Package ztime implements functions for date and time.
package ztime

import (
	"strings"
	"time"
)

// Time wraps time.Time to add {Start,End}Of, AddPeriod(), and Unpack(). It also
// wraps operations that return time.Time to return ztime.Time.
type Time struct{ time.Time }

// Proxy to time.Time, returning ztime.Time.

func (t Time) Add(d time.Duration) Time             { return Time{t.Time.Add(d)} }
func (t Time) AddDate(years, months, days int) Time { return Time{t.Time.AddDate(years, months, days)} }
func (t Time) In(loc *time.Location) Time           { return Time{t.Time.In(loc)} }
func (t Time) Local() Time                          { return Time{t.Time.Local()} }
func (t Time) Round(d time.Duration) Time           { return Time{t.Time.Round(d)} }
func (t Time) Truncate(d time.Duration) Time        { return Time{t.Time.Truncate(d)} }
func (t Time) UTC() Time                            { return Time{t.Time.UTC()} }

// Our new methods.

func (t Time) StartOf(p Period) Time          { return Time{StartOf(t.Time, p)} }
func (t Time) EndOf(p Period) Time            { return Time{EndOf(t.Time, p)} }
func (t Time) AddPeriod(n int, p Period) Time { return Time{AddPeriod(t.Time, n, p)} }
func (t Time) Unpack() (year int, month time.Month, day, hour, minute, second, nanosecond int, loc *time.Location) {
	return t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()
}

// MustParse is like time.Parse, but will panic on errors.
func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic("ztime.MustParse: " + err.Error())
	}
	return t
}

// FromString creates a new date from a string according to the layout:
//
//	2006-01-02 15:04:05 MST
//
// Any part on the right can be omitted; for example "2020-01-01" will create a
// new date without any time, or "2020-01-01 13" will create a date with the
// hour set.
//
// A timezone can always be added, for example "2020-01-01 13 CET".
//
// This will panic on errors. This is mostly useful in tests to quickly create a
// date without too much ceremony.
func FromString(s string) time.Time {
	tz := strings.LastIndexByte(s, ' ')
	if tz > -1 && strings.ContainsAny(s[tz:], "0123456789") {
		tz = -1
	}
	ss := s
	if tz > -1 {
		ss = s[:tz]
	}

	layout := "2006-01-02 15:04:05"
	if len(ss) < 19 {
		layout = layout[:len(ss)]
	}
	if tz > -1 {
		layout += " MST"
	}
	return MustParse(layout, s)
}

// Unpack a time to its individual components.
func Unpack(t time.Time) (year int, month time.Month, day, hour, minute, second, nanosecond int, loc *time.Location) {
	return t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()
}
