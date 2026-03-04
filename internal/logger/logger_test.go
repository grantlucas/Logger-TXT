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
