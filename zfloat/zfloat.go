// Package zfloat implements functions for floats.
package zfloat

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Join a float with the given separator.
func Join(nums []float64, sep string) string {
	s := make([]string, len(nums))
	for i := range nums {
		s[i] = strconv.FormatFloat(nums[i], 'f', -1, 64)
	}
	return strings.Join(s, sep)
}

// Split converts a string of numbers to a []float64.
func Split(s, sep string) ([]float64, error) {
	s = strings.Trim(s, " \t\n"+sep)
	if len(s) == 0 {
		return nil, nil
	}

	items := strings.Split(s, sep)
	ret := make([]float64, len(items))
	for i := range items {
		val, err := strconv.ParseFloat(strings.TrimSpace(items[i]), 64)
		if err != nil {
			return nil, err
		}
		ret[i] = val
	}

	return ret, nil
}

// RoundTo rounds the value to the given precision.
//
// e.g. RoundTo(7.258, 2) will return 7.26
func RoundTo(f float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return math.Round(f*shift) / shift
}

// Limit a value between a lower and upper limit.
func Limit(v, lower, upper float64) float64 {
	return math.Max(math.Min(v, upper), lower)
}

// IsSignedZero checks if this number is a signed zero (i.e. -0, instead of +0).
func IsSignedZero(f float64) bool {
	return math.Float64bits(f)^uint64(1<<63) == 0
}

// Byte is a float64 where the String() method prints out a human-redable
// description.
type Byte float64

var units = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}

func (b Byte) String() string {
	i := 0
	for ; i < len(units); i++ {
		if b < 1024 {
			return fmt.Sprintf("%.1f%s", b, units[i])
		}
		b /= 1024
	}
	return fmt.Sprintf("%.1f%s", b*1024, units[i-1])
}
