// Package logger provides core log file operations.
package logger

import (
	"bufio"
	"os"
	"path/filepath"

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
