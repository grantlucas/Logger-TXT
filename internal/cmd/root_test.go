package cmd

import (
	"testing"
)

func TestRootCmd_DefaultShowBehavior(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - WORK - Entry one\n"+
		"03/03/26 09:30 -0500 - WORK - Entry two\n")

	// Running root with no subcommand should behave like "show"
	out, _, err := executeCmd(t, "--file", logFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "03/03/26 09:00 -0500 - WORK - Entry one\n" +
		"03/03/26 09:30 -0500 - WORK - Entry two\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestRootCmd_FileFlagPropagates(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - Entry\n")

	// --file flag before subcommand should work
	out, _, err := executeCmd(t, "--file", logFile, "show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "03/03/26 09:00 -0500 - Entry\n" {
		t.Errorf("expected entry output, got: %q", out)
	}
}

func TestRootCmd_HelpFlag(t *testing.T) {
	out, _, err := executeCmd(t, "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out == "" {
		t.Error("expected help output")
	}
}
