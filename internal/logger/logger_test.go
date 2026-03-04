package logger_test

import (
	"os"
	"path/filepath"
	"testing"

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
