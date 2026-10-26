package logger

import (
	"fmt"
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

func TestWriters_DanglingSymlinkSharesTargetLock(t *testing.T) {
	saved := lockTimeout
	lockTimeout = 50 * time.Millisecond
	t.Cleanup(func() { lockTimeout = saved })

	e := entry.Entry{Time: time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC), Message: "B"}
	tests := []struct {
		name string
		hops int
		op   func(path string) error
	}{
		{"Append through one link", 1, func(p string) error { return Append(p, e) }},
		{"Append through a chain of links", 2, func(p string) error { return Append(p, e) }},
		{"Insert through one link", 1, func(p string) error { return Insert(p, e) }},
		{"Insert through a chain of links", 2, func(p string) error { return Insert(p, e) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			// The target does not exist yet, so every link is dangling
			link := filepath.Join(dir, "log.txt")
			for i := range tt.hops {
				next := filepath.Join(dir, fmt.Sprintf("link%d.txt", i))
				if err := os.Symlink(filepath.Base(link), next); err != nil {
					t.Fatal(err)
				}
				link = next
			}
			lock := filepath.Join(dir, "log.txt.lock")
			if err := os.WriteFile(lock, nil, 0644); err != nil {
				t.Fatal(err)
			}

			err = tt.op(link)
			if err == nil || !strings.Contains(err.Error(), lock) {
				t.Fatalf("error = %v, want one naming the target's lock %s", err, lock)
			}
		})
	}
}

func TestInsert_IntoMissingFileWaitsForLockBeforeChoosingPlacement(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.txt")
	lock := path + ".lock"
	if err := os.WriteFile(lock, nil, 0644); err != nil {
		t.Fatal(err)
	}

	errc := make(chan error)
	go func() {
		errc <- Insert(path, entry.Entry{Time: time.Date(2026, 10, 21, 10, 0, 0, 0, time.UTC), Message: "earlier"})
	}()
	// Let Insert reach the lock, then write a later entry as the lock holder
	time.Sleep(50 * time.Millisecond)
	const later = "21/10/26 11:00 +0000 - later\n"
	if err := os.WriteFile(path, []byte(later), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(lock); err != nil {
		t.Fatal(err)
	}

	if err := <-errc; err != nil {
		t.Fatalf("Insert() error = %v", err)
	}
	want := "21/10/26 10:00 +0000 - earlier\n" + later
	if got, _ := os.ReadFile(path); string(got) != want {
		t.Errorf("file content = %q, want %q", got, want)
	}
}
