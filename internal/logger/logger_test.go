package logger_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/grantlucas/Logger-TXT/internal/entry"
	"github.com/grantlucas/Logger-TXT/internal/logger"
)

func TestEnsureFile_CreatesNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")

	err := logger.EnsureFile(path)
	if err != nil {
		t.Fatalf("EnsureFile() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file was not created: %v", err)
	}
}

func TestEnsureFile_PreservesExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	want := "existing content\n"
	os.WriteFile(path, []byte(want), 0644)

	err := logger.EnsureFile(path)
	if err != nil {
		t.Fatalf("EnsureFile() error = %v", err)
	}

	got, _ := os.ReadFile(path)
	if string(got) != want {
		t.Fatalf("file content changed: got %q, want %q", got, want)
	}
}

func TestEnsureFile_ErrorOnUnwritablePath(t *testing.T) {
	// /dev/null/impossible is never a valid directory
	path := filepath.Join("/dev/null", "impossible", "log.txt")

	err := logger.EnsureFile(path)
	if err == nil {
		t.Fatal("EnsureFile() expected error for unwritable path, got nil")
	}
}

func TestEnsureFile_ErrorOnReadOnlyDirectory(t *testing.T) {
	dir := t.TempDir()
	readOnly := filepath.Join(dir, "readonly")
	os.Mkdir(readOnly, 0555)
	path := filepath.Join(readOnly, "log.txt")

	err := logger.EnsureFile(path)
	if err == nil {
		t.Fatal("EnsureFile() expected error for read-only directory, got nil")
	}
}

func TestEnsureFile_CreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "dirs", "log.txt")

	err := logger.EnsureFile(path)
	if err != nil {
		t.Fatalf("EnsureFile() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file was not created: %v", err)
	}
}

func testTime() time.Time {
	return time.Date(2026, 2, 22, 10, 30, 0, 0, time.FixedZone("EST", -5*3600))
}

func TestAppend_WritesToNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	e := entry.Entry{Time: testTime(), Message: "Grabbed a coffee"}

	err := logger.Append(path, e)
	if err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	got, _ := os.ReadFile(path)
	want := "22/02/26 10:30 -0500 - Grabbed a coffee\n"
	if string(got) != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestAppend_AppendsToExistingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	e1 := entry.Entry{Time: testTime(), Message: "First entry"}
	e2 := entry.Entry{Time: testTime(), Type: "WORK", Message: "Second entry"}

	logger.Append(path, e1)
	logger.Append(path, e2)

	got, _ := os.ReadFile(path)
	want := "22/02/26 10:30 -0500 - First entry\n22/02/26 10:30 -0500 - WORK - Second entry\n"
	if string(got) != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestAppend_ErrorOnUnwritablePath(t *testing.T) {
	err := logger.Append("/dev/null/impossible/log.txt", entry.Entry{Message: "test"})
	if err == nil {
		t.Fatal("Append() expected error for unwritable path, got nil")
	}
}

func TestAppend_ErrorOnReadOnlyDirectory(t *testing.T) {
	dir := t.TempDir()
	readOnly := filepath.Join(dir, "readonly")
	os.Mkdir(readOnly, 0555)

	err := logger.Append(filepath.Join(readOnly, "log.txt"), entry.Entry{Message: "test"})
	if err == nil {
		t.Fatal("Append() expected error for read-only directory, got nil")
	}
}

