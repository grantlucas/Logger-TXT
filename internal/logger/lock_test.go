package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/grantlucas/Logger-TXT/internal/entry"
)

func TestWriters_FailWhenLockIsHeld(t *testing.T) {
	saved := lockTimeout
	lockTimeout = 50 * time.Millisecond
	t.Cleanup(func() { lockTimeout = saved })

	e := entry.Entry{Time: time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC), Message: "B"}
	tests := []struct {
		name string
		op   func(path string) error
	}{
		{"Append", func(p string) error { return Append(p, e) }},
		{"Insert", func(p string) error { return Insert(p, e) }},
		{"DeleteLast", func(p string) error { _, err := DeleteLast(p); return err }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			const content = "20/10/26 09:00 -0500 - A\n22/10/26 09:00 -0500 - C\n"
			path := filepath.Join(t.TempDir(), "log.txt")
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
			lock := path + ".lock"
			if err := os.WriteFile(lock, nil, 0644); err != nil {
				t.Fatal(err)
			}

			err := tt.op(path)
			if err == nil || !strings.Contains(err.Error(), lock) {
				t.Fatalf("error = %v, want one naming %s", err, lock)
			}
			if got, _ := os.ReadFile(path); string(got) != content {
				t.Errorf("file content = %q, want unchanged", got)
			}
			if _, err := os.Stat(lock); err != nil {
				t.Errorf("another process's lock file was removed: %v", err)
			}
		})
	}
}
