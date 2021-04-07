// Package ztime implements functions for date and time.
package ztime

import (
	"strings"
	"testing"
	"time"
)

// Time wraps time.Time to add {Start,End}Of and Add(). To make it a bit more
// convenient to run several operations.
//
// time.Time's Add() is renamed to AddTime() here.
type Time struct{ time.Time }

func (t Time) StartOf(p Period) Time    { return Time{StartOf(t.Time, p)} }
func (t Time) EndOf(p Period) Time      { return Time{EndOf(t.Time, p)} }
func (t Time) Add(n int, p Period) Time { return Time{Add(t.Time, n, p)} }

func (t Time) AddTime(d time.Duration) Time         { return Time{t.Time.Add(d)} }
func (t Time) AddDate(years, months, days int) Time { return Time{t.Time.AddDate(years, months, days)} }
func (t Time) In(loc *time.Location) Time           { return Time{t.Time.In(loc)} }
func (t Time) Local() Time                          { return Time{t.Time.Local()} }
func (t Time) Round(d time.Duration) Time           { return Time{t.Time.Round(d)} }
func (t Time) Truncate(d time.Duration) Time        { return Time{t.Time.Truncate(d)} }
func (t Time) UTC() Time                            { return Time{t.Time.UTC()} }

// Period to adjust or align a time by.
type Period uint8

// Periods to adjust or align a time by.
const (
	_ Period = iota + 1
	Second
	Minute
	QuarterHour
	HalfHour
	Hour
	Day
	WeekMonday
	WeekSunday
	Month
	Quarter
	HalfYear
	Year
)

func (p Period) String() string {
	switch p {
	case Second:
		return "second"
	case Minute:
		return "minute"
	case QuarterHour:
		return "quarter hour"
	case HalfHour:
		return "half hour"
	case Hour:
		return "hour"
	case Day:
		return "day"
	case WeekMonday:
		return "week (Monday)"
	case WeekSunday:
		return "week (Sunday)"
	case Month:
		return "month"
	case Quarter:
		return "quarter"
	case HalfYear:
		return "half year"
	case Year:
		return "year"
	default:
		return ""
	}
}

// Week returns WeekSunday or WeekMonday.
func Week(sundayStartsWeek bool) Period {
	if sundayStartsWeek {
		return WeekSunday
	}
	return WeekMonday
}

// Now returns the current time as UTC.
//
// This can be swapped out in tests with SetNow()
var Now = func() time.Time { return time.Now().UTC() }

// SetNow sets Now() and restores it when the test finishes.
//
// The date is parsed with New().
func SetNow(t *testing.T, s string) {
	t.Helper()

	d := New(s)
	Now = func() time.Time { return d }
	t.Cleanup(func() {
		Now = func() time.Time { return time.Now().UTC() }
	})
}

// New creates a new date from a string according to the layout:
//
//  2006-01-02 15:04:05 MST
//
// Any part on the right can be omitted; for example New("2020-01-01") will
// create a new date without any time, or New("2020-01-01 13") will create a
// date with the hour set.
//
// A timezone can always be added, for example New("2020-01-01 13 CET").
//
// This will panic on errors. This is mostly useful in tests to quickly create a
// date without too much ceremony.
func New(s string) time.Time {
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

// MustParse is like time.Parse, but will panic on errors.
func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic("ztime.MustParse: " + err.Error())
	}
	return t
}

// LeapYear reports if this year is a leap year according to the Gregorian
// calendar.
func LeapYear(t time.Time) bool {
	y := t.Year()
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}

// DaysInMonth returns the number of days for this month.
func DaysInMonth(t time.Time) int {
	switch t.Month() {
	default:
		return 30
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 2:
		if LeapYear(t) {
			return 29
		}
		return 28
	}
}

// LastInMonth reports if the current day is the last day in this month.
func LastInMonth(t time.Time) bool {
	return t.Day() == DaysInMonth(t)
}

