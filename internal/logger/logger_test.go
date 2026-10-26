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
	if err := os.WriteFile(path, []byte(want), 0644); err != nil {
		t.Fatal(err)
	}

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
	if err := os.Mkdir(readOnly, 0555); err != nil {
		t.Fatal(err)
	}
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

	if err := logger.Append(path, e1); err != nil {
		t.Fatal(err)
	}
	if err := logger.Append(path, e2); err != nil {
		t.Fatal(err)
	}

	got, _ := os.ReadFile(path)
	want := "22/02/26 10:30 -0500 - First entry\n22/02/26 10:30 -0500 - WORK - Second entry\n"
	if string(got) != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}

func TestAppend_StartsNewLineWhenFileLacksTrailingNewline(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.txt")
	if err := os.WriteFile(path, []byte("22/02/26 09:00 -0500 - Edited by hand"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := logger.Append(path, entry.Entry{Time: testTime(), Message: "Next"}); err != nil {
		t.Fatalf("Append() error = %v", err)
	}

	got, _ := os.ReadFile(path)
	want := "22/02/26 09:00 -0500 - Edited by hand\n22/02/26 10:30 -0500 - Next\n"
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
	if err := os.Mkdir(readOnly, 0555); err != nil {
		t.Fatal(err)
	}

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

	got, err := logger.Tail(path, 3, nil)
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

	got, err := logger.Tail(path, 10, nil)
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
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := logger.Tail(path, 10, nil)
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
	if err := os.WriteFile(path, []byte("line1\r\nline2\r\nline3\r\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := logger.Tail(path, 2, nil)
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

func TestTail_ErrorOnReadFailure(t *testing.T) {
	dir := t.TempDir()
	// A directory path is not a readable file — triggers a read error.
	_, err := logger.Tail(dir, 10, nil)
	if err == nil {
		t.Fatal("Tail() expected error for directory path, got nil")
	}
}

func TestTail_ErrorOnNonExistentFile(t *testing.T) {
	_, err := logger.Tail("/nonexistent/log.txt", 10, nil)
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

	got, err := logger.Search(path, "work", false, 10, nil)
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

	got, err := logger.Search(path, "WORK", true, 10, nil)
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

	got, err := logger.Search(path, "WORK", false, 2, nil)
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

	got, err := logger.Search(path, "NONEXISTENT", false, 10, nil)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("Search() returned %d results, want 0", len(got))
	}
}

func TestSearch_ErrorOnReadFailure(t *testing.T) {
	dir := t.TempDir()
	// A directory path is not a readable file — triggers a read error.
	_, err := logger.Search(dir, "x", false, 10, nil)
	if err == nil {
		t.Fatal("Search() expected error for directory path, got nil")
	}
}

func TestSearch_ErrorOnNonExistentFile(t *testing.T) {
	_, err := logger.Search("/nonexistent/log.txt", "term", false, 10, nil)
	if err == nil {
		t.Fatal("Search() expected error for non-existent file, got nil")
	}
}

func TestDeleteLast_RemovesAndReturnsLastLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{"line1", "line2", "line3"})

	got, err := logger.DeleteLast(path)
	if err != nil {
		t.Fatalf("DeleteLast() error = %v", err)
	}
	if got != "line3" {
		t.Errorf("DeleteLast() = %q, want %q", got, "line3")
	}

	// Verify remaining file content
	remaining, _ := os.ReadFile(path)
	want := "line1\nline2\n"
	if string(remaining) != want {
		t.Errorf("remaining file = %q, want %q", remaining, want)
	}
}

func TestDeleteLast_SingleLineFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{"only line"})

	got, err := logger.DeleteLast(path)
	if err != nil {
		t.Fatalf("DeleteLast() error = %v", err)
	}
	if got != "only line" {
		t.Errorf("DeleteLast() = %q, want %q", got, "only line")
	}

	remaining, _ := os.ReadFile(path)
	if string(remaining) != "" {
		t.Errorf("remaining file = %q, want empty", remaining)
	}
}

func TestDeleteLast_ErrorOnEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	_, err := logger.DeleteLast(path)
	if err == nil {
		t.Fatal("DeleteLast() expected error for empty file, got nil")
	}
}

func TestDeleteLast_ErrorOnWriteFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{"line1", "line2"})
	// Make file read-only so WriteFile fails
	if err := os.Chmod(path, 0444); err != nil {
		t.Fatal(err)
	}

	_, err := logger.DeleteLast(path)
	if err == nil {
		t.Fatal("DeleteLast() expected error for read-only file, got nil")
	}
}

