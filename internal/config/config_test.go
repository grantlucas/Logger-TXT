package config

import (
	"testing"
)

func TestResolveFilePath_DefaultWhenNothingSet(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "")

	got := ResolveFilePath("")
	want := "./log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"\") = %q, want %q", got, want)
	}
}

func TestResolveFilePath_EnvVarOverridesDefault(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "/home/user/my-log.txt")

	got := ResolveFilePath("")
	want := "/home/user/my-log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"\") = %q, want %q", got, want)
	}
}

func TestResolveFilePath_FlagOverridesEnvVar(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "/home/user/env-log.txt")

	got := ResolveFilePath("/tmp/flag-log.txt")
	want := "/tmp/flag-log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"/tmp/flag-log.txt\") = %q, want %q", got, want)
	}
}

func TestResolveFilePath_FlagOverridesDefault(t *testing.T) {
	t.Setenv("LOGGERTXT_PATH", "")

	got := ResolveFilePath("/tmp/flag-log.txt")
	want := "/tmp/flag-log.txt"

	if got != want {
		t.Errorf("ResolveFilePath(\"/tmp/flag-log.txt\") = %q, want %q", got, want)
	}
}
