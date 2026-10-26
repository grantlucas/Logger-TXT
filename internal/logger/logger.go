// Package logger provides core log file operations.
package logger

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
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
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	line := e.Format() + "\n"
	if size, _ := f.Seek(0, io.SeekEnd); endsWithoutNewline(f, size) {
		line = "\n" + line
	}
	_, err = f.WriteString(line)
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

// Insert writes the entry into the file in chronological position, just
// before the first entry whose timestamp is later than e.Time. Lines that are
// not entries stay with the entry above them.
func Insert(path string, e entry.Entry) error {
	path = resolvePath(path)
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return Append(path, e)
	}
	return rewrite(path, func(w io.Writer, src *os.File, size int64) error {
		return copyWithInsert(w, src, size, e)
	})
}

// copyWithInsert copies src to w with e inserted in chronological position.
func copyWithInsert(w io.Writer, src *os.File, size int64, e entry.Entry) error {
	insertAt, err := insertOffset(src, size, e.Time)
	line := e.Format() + "\n"
	if insertAt == size && endsWithoutNewline(src, size) {
		line = "\n" + line
	}
	if err == nil {
		_, err = io.Copy(w, io.NewSectionReader(src, 0, insertAt))
	}
	if err == nil {
		_, err = io.WriteString(w, line)
	}
	if err == nil {
		_, err = io.Copy(w, io.NewSectionReader(src, insertAt, size-insertAt))
	}
	return err
}

// insertOffset returns the byte offset of the first entry later than t, or
// size when there is none. It scans backwards, so a recent backfill reads
// only the end of the file.
func insertOffset(f *os.File, size int64, t time.Time) (int64, error) {
	insertAt := size
	s := NewReverseLineScanner(f)
	for s.Scan() {
		// Non-entry lines (headers, notes) stay with the entry above them
		existing, err := entry.ParseEntry(s.Text())
		if err != nil {
			continue
		}
		if !existing.Time.After(t) {
			break
		}
		insertAt = s.LineStart()
	}
	return insertAt, s.Err()
}

// rewrite atomically replaces the log at path with the output of write, which
// reads the current contents from src. The new contents go to a temp file in
// the same directory that is synced and renamed over the original, so a crash
// or full disk leaves the old log intact rather than half-rewritten. The
// original is opened read/write first so a read-only log is refused, not
// silently replaced.
func rewrite(path string, write func(w io.Writer, src *os.File, size int64) error) error {
	src, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	info, err := src.Stat()
	var tmp *os.File
	if err == nil {
		tmp, err = os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".*.tmp")
	}
	if err != nil {
		_ = src.Close()
		return err
	}

	err = write(tmp, src, info.Size())
	_ = src.Close() // before the rename, which Windows refuses on open files
	if err == nil {
		err = tmp.Sync()
	}
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Chmod(tmp.Name(), info.Mode().Perm())
	}
	if err == nil {
		err = os.Rename(tmp.Name(), path)
	}
	if err != nil {
		_ = os.Remove(tmp.Name())
	}
	return err
}

// resolvePath follows symlinks so a rewrite replaces the real log file
// rather than the link pointing at it.
func resolvePath(path string) string {
	if real, err := filepath.EvalSymlinks(path); err == nil {
		return real
	}
	return path
}

// endsWithoutNewline reports whether the file's last byte is something other
// than a newline, so a line written at the end would join onto the last line.
func endsWithoutNewline(f *os.File, size int64) bool {
	last := make([]byte, 1)
	_, err := f.ReadAt(last, size-1)
	return err == nil && last[0] != '\n'
}

// Tail returns the last n lines from the file. When filter is non-nil,
// each line is parsed and only entries matching the filter count toward n.
// Unparseable lines are skipped when a filter is active.
func Tail(path string, n int, filter func(entry.Entry) bool) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var lines []string
	s := NewReverseLineScanner(f)
	for s.Scan() {
		line := s.Text()
		if filter != nil {
			e, err := entry.ParseEntry(line)
			if err != nil {
				continue
			}
			if !filter(e) {
				continue
			}
		}
		lines = append(lines, line)
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
// When caseSensitive is false, comparison is case-insensitive. When filter
// is non-nil, lines must also match the entry filter to count toward limit.
func Search(path string, term string, caseSensitive bool, limit int, filter func(entry.Entry) bool) ([]string, error) {
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
			if filter != nil {
				e, err := entry.ParseEntry(line)
				if err != nil {
					continue
				}
				if !filter(e) {
					continue
				}
			}
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
	var deleted string
	err := rewrite(resolvePath(path), func(w io.Writer, src *os.File, _ int64) error {
		data, err := io.ReadAll(src)
		content := strings.TrimRight(string(data), "\r\n")
		if err == nil && content == "" {
			err = fmt.Errorf("file is empty")
		}
		if err != nil {
			return err
		}
		idx := strings.LastIndex(content, "\n")
		deleted = content[idx+1:]
		_, err = io.WriteString(w, content[:idx+1])
		return err
	})
	if err != nil {
		return "", err
	}
	return deleted, nil
}
