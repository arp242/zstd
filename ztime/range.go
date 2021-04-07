package ztime

import (
	"strconv"
	"time"
)

// A Range represents a time range from Start to End.
//
// The timezone is always taken from the start time. The end's timezone will be
// adjusted if it differs.
type Range struct {
	// TODO: these should not be exported; they're only exported for now to make
	// migrating to Range a bit easier.
	Start, End time.Time
}

// NewRange creates a new range with the start date set.
func NewRange(start time.Time) Range {
	return Range{Start: start}
}

// From returns a copy with the start time set.
//
// This will apply the timezone from the passed start time to the end time.
func (r Range) From(start time.Time) Range {
	r.Start = start
	r.End = r.End.In(start.Location())
	return r
}

// To returns a copy with the end time set.
//
// This will apply the timezone from the start time to the passed end time.
func (r Range) To(end time.Time) Range {
	r.End = end.In(r.Start.Location())
	return r
}

// Period returns a copy withh the end time set to n Period from the start time.
//
// This uses ztime.Add() and its "common sense" understanding of months.
func (r Range) Period(n int, p Period) Range {
	r.End = Add(r.Start, n, p)
	return r
}

// In returns a copy with the timezone set to UTC.
func (r Range) UTC() Range {
	r.Start = r.Start.In(time.UTC)
	r.End = r.End.In(time.UTC)
	return r
}

// In returns a copy with the timezone set to loc.
func (r Range) In(loc *time.Location) Range {
	r.Start = r.Start.In(loc)
	r.End = r.End.In(loc)
	return r
}

// Current returns a copy with the start and end times set the current Period p.
//
// This uses the value of the start time. Any value in the end time is ignored.
//
// For example with NewRange("2020-06-18 14:00:00"):
//
//   Current(Month)       2020-06-01 00:00:00       → 2020-06-30 23:59:59
//   Current(WeekMonday)  2020-06-15 00:00:00 (Mon) → 2020-06-21 23:59:59 (Sun)
func (r Range) Current(p Period) Range {
	return Range{
		Start: StartOf(r.Start, p),
		End:   EndOf(r.Start, p),
	}
}

// Last returns a copy with the start and end times set to the last Period p.
//
// This uses the value of the start time. Any value in the end time is ignored.
//
// For example with NewRange("2020-06-18 14:00:00") (Thursday):
//
//   Last(Month)       2020-05-18 00:00:00       → 2020-06-18 23:59:59
//   Last(WeekMonday)  2020-06-11 00:00:00 (Wed) → 2020-06-18 23:59:59 (Thu)
func (r Range) Last(p Period) Range {
	// TODO: are we sure this is what we want? Wouldn't make e.g. Thursday to
	// Thursday make more sense?
	//
	//   Last(Month)       2020-05-18 00:00:00       → 2020-06-18 23:59:59
	//   Last(WeekMonday)  2020-06-12 00:00:00 (Thu) → 2020-06-18 23:59:59 (Thu)
	//
	// Need to see how it's used in GC.
	pp := map[Period]Period{
		Second:     0,
		Minute:     Second,
		Hour:       Minute,
		Day:        Hour,
		WeekMonday: Day,
		WeekSunday: Day,
		Month:      Day,
		Quarter:    Day,
		HalfYear:   Day,
		Year:       Day,
	}[p]

	return Range{
		Start: StartOf(Add(r.Start, -1, p), pp),
		End:   EndOf(r.Start, pp),
	}
}

