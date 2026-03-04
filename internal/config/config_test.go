package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFilePath_DefaultWhenNothingSet(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "")
	t.Chdir(t.TempDir())

	got := ResolveFilePath("")
	want := "./log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"\") = %q, want %q", got, want)
	}
}

func TestResolveFilePath_EnvVarOverridesDefault(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "/home/user/my-log.txt")
	t.Chdir(t.TempDir())

	got := ResolveFilePath("")
	want := "/home/user/my-log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"\") = %q, want %q", got, want)
	}
}

func TestResolveFilePath_LocalFileOverridesEnvVar(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "/home/user/env-log.txt")
	dir := t.TempDir()
	t.Chdir(dir)

	// Create a local log.txt so it exists in cwd
	if err := os.WriteFile(filepath.Join(dir, "log.txt"), []byte("existing log\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got := ResolveFilePath("")
	want := "./log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"\") = %q, want %q", got, want)
	}
}

func TestResolveFilePath_FlagOverridesEnvVar(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "/home/user/env-log.txt")
	t.Chdir(t.TempDir())

	got := ResolveFilePath("/tmp/flag-log.txt")
	want := "/tmp/flag-log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"/tmp/flag-log.txt\") = %q, want %q", got, want)
	}
}

func TestResolveFilePath_FlagOverridesLocalFile(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "")
	dir := t.TempDir()
	t.Chdir(dir)

	// Local log.txt exists, but flag should still win
	if err := os.WriteFile(filepath.Join(dir, "log.txt"), []byte("existing log\n"), 0644); err != nil {
		t.Fatal(err)
	}

	got := ResolveFilePath("/tmp/flag-log.txt")
	want := "/tmp/flag-log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"/tmp/flag-log.txt\") = %q, want %q", got, want)
	}
}

func TestResolveFilePath_FlagOverridesDefault(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "")
	t.Chdir(t.TempDir())

	got := ResolveFilePath("/tmp/flag-log.txt")
	want := "/tmp/flag-log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"/tmp/flag-log.txt\") = %q, want %q", got, want)
	}
}
