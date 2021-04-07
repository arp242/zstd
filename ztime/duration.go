package ztime

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// DurationAs formats a duration as the given time unit, with fractions (if
// any).
//
// For example DurationAs(d, time.Millisecond) will return "0.05" for 50
// microseconds.
//
// Use Round() if you want to limit the precision.
func DurationAs(d, as time.Duration) string {
	f := float64(d) / float64(as)
	_, frac := math.Modf(f)
	if frac == 0 {
		return fmt.Sprintf("%.0f", f)
	}

	s := strings.TrimRight(fmt.Sprintf("%f", f), "0")
	if s == "0." {
		s = fmt.Sprintf("%.999f", f)
		for i, c := range s {
			if c != '0' && c != '.' {
				s = s[:i+1]
				break
			}
		}
	}
	return s
}