// String shows the range from start to end as a human-readable representation;
// e.g. "current week", "last week", "previous month", etc.
//
// It falls back to "Mon Jan 2–Mon Jan 2" if there's no clear way to describe
// it.
func (r Range) String() string {
	today := StartOf(Now().In(r.Start.Location()), Day)
	r.Start, r.End = StartOf(r.Start, Day), StartOf(r.End, Day)

	d := r.Diff(Day, Month)
	n := strconv.Itoa
	addYear := func(t time.Time, s string) string {
		if t.Year() != today.Year() {
			return s + " 2006"
		}
		return s
	}

	// Selected one full month, display as month name.
	if d.Months == 0 && r.Start.Day() == 1 && LastInMonth(r.End) {
		return r.Start.Format(addYear(r.Start, "January"))
	}

	// From start of a month to end of another month.
	if d.Months > 1 && r.Start.Day() == 1 && LastInMonth(r.End) {
		return r.Start.Format(addYear(r.Start, "January")) + "–" + r.End.Format(addYear(r.End, "January"))
	}

	if d.Months == 0 && d.Days == 0 && StartOf(r.End.AddDate(0, 0, 1), Day).Equal(today) {
		return "Yesterday"
	}

	if r.End.Equal(today) {
		if d.Months == 0 {
			if d.Days == 0 {
				return "Today"
			}
			if d.Days == 1 {
				return "Yesterday–Today"
			}
		}

		if r.Start.Day() == r.End.Day() {
			if d.Months == 1 {
				return n(d.Months) + " month ago–Today"
			}
			return n(d.Months) + " months ago–Today"
		}
		if d.Days%7 == 0 {
			w := n(d.Days / 7)
			if w == "1" {
				return w + " week ago–Today"
			}
			return w + " weeks ago–Today"
		}
		if d.Months > 0 {
			return r.Start.Format("Jan 2") + "–Today"
		}

		return n(d.Days) + " days ago–Today"
	}

	if d.Months == 0 && d.Days == 0 {
		return r.Start.Format(addYear(r.Start, "Jan 2"))
	}

	return r.Start.Format(addYear(r.Start, "Jan 2")) + "–" + r.End.Format(addYear(r.End, "Jan 2")) +
		" (" + d.String() + ")"
}

// Diff gets the difference between two dates.
//
// Optionally pass any Period arguments to get the difference in those periods,
// ignoring any others. For example "Month, Day" would return "29 months, 6
// days", instead of "2 years, 5 months, 6 days". The default is to get
// everything excluding weeks.
//
// Adapted from https://stackoverflow.com/a/36531443/660921
func (r Range) Diff(periods ...Period) Diff {
	if r.Start.After(r.End) {
		r.Start, r.End = r.End, r.Start
	}

	y1, m1, d1 := r.Start.Date()
	y2, m2, d2 := r.End.Date()
	h1, min1, s1 := r.Start.Clock()
	h2, min2, s2 := r.End.Clock()

	d := Diff{
		Years: y2 - y1, Months: int(m2 - m1), Days: d2 - d1,
		Hours: h2 - h1, Mins: min2 - min1, Secs: s2 - s1,
	}

	if d.Secs < 0 {
		d.Secs += 60
		d.Mins--
	}
	if d.Mins < 0 {
		d.Mins += 60
		d.Hours--
	}
	if d.Hours < 0 {
		d.Hours += 24
		d.Days--
	}
	if d.Days < 0 {
		t := time.Date(y1, m1, 32, 0, 0, 0, 0, time.UTC)
		d.Days += 32 - t.Day()
		d.Months--
	}
	if d.Months < 0 {
		d.Months += 12
		d.Years--
	}

	if len(periods) == 0 {
		return d
	}

	var hasY, hasM, hasW, hasD, hasH, hasMin, hasSec bool
	for _, v := range periods {
		switch v {
		case Year:
			hasY = true
		case Month:
			hasM = true
		case WeekMonday, WeekSunday:
			hasW = true
		case Day:
			hasD = true
		case Hour:
			hasH = true
		case Minute:
			hasMin = true
		case Second:
			hasSec = true
		}
	}

	if !hasY {
		d.Months += d.Years * 12
		d.Years = 0
	}
	if !hasM {
		t := r.Start
		for ; d.Months > 0; d.Months-- {
			d.Days += DaysInMonth(t)
			t = Add(t, 1, Month)
		}
	}
	if hasW {
		d.Weeks = d.Days / 7
		d.Days = d.Days % 7
	}
	if !hasD {
		d.Hours += d.Days * 24
		d.Days = 0
	}
	if !hasH {
		d.Mins += d.Hours * 60
		d.Hours = 0
	}
	if !hasMin {
		d.Secs += d.Mins * 60
		d.Mins = 0
	}
	if !hasSec {
		d.Secs = 0
	}
	return d
}

// Diff represents the difference between two times.
type Diff struct {
	Years, Months, Weeks, Days int
	Hours, Mins, Secs          int
}

func (d Diff) String() string {
	n := strconv.Itoa

	if d.Months == 0 {
		if d.Days == 1 {
			return "2 days"
		}
		if d.Days == 6 {
			return "1 week"
		}
		// TODO?
		//if (d.Days-1)%7 == 0 {
		//	return n(d.Days/7) + " weeks"
		//}
		return n(d.Days+1) + " days"
	}

	s := n(d.Months) + " month"
	if d.Months > 1 {
		s += "s"
	}
	if d.Days > 0 {
		s += ", " + n(d.Days) + " days"
	}
	return s
}