func TestDeleteLast_ErrorOnNonExistentFile(t *testing.T) {
	_, err := logger.DeleteLast("/nonexistent/log.txt")
	if err == nil {
		t.Fatal("DeleteLast() expected error for non-existent file, got nil")
	}
}

func TestRange_SkipsUnparseableLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	tz := time.FixedZone("EST", -5*3600)
	writeLines(t, path, []string{
		"22/02/26 10:00 -0500 - Valid entry",
		"this is not a valid log line",
		"22/02/26 11:00 -0500 - Another valid",
	})

	start := time.Date(2026, 2, 22, 0, 0, 0, 0, tz)
	end := time.Date(2026, 2, 22, 23, 59, 0, 0, tz)

	got, err := logger.Range(path, start, end, nil)
	if err != nil {
		t.Fatalf("Range() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("Range() returned %d lines, want 2 (unparseable should be skipped)", len(got))
	}
}

func TestRange_ErrorOnMissingFile(t *testing.T) {
	tz := time.FixedZone("EST", -5*3600)
	start := time.Date(2026, 2, 22, 0, 0, 0, 0, tz)
	end := time.Date(2026, 2, 22, 23, 59, 0, 0, tz)

	_, err := logger.Range("/nonexistent/log.txt", start, end, nil)
	if err == nil {
		t.Fatal("Range() expected error for missing file, got nil")
	}
}

func TestRange_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	if err := os.WriteFile(path, []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	tz := time.FixedZone("EST", -5*3600)
	start := time.Date(2026, 2, 22, 0, 0, 0, 0, tz)
	end := time.Date(2026, 2, 22, 23, 59, 0, 0, tz)

	got, err := logger.Range(path, start, end, nil)
	if err != nil {
		t.Fatalf("Range() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Range() returned %d lines for empty file, want 0", len(got))
	}
}

func TestRange_FilterFunction(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	tz := time.FixedZone("EST", -5*3600)
	writeLines(t, path, []string{
		"22/02/26 10:00 -0500 - WORK - Task one",
		"22/02/26 11:00 -0500 - Coffee break",
		"22/02/26 12:00 -0500 - WORK - Task two",
	})

	start := time.Date(2026, 2, 22, 0, 0, 0, 0, tz)
	end := time.Date(2026, 2, 22, 23, 59, 0, 0, tz)

	got, err := logger.Range(path, start, end, func(e entry.Entry) bool {
		return e.Type == "WORK"
	})
	if err != nil {
		t.Fatalf("Range() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("Range() returned %d lines, want 2", len(got))
	}
}

func TestRange_BoundariesInclusive(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	tz := time.FixedZone("EST", -5*3600)
	writeLines(t, path, []string{
		"22/02/26 09:00 -0500 - At start",
		"22/02/26 10:00 -0500 - Middle",
		"22/02/26 17:00 -0500 - At end",
	})

	start := time.Date(2026, 2, 22, 9, 0, 0, 0, tz)
	end := time.Date(2026, 2, 22, 17, 0, 0, 0, tz)

	got, err := logger.Range(path, start, end, nil)
	if err != nil {
		t.Fatalf("Range() error = %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("Range() returned %d lines, want 3 (boundaries should be inclusive)", len(got))
	}
}

func TestRange_ReturnsEntriesInChronologicalOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	tz := time.FixedZone("EST", -5*3600)
	writeLines(t, path, []string{
		"20/02/26 09:00 -0500 - Too early",
		"22/02/26 10:00 -0500 - In range first",
		"22/02/26 11:00 -0500 - In range second",
		"25/02/26 09:00 -0500 - Too late",
	})

	start := time.Date(2026, 2, 22, 0, 0, 0, 0, tz)
	end := time.Date(2026, 2, 22, 23, 59, 0, 0, tz)

	got, err := logger.Range(path, start, end, nil)
	if err != nil {
		t.Fatalf("Range() error = %v", err)
	}

	want := []string{
		"22/02/26 10:00 -0500 - In range first",
		"22/02/26 11:00 -0500 - In range second",
	}
	if len(got) != len(want) {
		t.Fatalf("Range() returned %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Range()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestTail_WithFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{
		"22/02/26 10:00 -0500 - WORK - Task one",
		"22/02/26 10:30 -0500 - MEETING - Standup",
		"22/02/26 11:00 -0500 - WORK (API) - Task two",
		"22/02/26 11:30 -0500 - MEETING - Review",
		"22/02/26 12:00 -0500 - WORK - Task three",
	})

	filter := func(e entry.Entry) bool { return e.Type == "WORK" }

	got, err := logger.Tail(path, 2, filter)
	if err != nil {
		t.Fatalf("Tail() error = %v", err)
	}

	want := []string{
		"22/02/26 11:00 -0500 - WORK (API) - Task two",
		"22/02/26 12:00 -0500 - WORK - Task three",
	}
	if len(got) != len(want) {
		t.Fatalf("Tail() returned %d lines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Tail()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestTail_WithFilterSkipsUnparseable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{
		"22/02/26 10:00 -0500 - WORK - Valid",
		"this is not a valid log line",
		"22/02/26 11:00 -0500 - WORK - Also valid",
	})

	filter := func(e entry.Entry) bool { return e.Type == "WORK" }

	got, err := logger.Tail(path, 10, filter)
	if err != nil {
		t.Fatalf("Tail() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Tail() returned %d lines, want 2 (unparseable skipped)", len(got))
	}
}

func TestSearch_WithFilter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{
		"22/02/26 10:00 -0500 - WORK - Fixed bug in API",
		"22/02/26 10:30 -0500 - MEETING - Discussed bug",
		"22/02/26 11:00 -0500 - WORK (API) - Another bug fix",
		"22/02/26 11:30 -0500 - MEETING - Bug triage",
	})

	filter := func(e entry.Entry) bool { return e.Type == "WORK" }

	got, err := logger.Search(path, "bug", false, 10, filter)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	want := []string{
		"22/02/26 10:00 -0500 - WORK - Fixed bug in API",
		"22/02/26 11:00 -0500 - WORK (API) - Another bug fix",
	}
	if len(got) != len(want) {
		t.Fatalf("Search() returned %d results, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Search()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestSearch_WithFilterRespectsLimit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{
		"22/02/26 10:00 -0500 - WORK - Bug one",
		"22/02/26 10:30 -0500 - WORK - Bug two",
		"22/02/26 11:00 -0500 - WORK - Bug three",
	})

	filter := func(e entry.Entry) bool { return e.Type == "WORK" }

	got, err := logger.Search(path, "bug", false, 2, filter)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Search() returned %d results, want 2", len(got))
	}
	if got[0] != "22/02/26 10:30 -0500 - WORK - Bug two" {
		t.Errorf("Search()[0] = %q, want last 2 matches", got[0])
	}
}

