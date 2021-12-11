package ztime

import (
	"strconv"
	"time"
)

// RangeLocale can be used to translate the output of Range.String().
//
// This defaults to DefaultRangeLocale.
type RangeLocale struct {
	Today     func() string             // "Today"
	Yesterday func() string             // "Yesterday"
	Month     func(m time.Month) string // "January", "December"
	DayAgo    func(n int) string        // "1 day ago", "5 days ago"
	WeekAgo   func(n int) string        // "1 week ago", "5 weeks ago"
	MonthAgo  func(n int) string        // "1 month ago", "5 months ago"
}

var DefaultRangeLocale = RangeLocale{
	Today:     func() string { return "Today" },
	Yesterday: func() string { return "Yesterday" },
	Month: func(m time.Month) string {
		return []string{"January", "February", "March", "April", "May", "June",
			"July", "August", "September", "October", "November", "December"}[m-1]
	},
	DayAgo: func(n int) string {
		if n == 1 {
			return "1 day ago"
		}
		return strconv.Itoa(n) + " days ago"
	},
	WeekAgo: func(n int) string {
		if n == 1 {
			return "1 week ago"
		}
		return strconv.Itoa(n) + " weeks ago"
	},
	MonthAgo: func(n int) string {
		if n == 1 {
			return "1 month ago"
		}
		return strconv.Itoa(n) + " months ago"
	},
}
