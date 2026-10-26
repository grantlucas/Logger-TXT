package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	path := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = f.Close() })
	return f
}

func scanAll(t *testing.T, f *os.File) []string {
	t.Helper()
	s := NewReverseLineScanner(f)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("ReverseLineScanner error: %v", err)
	}
	return lines
}

func TestReverseLineScanner_EmptyFile(t *testing.T) {
	f := writeTempFile(t, "")
	got := scanAll(t, f)
	if len(got) != 0 {
		t.Fatalf("got %d lines for empty file, want 0", len(got))
	}
}

func TestReverseLineScanner_SingleLineNoNewline(t *testing.T) {
	f := writeTempFile(t, "only line")
	got := scanAll(t, f)
	if len(got) != 1 {
		t.Fatalf("got %d lines, want 1", len(got))
	}
	if got[0] != "only line" {
		t.Errorf("line = %q, want %q", got[0], "only line")
	}
}

func TestReverseLineScanner_ChunkBoundary(t *testing.T) {
	// Create content where lines span the 8KB chunk boundary
	line := strings.Repeat("x", 5000)
	content := line + "\n" + line + "\n" + line + "\n"
	f := writeTempFile(t, content)
	got := scanAll(t, f)

	if len(got) != 3 {
		t.Fatalf("got %d lines, want 3", len(got))
	}
	for i, g := range got {
		if g != line {
			t.Errorf("line[%d] length = %d, want %d", i, len(g), len(line))
		}
	}
}

func TestReverseLineScanner_WindowsLineEndings(t *testing.T) {
	f := writeTempFile(t, "line1\r\nline2\r\nline3\r\n")
	got := scanAll(t, f)

	want := []string{"line3", "line2", "line1"}
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestReverseLineScanner_ReadError(t *testing.T) {
	// Create a file, open it, then close it to force read errors
	path := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(path, []byte("line1\nline2\n"), 0644); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	s := NewReverseLineScanner(f)
	// Close the file to force ReadAt errors
	_ = f.Close()

	for s.Scan() {
		// drain
	}
	if s.Err() == nil {
		t.Fatal("expected error after closing file, got nil")
	}
}

func TestReverseLineScanner_MultipleLines(t *testing.T) {
	f := writeTempFile(t, "line1\nline2\nline3\n")
	got := scanAll(t, f)

	want := []string{"line3", "line2", "line1"}
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestReverseLineScanner_LineStart(t *testing.T) {
	long := strings.Repeat("x", 5000)
	tests := []struct {
		name    string
		content string
	}{
		{"LF", "first\nsecond\nthird\n"},
		{"no trailing newline", "first\nsecond\nthird"},
		{"CRLF", "first\r\nsecond\r\nthird\r\n"},
		{"blank lines", "first\n\nsecond\n\n\nthird\n"},
		{"chunk boundary", long + "\n" + long + "b\n" + long + "c\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewReverseLineScanner(writeTempFile(t, tt.content))
			n := 0
			for s.Scan() {
				n++
				if !strings.HasPrefix(tt.content[s.LineStart():], s.Text()) {
					t.Errorf("LineStart() = %d does not point at %q", s.LineStart(), s.Text())
				}
			}
			if n != 3 {
				t.Errorf("scanned %d lines, want 3", n)
			}
		})
	}
}
