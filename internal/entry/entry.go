// Package entry defines the log entry type and handles formatting and parsing.
package entry

import (
	"fmt"
	"strings"
	"time"
)

// Entry represents a single log entry.
type Entry struct {
	Time    time.Time
	Type    string
	Project string
	Message string
}

// Format returns the log line string for this entry.
func (e Entry) Format() string {
	s := e.Time.Format("02/01/06 15:04 -0700") + " - "
	if e.Type != "" && e.Project != "" {
		s += e.Type + " (" + e.Project + ") - "
	} else if e.Type != "" {
		s += e.Type + " - "
	} else if e.Project != "" {
		s += "(" + e.Project + ") - "
	}
	s += e.Message
	return s
}

const timeLayout = "02/01/06 15:04 -0700"

// ParseEntry parses a log line into an Entry.
func ParseEntry(line string) (Entry, error) {
	// Timestamp is fixed-width: "DD/MM/YY HH:MM +ZZZZ" = 20 chars
	if len(line) < 20 {
		return Entry{}, fmt.Errorf("line too short: %q", line)
	}

	t, err := time.Parse(timeLayout, line[:20])
	if err != nil {
		return Entry{}, fmt.Errorf("invalid timestamp: %w", err)
	}

	// After timestamp, expect " - " separator
	rest := line[20:]
	if !strings.HasPrefix(rest, " - ") {
		return Entry{}, fmt.Errorf("missing separator after timestamp: %q", line)
	}
	rest = rest[3:] // skip " - "

	return Entry{
		Time:    t,
		Message: rest,
	}, nil
}
