package entry

import (
	"fmt"
	"time"
)

const (
	inputDateLayout     = "02/01/06"
	inputDateTimeLayout = "02/01/06 15:04"
)

// ParseInputTime parses a date string in DD/MM/YY or DD/MM/YY HH:MM format.
// Returns the parsed time, whether only a date was provided, and any error.
func ParseInputTime(input string, loc *time.Location) (time.Time, bool, error) {
	if t, err := time.ParseInLocation(inputDateTimeLayout, input, loc); err == nil {
		return t, false, nil
	}
	t, err := time.ParseInLocation(inputDateLayout, input, loc)
	if err != nil {
		return time.Time{}, false, err
	}
	return t, true, nil
}

// ParseDateRange parses start and end date strings into a time range.
// When end is date-only, EndOfDay is applied automatically.
func ParseDateRange(startStr, endStr string, loc *time.Location) (time.Time, time.Time, error) {
	start, _, err := ParseInputTime(startStr, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid --start value: %w", err)
	}

	end, endDateOnly, err := ParseInputTime(endStr, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid --end value: %w", err)
	}

	if endDateOnly {
		end = EndOfDay(end)
	}

	return start, end, nil
}

// EndOfDay returns t with the time set to 23:59 in the same location.
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 0, 0, t.Location())
}
