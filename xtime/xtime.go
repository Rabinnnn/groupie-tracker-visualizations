package xtime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse returns the time represented in the given string. The date in the string,
// is expected to follow the time format DD-MM-YYYY. Returns an error if parsing the time in the format fails
func Parse(s string) (time.Time, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("xtime: invalid format `%s`: expected format: DD-MM-YYYY", s)
	}

	dd, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	mm, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	yy, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(int(yy), time.Month(mm), int(dd), 0, 0, 0, 0, time.UTC), nil
}