// writeLines is a test helper that writes lines to a file.
func writeLines(t *testing.T, path string, lines []string) {
	t.Helper()
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestTail_ReturnsLastNLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{"line1", "line2", "line3", "line4", "line5"})

	got, err := logger.Tail(path, 3)
	if err != nil {
		t.Fatalf("Tail() error = %v", err)
	}

	want := []string{"line3", "line4", "line5"}
	if len(got) != len(want) {
		t.Fatalf("Tail() returned %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Tail()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestTail_ReturnsAllLinesWhenFewerThanN(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{"line1", "line2"})

	got, err := logger.Tail(path, 10)
	if err != nil {
		t.Fatalf("Tail() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Tail() returned %d lines, want 2", len(got))
	}
	if got[0] != "line1" || got[1] != "line2" {
		t.Errorf("Tail() = %v, want [line1, line2]", got)
	}
}

func TestTail_EmptyFileReturnsEmptySlice(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	os.WriteFile(path, []byte{}, 0644)

	got, err := logger.Tail(path, 10)
	if err != nil {
		t.Fatalf("Tail() error = %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("Tail() returned %d lines for empty file, want 0", len(got))
	}
}

func TestTail_HandlesWindowsLineEndings(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	os.WriteFile(path, []byte("line1\r\nline2\r\nline3\r\n"), 0644)

	got, err := logger.Tail(path, 2)
	if err != nil {
		t.Fatalf("Tail() error = %v", err)
	}

	want := []string{"line2", "line3"}
	if len(got) != len(want) {
		t.Fatalf("Tail() returned %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Tail()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestTail_ErrorOnScannerFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	// bufio.Scanner has a default max token size of 64KB.
	// A line exceeding that triggers a scan error.
	longLine := make([]byte, 1024*1024)
	for i := range longLine {
		longLine[i] = 'x'
	}
	os.WriteFile(path, longLine, 0644)

	_, err := logger.Tail(path, 10)
	if err == nil {
		t.Fatal("Tail() expected error for oversized line, got nil")
	}
}

func TestTail_ErrorOnNonExistentFile(t *testing.T) {
	_, err := logger.Tail("/nonexistent/log.txt", 10)
	if err == nil {
		t.Fatal("Tail() expected error for non-existent file, got nil")
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{
		"22/02/26 10:30 -0500 - WORK - Fixed bug",
		"22/02/26 10:31 -0500 - Grabbed coffee",
		"22/02/26 10:32 -0500 - WORK (API) - Deployed",
	})

	got, err := logger.Search(path, "work", false, 10)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Search() returned %d results, want 2", len(got))
	}
	if got[0] != "22/02/26 10:30 -0500 - WORK - Fixed bug" {
		t.Errorf("Search()[0] = %q", got[0])
	}
	if got[1] != "22/02/26 10:32 -0500 - WORK (API) - Deployed" {
		t.Errorf("Search()[1] = %q", got[1])
	}
}

func TestSearch_CaseSensitive(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{
		"22/02/26 10:30 -0500 - WORK - Fixed bug",
		"22/02/26 10:31 -0500 - work - lowercase",
		"22/02/26 10:32 -0500 - Grabbed coffee",
	})

	got, err := logger.Search(path, "WORK", true, 10)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("Search() returned %d results, want 1", len(got))
	}
	if got[0] != "22/02/26 10:30 -0500 - WORK - Fixed bug" {
		t.Errorf("Search()[0] = %q", got[0])
	}
}

func TestSearch_LimitsResults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{
		"22/02/26 10:30 -0500 - WORK - First",
		"22/02/26 10:31 -0500 - WORK - Second",
		"22/02/26 10:32 -0500 - WORK - Third",
	})

	got, err := logger.Search(path, "WORK", false, 2)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Search() returned %d results, want 2", len(got))
	}
	// Should return the last 2 matches
	if got[0] != "22/02/26 10:31 -0500 - WORK - Second" {
		t.Errorf("Search()[0] = %q", got[0])
	}
	if got[1] != "22/02/26 10:32 -0500 - WORK - Third" {
		t.Errorf("Search()[1] = %q", got[1])
	}
}

func TestSearch_NoMatchesReturnsEmptySlice(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{"22/02/26 10:30 -0500 - Grabbed coffee"})

	got, err := logger.Search(path, "NONEXISTENT", false, 10)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("Search() returned %d results, want 0", len(got))
	}
}

func TestSearch_ErrorOnScannerFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	longLine := make([]byte, 1024*1024)
	for i := range longLine {
		longLine[i] = 'x'
	}
	os.WriteFile(path, longLine, 0644)

	_, err := logger.Search(path, "x", false, 10)
	if err == nil {
		t.Fatal("Search() expected error for oversized line, got nil")
	}
}

func TestSearch_ErrorOnNonExistentFile(t *testing.T) {
	_, err := logger.Search("/nonexistent/log.txt", "term", false, 10)
	if err == nil {
		t.Fatal("Search() expected error for non-existent file, got nil")
	}
}
