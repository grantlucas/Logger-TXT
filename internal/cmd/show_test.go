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

func TestShowCmd_PathWithSpaces(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "path with spaces")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - WORK - Spaced path entry\n")

	out, _, err := executeCmd(t, "--file", logFile, "show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "03/03/26 09:00 -0500 - WORK - Spaced path entry\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestShowCmd_DateRange(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir,
		"20/02/26 09:00 -0500 - Too early\n"+
			"22/02/26 10:00 -0500 - In range first\n"+
			"22/02/26 11:00 -0500 - In range second\n"+
			"25/02/26 09:00 -0500 - Too late\n")

	out, _, err := executeCmd(t, "--file", logFile, "show", "--start", "22/02/26", "--end", "22/02/26")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "22/02/26 10:00 -0500 - In range first\n" +
		"22/02/26 11:00 -0500 - In range second\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestShowCmd_DateRangeWithCount(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir,
		"22/02/26 09:00 -0500 - First\n"+
			"22/02/26 10:00 -0500 - Second\n"+
			"22/02/26 11:00 -0500 - Third\n")

	out, _, err := executeCmd(t, "--file", logFile, "show", "--start", "22/02/26", "--end", "22/02/26", "-c", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "22/02/26 10:00 -0500 - Second\n" +
		"22/02/26 11:00 -0500 - Third\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestShowCmd_DateRangeWithTime(t *testing.T) {
	t.Setenv("TZ", "EST5EDT")
	dir := t.TempDir()
	logFile := writeLogFile(t, dir,
		"22/02/26 08:00 -0500 - Before range\n"+
			"22/02/26 10:00 -0500 - In range\n"+
			"22/02/26 14:00 -0500 - After range\n")

	out, _, err := executeCmd(t, "--file", logFile, "show", "--start", "22/02/26 09:00", "--end", "22/02/26 12:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "22/02/26 10:00 -0500 - In range\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestShowCmd_StartWithoutEnd(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "22/02/26 10:00 -0500 - Entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "show", "--start", "22/02/26")
	if err == nil {
		t.Fatal("expected error when --start provided without --end")
	}
}

func TestShowCmd_EndWithoutStart(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "22/02/26 10:00 -0500 - Entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "show", "--end", "22/02/26")
	if err == nil {
		t.Fatal("expected error when --end provided without --start")
	}
}

func TestShowCmd_InvalidStartDate(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "22/02/26 10:00 -0500 - Entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "show", "--start", "not-a-date", "--end", "22/02/26")
	if err == nil {
		t.Fatal("expected error for invalid start date")
	}
}

func TestShowCmd_InvalidEndDate(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "22/02/26 10:00 -0500 - Entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "show", "--start", "22/02/26", "--end", "not-a-date")
	if err == nil {
		t.Fatal("expected error for invalid end date")
	}
}

func TestShowCmd_DateRangeFileNotFound(t *testing.T) {
	_, _, err := executeCmd(t, "--file", "/nonexistent/log.txt", "show", "--start", "22/02/26", "--end", "22/02/26")
	if err == nil {
		t.Fatal("expected error for non-existent file with date range")
	}
}

func TestShowCmd_FileNotFound(t *testing.T) {
	_, _, err := executeCmd(t, "--file", "/nonexistent/path/log.txt", "show")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}
