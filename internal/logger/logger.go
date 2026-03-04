// Package logger provides core log file operations.
package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	defer f.Close()
	_, err = f.WriteString(e.Format() + "\n")
	return err
}

// Tail returns the last n lines from the file.
func Tail(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:]
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
	defer f.Close()

	if !caseSensitive {
		term = strings.ToLower(term)
	}

	var matches []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		haystack := line
		if !caseSensitive {
			haystack = strings.ToLower(line)
		}
		if strings.Contains(haystack, term) {
			matches = append(matches, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(matches) > limit {
		matches = matches[len(matches)-limit:]
	}
	return matches, nil
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
