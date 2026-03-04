// Package entry defines the log entry type and handles formatting and parsing.
package entry

// Entry represents a single log entry.
type Entry struct {
	Time    string
	Type    string
	Project string
	Message string
}
