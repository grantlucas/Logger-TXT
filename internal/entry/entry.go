// Package entry defines the log entry type and handles formatting and parsing.
package entry

import "time"

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
	if e.Type != "" {
		s += e.Type + " - "
	}
	s += e.Message
	return s
}
