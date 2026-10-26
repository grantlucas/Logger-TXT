package logger

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRewrite_FailureLeavesOriginal(t *testing.T) {
	const content = "original\n"
	tests := []struct {
		name  string
		setup func(t *testing.T, dir string)
		write func(w io.Writer, src *os.File, size int64) error
	}{
		{
			name:  "write fails",
			setup: func(t *testing.T, dir string) {},
			write: func(w io.Writer, src *os.File, size int64) error {
				_, _ = io.WriteString(w, "partial")
				return errors.New("boom")
			},
		},
		{
			name: "directory not writable",
			setup: func(t *testing.T, dir string) {
				if err := os.Chmod(dir, 0555); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { _ = os.Chmod(dir, 0755) })
			},
			write: func(w io.Writer, src *os.File, size int64) error { return nil },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "log.txt")
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				t.Fatal(err)
			}
			tt.setup(t, dir)

			if err := rewrite(path, tt.write); err == nil {
				t.Fatal("rewrite() expected error, got nil")
			}

			if got, _ := os.ReadFile(path); string(got) != content {
				t.Errorf("file content = %q, want unchanged %q", got, content)
			}
			if files, _ := os.ReadDir(dir); len(files) != 1 {
				t.Errorf("directory has %d files, want only the log", len(files))
			}
		})
	}
}
