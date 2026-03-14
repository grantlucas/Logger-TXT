package logger

import (
	"io"
	"os"
	"strings"
)

const chunkSize = 8192

// ReverseLineScanner reads a file from the end, yielding lines from newest
// to oldest. It mirrors the bufio.Scanner API.
type ReverseLineScanner struct {
	f      *os.File
	offset int64
	buf    string
	line   string
	err    error
	done   bool
}

// NewReverseLineScanner returns a scanner that reads f from end to start.
func NewReverseLineScanner(f *os.File) *ReverseLineScanner {
	offset, _ := f.Seek(0, io.SeekEnd)
	return &ReverseLineScanner{f: f, offset: offset}
}

// Scan advances to the next line (moving backwards through the file).
// Returns false when there are no more lines or an error occurred.
func (s *ReverseLineScanner) Scan() bool {
	for {
		if s.done && s.buf == "" {
			return false
		}

		// Try to extract a line from the buffer
		if idx := strings.LastIndex(s.buf, "\n"); idx >= 0 {
			s.line = strings.TrimRight(s.buf[idx+1:], "\r")
			s.buf = s.buf[:idx]
			// Skip empty lines from trailing newlines
			if s.line == "" {
				continue
			}
			return true
		}

		// No newline in buffer — need to read more
		if s.done {
			// No more data to read; remaining buffer is the first line
			s.line = strings.TrimRight(s.buf, "\r")
			s.buf = ""
			return s.line != ""
		}

		// Read the next chunk backwards
		readSize := int64(chunkSize)
		if readSize > s.offset {
			readSize = s.offset
		}
		s.offset -= readSize

		chunk := make([]byte, readSize)
		_, err := s.f.ReadAt(chunk, s.offset)
		if err != nil && err != io.EOF {
			s.err = err
			return false
		}

		s.buf = string(chunk) + s.buf

		if s.offset == 0 {
			s.done = true
		}
	}
}

// Text returns the most recent line produced by Scan.
func (s *ReverseLineScanner) Text() string {
	return s.line
}

// Err returns the first non-EOF error encountered.
func (s *ReverseLineScanner) Err() error {
	return s.err
}
