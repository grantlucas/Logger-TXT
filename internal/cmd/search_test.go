package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchCmd_CaseInsensitiveDefault(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - WORK (API) - Started standup\n"+
		"03/03/26 09:30 -0500 - Grabbed a coffee\n"+
		"03/03/26 10:00 -0500 - WORK (API) - Fixed auth bug\n")

	out, _, err := executeCmd(t, "--file", logFile, "search", "work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "03/03/26 09:00 -0500 - WORK (API) - Started standup\n" +
		"03/03/26 10:00 -0500 - WORK (API) - Fixed auth bug\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestSearchCmd_CaseSensitive(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - WORK - Task one\n"+
		"03/03/26 09:30 -0500 - work - Task two\n"+
		"03/03/26 10:00 -0500 - WORK - Task three\n")

	out, _, err := executeCmd(t, "--file", logFile, "search", "--case-sensitive", "work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "03/03/26 09:30 -0500 - work - Task two\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestSearchCmd_CountLimit(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - (API) - Task one\n"+
		"03/03/26 09:30 -0500 - (API) - Task two\n"+
		"03/03/26 10:00 -0500 - (API) - Task three\n")

	out, _, err := executeCmd(t, "--file", logFile, "search", "-c", "2", "api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "03/03/26 09:30 -0500 - (API) - Task two\n" +
		"03/03/26 10:00 -0500 - (API) - Task three\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestSearchCmd_NoMatches(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - WORK - Task one\n")

	out, _, err := executeCmd(t, "--file", logFile, "search", "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "" {
		t.Errorf("expected no output for no matches, got: %q", out)
	}
}

func TestSearchCmd_PathWithSpaces(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "path with spaces")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - WORK - Spaced path entry\n"+
		"03/03/26 09:30 -0500 - Coffee break\n")

	out, _, err := executeCmd(t, "--file", logFile, "search", "work")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "03/03/26 09:00 -0500 - WORK - Spaced path entry\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestSearchCmd_DateRange(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir,
		"20/02/26 09:00 -0500 - WORK - Too early\n"+
			"22/02/26 10:00 -0500 - WORK - In range match\n"+
			"22/02/26 11:00 -0500 - Coffee break\n"+
			"22/02/26 12:00 -0500 - WORK - In range match two\n"+
			"25/02/26 09:00 -0500 - WORK - Too late\n")

	out, _, err := executeCmd(t, "--file", logFile, "search", "work", "--start", "22/02/26", "--end", "22/02/26")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "22/02/26 10:00 -0500 - WORK - In range match\n" +
		"22/02/26 12:00 -0500 - WORK - In range match two\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestSearchCmd_DateRangeWithCount(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir,
		"22/02/26 09:00 -0500 - WORK - First\n"+
			"22/02/26 10:00 -0500 - WORK - Second\n"+
			"22/02/26 11:00 -0500 - WORK - Third\n")

	out, _, err := executeCmd(t, "--file", logFile, "search", "work", "--start", "22/02/26", "--end", "22/02/26", "-c", "2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "22/02/26 10:00 -0500 - WORK - Second\n" +
		"22/02/26 11:00 -0500 - WORK - Third\n"
	if out != expected {
		t.Errorf("output mismatch\ngot:  %q\nwant: %q", out, expected)
	}
}

func TestSearchCmd_StartWithoutEnd(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "22/02/26 10:00 -0500 - Entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "search", "entry", "--start", "22/02/26")
	if err == nil {
		t.Fatal("expected error when --start provided without --end")
	}
}

func TestSearchCmd_EndWithoutStart(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "22/02/26 10:00 -0500 - Entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "search", "entry", "--end", "22/02/26")
	if err == nil {
		t.Fatal("expected error when --end provided without --start")
	}
}

func TestSearchCmd_InvalidStartDate(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "22/02/26 10:00 -0500 - Entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "search", "entry", "--start", "bad", "--end", "22/02/26")
	if err == nil {
		t.Fatal("expected error for invalid start date")
	}
}

func TestSearchCmd_InvalidEndDate(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "22/02/26 10:00 -0500 - Entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "search", "entry", "--start", "22/02/26", "--end", "bad")
	if err == nil {
		t.Fatal("expected error for invalid end date")
	}
}

func TestSearchCmd_DateRangeFileNotFound(t *testing.T) {
	_, _, err := executeCmd(t, "--file", "/nonexistent/log.txt", "search", "test", "--start", "22/02/26", "--end", "22/02/26")
	if err == nil {
		t.Fatal("expected error for non-existent file with date range")
	}
}

func TestSearchCmd_FileNotFound(t *testing.T) {
	_, _, err := executeCmd(t, "--file", "/nonexistent/path/log.txt", "search", "test")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}
