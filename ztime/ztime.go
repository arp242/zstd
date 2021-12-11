// Package ztime implements functions for date and time.
package ztime

import (
	"strings"
	"time"
)

// Time wraps time.Time to add {Start,End}Of and Add(). To make it a bit more
// convenient to run several operations.
//
// time.Time's Add() is renamed to AddTime() here.
type Time struct{ time.Time }

// Proxy to time.Time
//
// TODO: maybe make this a drop-in replacement? That would mean copying all top-level
// time.* functions/constants, as well as proxy all time.Time methods (including
// not renaming AddTime(). Not sure if it's worth it?

func (t Time) AddTime(d time.Duration) Time         { return Time{t.Time.Add(d)} }
func (t Time) AddDate(years, months, days int) Time { return Time{t.Time.AddDate(years, months, days)} }
func (t Time) In(loc *time.Location) Time           { return Time{t.Time.In(loc)} }
func (t Time) Local() Time                          { return Time{t.Time.Local()} }
func (t Time) Round(d time.Duration) Time           { return Time{t.Time.Round(d)} }
func (t Time) Truncate(d time.Duration) Time        { return Time{t.Time.Truncate(d)} }
func (t Time) UTC() Time                            { return Time{t.Time.UTC()} }

func (t Time) StartOf(p Period) Time    { return Time{StartOf(t.Time, p)} }
func (t Time) EndOf(p Period) Time      { return Time{EndOf(t.Time, p)} }
func (t Time) Add(n int, p Period) Time { return Time{Add(t.Time, n, p)} }

// MustParse is like time.Parse, but will panic on errors.
func MustParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic("ztime.MustParse: " + err.Error())
	}
	return t
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
	// TODO: rename to something else; maybe FromString()? Quick()?
	// And add a New(time.Time) ztime.Time
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
