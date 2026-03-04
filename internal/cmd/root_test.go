package cmd

import (
	"os"
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

func TestExecute(t *testing.T) {
	// Execute() is the public entry point used by main.go
	// With no args and no log file, it should return an error (file not found)
	// but we just need to exercise the function.
	// Save and restore os.Args since Execute creates a fresh root command.
	oldArgs := os.Args
	os.Args = []string{"logger-txt", "--help"}
	defer func() { os.Args = oldArgs }()

	err := Execute()
	if err != nil {
		t.Fatalf("unexpected error from Execute: %v", err)
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