func TestSearch_WithFilterSkipsUnparseable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")
	writeLines(t, path, []string{
		"22/02/26 10:00 -0500 - WORK - Found bug",
		"this line mentions bug but is not parseable",
		"22/02/26 11:00 -0500 - WORK - Another bug",
	})

	filter := func(e entry.Entry) bool { return e.Type == "WORK" }

	got, err := logger.Search(path, "bug", false, 10, filter)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Search() returned %d results, want 2 (unparseable skipped)", len(got))
	}
}

func TestInsert(t *testing.T) {
	est := time.FixedZone("EST", -5*3600)
	at := time.Date(2026, 10, 21, 12, 0, 0, 0, est)
	e := entry.Entry{Time: at, Message: "B"}
	tests := []struct {
		name    string
		content string
		noFile  bool
		want    string
	}{
		{
			name:    "between existing entries",
			content: "20/10/26 09:00 -0500 - A\n22/10/26 09:00 -0500 - C\n",
			want:    "20/10/26 09:00 -0500 - A\n21/10/26 12:00 -0500 - B\n22/10/26 09:00 -0500 - C\n",
		},
		{
			name:    "before all entries",
			content: "22/10/26 09:00 -0500 - C\n23/10/26 09:00 -0500 - D\n",
			want:    "21/10/26 12:00 -0500 - B\n22/10/26 09:00 -0500 - C\n23/10/26 09:00 -0500 - D\n",
		},
		{
			name:    "after all entries",
			content: "20/10/26 09:00 -0500 - A\n",
			want:    "20/10/26 09:00 -0500 - A\n21/10/26 12:00 -0500 - B\n",
		},
		{
			name:    "after entries in the same minute",
			content: "21/10/26 12:00 -0500 - A1\n21/10/26 12:00 -0500 - A2\n22/10/26 09:00 -0500 - C\n",
			want:    "21/10/26 12:00 -0500 - A1\n21/10/26 12:00 -0500 - A2\n21/10/26 12:00 -0500 - B\n22/10/26 09:00 -0500 - C\n",
		},
		{
			name:    "keeps non-entry lines with the entry above them",
			content: "20/10/26 09:00 -0500 - A\nnote about A\n22/10/26 09:00 -0500 - C\n",
			want:    "20/10/26 09:00 -0500 - A\nnote about A\n21/10/26 12:00 -0500 - B\n22/10/26 09:00 -0500 - C\n",
		},
		{
			name:    "keeps a header above the earliest entry",
			content: "# Work log\n22/10/26 09:00 -0500 - C\n",
			want:    "# Work log\n21/10/26 12:00 -0500 - B\n22/10/26 09:00 -0500 - C\n",
		},
		{
			// Blank lines belong to the entry above, same as mid-file
			name:    "after trailing blank lines",
			content: "20/10/26 09:00 -0500 - A\n\n\n",
			want:    "20/10/26 09:00 -0500 - A\n\n\n21/10/26 12:00 -0500 - B\n",
		},
		{
			name:    "last line has no trailing newline",
			content: "20/10/26 09:00 -0500 - A",
			want:    "20/10/26 09:00 -0500 - A\n21/10/26 12:00 -0500 - B\n",
		},
		{
			name:    "CRLF line endings",
			content: "20/10/26 09:00 -0500 - A\r\n22/10/26 09:00 -0500 - C\r\n",
			want:    "20/10/26 09:00 -0500 - A\r\n21/10/26 12:00 -0500 - B\n22/10/26 09:00 -0500 - C\r\n",
		},
		{
			// 12:30 -0400 is 16:30 UTC, earlier than 12:00 -0500 (17:00 UTC)
			name:    "compares instants across offsets",
			content: "21/10/26 12:30 -0400 - A\n21/10/26 13:30 -0400 - C\n",
			want:    "21/10/26 12:30 -0400 - A\n21/10/26 12:00 -0500 - B\n21/10/26 13:30 -0400 - C\n",
		},
		{
			name:    "empty file",
			content: "",
			want:    "21/10/26 12:00 -0500 - B\n",
		},
		{
			name:   "missing file",
			noFile: true,
			want:   "21/10/26 12:00 -0500 - B\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "log.txt")
			if !tt.noFile {
				if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			if err := logger.Insert(path, e); err != nil {
				t.Fatalf("Insert() error = %v", err)
			}

			got, _ := os.ReadFile(path)
			if string(got) != tt.want {
				t.Errorf("file content = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestInsert_ErrorOnUnwritablePath(t *testing.T) {
	e := entry.Entry{Time: testTime(), Message: "test"}
	tests := []struct {
		name string
		path func(t *testing.T) string
	}{
		{"directory", func(t *testing.T) string { return t.TempDir() }},
		{"read-only file", func(t *testing.T) string {
			path := filepath.Join(t.TempDir(), "log.txt")
			if err := os.WriteFile(path, []byte("line1\n"), 0444); err != nil {
				t.Fatal(err)
			}
			return path
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := logger.Insert(tt.path(t), e); err == nil {
				t.Fatal("Insert() expected error, got nil")
			}
		})
	}
}

func TestWriters_PreserveLogFile(t *testing.T) {
	e := entry.Entry{Time: time.Date(2026, 10, 21, 12, 0, 0, 0, time.FixedZone("EST", -5*3600)), Message: "B"}
	tests := []struct {
		name string
		op   func(path string) error
	}{
		{"Append", func(p string) error { return logger.Append(p, e) }},
		{"Insert", func(p string) error { return logger.Insert(p, e) }},
		{"DeleteLast", func(p string) error { _, err := logger.DeleteLast(p); return err }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			target := filepath.Join(dir, "real.txt")
			content := "20/10/26 09:00 -0500 - A\n22/10/26 09:00 -0500 - C\n"
			if err := os.WriteFile(target, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
			link := filepath.Join(dir, "log.txt")
			if err := os.Symlink(target, link); err != nil {
				t.Fatal(err)
			}

			if err := tt.op(link); err != nil {
				t.Fatalf("%s() error = %v", tt.name, err)
			}

			if fi, err := os.Lstat(link); err != nil || fi.Mode()&os.ModeSymlink == 0 {
				t.Errorf("log path is no longer a symlink (err %v)", err)
			}
			if got, _ := os.ReadFile(target); string(got) == content {
				t.Error("symlink target was not updated")
			}
			if fi, _ := os.Stat(target); fi.Mode().Perm() != 0600 {
				t.Errorf("permissions = %v, want 0600", fi.Mode().Perm())
			}
			if files, _ := os.ReadDir(dir); len(files) != 2 {
				t.Errorf("directory has %d files, want only the log and its symlink", len(files))
			}
		})
	}
}
