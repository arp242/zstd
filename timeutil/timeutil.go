package timeutil // import "zgo.at/utils/timeutil"

import (
	"fmt"
	"time"
)

// UnixMilli returns the number of milliseconds elapsed since January 1, 1970
// UTC.
func UnixMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

var (
	day    = 24 * time.Hour
	day100 = 24 * time.Hour * 100
)

// FormatDuration formats duration as a concise string. It's less accurate than
// Duration.String(), but shorter.
func FormatDuration(d time.Duration) string {
	switch {
	case d >= day100:
		h := int(d.Round(time.Hour).Hours())
		return fmt.Sprintf("%dd", h/24)
	case d > day:
		h := int(d.Round(time.Hour).Hours())
		return fmt.Sprintf("%dd%dh", h/24, h%24)
	case d >= 10*time.Minute:
		s := d.Round(time.Minute).String()
		return s[:len(s)-2]
	}

	return d.Round(time.Second).String()
}
