package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestScanRange_ScannerError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "log.txt")

	// File must be large enough that the scanner needs to ReadAt
	line := "22/02/26 10:00 -0500 - Entry\n"
	content := ""
	for i := 0; i < 500; i++ {
		content += line
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Create scanner while file is open (records correct offset)
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	s := NewReverseLineScanner(f)

	// Close fd so ReadAt fails with EBADF
	_ = f.Close()

	tz := time.FixedZone("EST", -5*3600)
	start := time.Date(2026, 2, 22, 0, 0, 0, 0, tz)
	end := time.Date(2026, 2, 22, 23, 59, 0, 0, tz)

	_, err = scanRange(s, start, end, nil)
	if err == nil {
		t.Fatal("expected error from scanRange with broken scanner, got nil")
	}
}
