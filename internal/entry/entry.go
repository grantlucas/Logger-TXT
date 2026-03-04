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

	var entryType, project, message string

	// Try to parse category prefix before the message.
	// Possible patterns:
	//   "TYPE (PROJECT) - message"
	//   "TYPE - message"
	//   "(PROJECT) - message"
	//   "message"
	if idx := strings.Index(rest, " - "); idx >= 0 {
		prefix := rest[:idx]
		after := rest[idx+3:]

		if parseCategory(prefix, &entryType, &project) {
			message = after
		} else {
			// No valid category — entire rest is the message
			message = rest
		}
	} else {
		message = rest
	}

	return Entry{
		Time:    t,
		Type:    entryType,
		Project: project,
		Message: message,
	}, nil
}

// parseCategory attempts to parse a category prefix like "TYPE", "(PROJECT)",
// or "TYPE (PROJECT)". Returns true if the prefix is a valid category.
func parseCategory(prefix string, entryType, project *string) bool {
	// "TYPE (PROJECT)"
	if parenOpen := strings.Index(prefix, "("); parenOpen >= 0 {
		if !strings.HasSuffix(prefix, ")") {
			return false
		}
		*project = prefix[parenOpen+1 : len(prefix)-1]
		if parenOpen > 0 {
			*entryType = strings.TrimSpace(prefix[:parenOpen])
		}
		return true
	}

	// "TYPE" — must be all uppercase letters (no spaces, no parens)
	if isTypeName(prefix) {
		*entryType = prefix
		return true
	}

	return false
}

// isTypeName returns true if s looks like a TYPE token (uppercase letters only).
func isTypeName(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