// StartOf adjusts the time to the start of the given period.
//
// For example StartOf(t, QuarterHour) with "15:19" will adjust the time to
// "15:15".
func StartOf(t time.Time, p Period) time.Time {
	y, m, d, h, min, s, _, l := t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()
	ns := 0
	switch p {
	case Second:
		// ns is already set.
	case Minute:
		s = 0
	case QuarterHour:
		min, s = min-min%15, 0
	case HalfHour:
		min, s = min-min%30, 0
	case Hour:
		min, s = 0, 0
	case Day:
		h, min, s = 0, 0, 0
	case WeekMonday:
		wd := int(t.Weekday()) - 1
		if wd == -1 {
			wd = 6
		}
		d, h, min, s = d-wd, 0, 0, 0
	case WeekSunday:
		d, h, min, s = d-int(t.Weekday()), 0, 0, 0
	case Month:
		d, h, min, s = 1, 0, 0, 0
	case Quarter:
		m, d, h, min, s = m-(m-1)%3, 1, 0, 0, 0
	case HalfYear:
		m, d, h, min, s = m-(m-1)%6, 1, 0, 0, 0
	case Year:
		m, d, h, min, s = 1, 1, 0, 0, 0
	}
	return time.Date(y, time.Month(m), d, h, min, s, ns, l)
}

// EndOf adjusts the time to the end of the given period.
//
// For example EndOf(t, QuarterHour) with "15:19" will adjust the time to
// "15:30".
func EndOf(t time.Time, p Period) time.Time {
	y, m, d, h, min, s, _, l := t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()
	ns := 999999999
	switch p {
	case Second:
		// ns is already set.
	case Minute:
		s = 59
	case QuarterHour:
		min, s = StartOf(t, p).Add(14*time.Minute).Minute(), 59
	case HalfHour:
		min, s = StartOf(t, p).Add(29*time.Minute).Minute(), 59
	case Hour:
		min, s = 59, 59
	case Day:
		h, min, s = 23, 59, 59
	case WeekMonday, WeekSunday:
		t = StartOf(t, p).AddDate(0, 0, 6)
		m, d, h, min, s = t.Month(), t.Day(), 23, 59, 59
	case Month:
		d, h, min, s = DaysInMonth(t), 23, 59, 59
	case Quarter:
		t = StartOf(t, p).AddDate(0, 3, 0).Add(-1)
		m, d, h, min, s = t.Month(), DaysInMonth(t), 23, 59, 59
	case HalfYear:
		t = StartOf(t, p).AddDate(0, 6, 0).Add(-1)
		m, d, h, min, s = t.Month(), DaysInMonth(t), 23, 59, 59
	case Year:
		m, d, h, min, s = 12, 31, 23, 59, 59
	}
	return time.Date(y, m, d, h, min, s, ns, l)
}

// Add a time period.
//
// This matches common understanding of what things like "next month" mean;
// adding or subtracting months will always end up in the expected month,
// regardless of the number of days in either month.
//
// For example:
//
//   Jan 31 + 1 month  = Feb 28 (or Feb 29, if it's a leap year)
//   Dec 31 - 3 months = Sep 30
//
// This is done for Month, Quarter, and HalfYear.
//
// There is one special case for Year: if the date is Feb 29th, adding or
// subtracting a year will land you on Feb 28th.
//
// Since leap seconds are irregular and unpredictable they are not handled. The
// entire concept is silly and most programs should just pretend they don't
// exist.
func Add(t time.Time, n int, p Period) time.Time {
	y, m, d, h, min, s, ns, l := t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()
	switch p {
	case Second:
		s += n
	case Minute:
		min += n
	case QuarterHour:
		min += 15 * n
	case HalfHour:
		min += 30 * n
	case Hour:
		h += n
	case Day:
		h += 24 * n
	case WeekMonday, WeekSunday:
		d += 7 * n
	case Month:
		m += n
		if x := DaysInMonth(time.Date(y, time.Month(m), 1, h, min, s, ns, l)); d > x {
			d = x
		}
	case Quarter:
		return Add(t, n*3, Month)
	case HalfYear:
		return Add(t, n*6, Month)
	case Year:
		y += n
		if n := DaysInMonth(time.Date(y, time.Month(m), 1, h, min, s, ns, l)); DaysInMonth(t) > n {
			d = n // Deal with leap years
		}
	}
	return time.Date(y, time.Month(m), d, h, min, s, ns, l)
}
