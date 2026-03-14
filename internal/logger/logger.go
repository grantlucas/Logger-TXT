// Package logger provides core log file operations.
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grantlucas/Logger-TXT/internal/entry"
)

// EnsureFile creates the log file and any parent directories if they don't exist.
func EnsureFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Append formats the entry and appends it as a new line to the file.
// Creates the file and parent directories if they don't exist.
func Append(path string, e entry.Entry) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	_, err = f.WriteString(e.Format() + "\n")
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

// Tail returns the last n lines from the file.
func Tail(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var lines []string
	s := NewReverseLineScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
		if len(lines) == n {
			break
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	// Reverse to chronological order
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
	return lines, nil
}

// Search returns the last limit lines that contain the search term.
// When caseSensitive is false, comparison is case-insensitive.
func Search(path string, term string, caseSensitive bool, limit int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	if !caseSensitive {
		term = strings.ToLower(term)
	}

	var matches []string
	s := NewReverseLineScanner(f)
	for s.Scan() {
		line := s.Text()
		haystack := line
		if !caseSensitive {
			haystack = strings.ToLower(line)
		}
		if strings.Contains(haystack, term) {
			matches = append(matches, line)
			if len(matches) == limit {
				break
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	// Reverse to chronological order
	for i, j := 0, len(matches)-1; i < j; i, j = i+1, j-1 {
		matches[i], matches[j] = matches[j], matches[i]
	}
	return matches, nil
}

// Range returns log lines whose timestamps fall within [start, end].
// Lines are returned in chronological order. If fn is non-nil, only entries
// for which fn returns true are included. Unparseable lines are skipped.
func Range(path string, start, end time.Time, fn func(entry.Entry) bool) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return scanRange(NewReverseLineScanner(f), start, end, fn)
}

func scanRange(s *ReverseLineScanner, start, end time.Time, fn func(entry.Entry) bool) ([]string, error) {
	var collected []string
	for s.Scan() {
		line := s.Text()
		e, err := entry.ParseEntry(line)
		if err != nil {
			continue // skip unparseable lines
		}
		if e.Time.After(end) {
			continue // haven't reached the window yet
		}
		if e.Time.Before(start) {
			break // past the window — done
		}
		if fn != nil && !fn(e) {
			continue
		}
		collected = append(collected, line)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	// Reverse to chronological order
	for i, j := 0, len(collected)-1; i < j; i, j = i+1, j-1 {
		collected[i], collected[j] = collected[j], collected[i]
	}
	return collected, nil
}

// DeleteLast removes the last line from the file and returns it.
func DeleteLast(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	content := strings.TrimRight(string(data), "\r\n")
	if content == "" {
		return "", fmt.Errorf("file is empty")
	}

	idx := strings.LastIndex(content, "\n")
	var deleted string
	var remaining string
	if idx < 0 {
		// Only one line in the file
		deleted = content
		remaining = ""
	} else {
		deleted = content[idx+1:]
		remaining = content[:idx+1]
	}

	if err := os.WriteFile(path, []byte(remaining), 0644); err != nil {
		return "", err
	}
	return deleted, nil
}
