package entry

import "time"

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

// EndOfDay returns t with the time set to 23:59 in the same location.
func EndOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 0, 0, t.Location())
}
