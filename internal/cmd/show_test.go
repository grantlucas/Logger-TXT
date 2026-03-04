package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeLogFile(t *testing.T, dir string, content string) string {
	t.Helper()
	logFile := filepath.Join(dir, "log.txt")
	if err := os.WriteFile(logFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return logFile
}

func executeCmd(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := NewRootCmd()
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), errOut.String(), err
}

func TestShowCmd_PrintsLines(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - WORK - First entry\n"+
		"03/03/26 09:30 -0500 - WORK - Second entry\n"+
		"03/03/26 10:00 -0500 - WORK - Third entry\n")

	out, _, err := executeCmd(t, "--file", logFile, "show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "03/03/26 09:00 -0500 - WORK - First entry\n" +
		"03/03/26 09:30 -0500 - WORK - Second entry\n" +
		"03/03/26 10:00 -0500 - WORK - Third entry\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestShowCmd_CustomCount(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - Line 1\n"+
		"03/03/26 09:30 -0500 - Line 2\n"+
		"03/03/26 10:00 -0500 - Line 3\n")

	out, _, err := executeCmd(t, "--file", logFile, "show", "-c", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "03/03/26 09:30 -0500 - Line 2\n" +
		"03/03/26 10:00 -0500 - Line 3\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestShowCmd_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "")

	out, _, err := executeCmd(t, "--file", logFile, "show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "" {
		t.Errorf("expected no output for empty file, got: %q", out)
	}
}

func TestShowCmd_FileNotFound(t *testing.T) {
	_, _, err := executeCmd(t, "--file", "/nonexistent/path/log.txt", "show")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}
