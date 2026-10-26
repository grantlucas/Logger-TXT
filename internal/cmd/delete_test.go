package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func executeCmdWithInput(t *testing.T, input string, args ...string) (string, string, error) {
	t.Helper()
	cmd := NewRootCmd()
	out := new(bytes.Buffer)
	errOut := new(bytes.Buffer)
	cmd.SetOut(out)
	cmd.SetErr(errOut)
	cmd.SetIn(strings.NewReader(input))
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), errOut.String(), err
}

func TestDeleteCmd_YesFlag(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - First entry\n"+
		"03/03/26 09:30 -0500 - Second entry\n")

	out, _, err := executeCmd(t, "--file", logFile, "delete", "--yes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "Deleted last line from file") {
		t.Errorf("expected deletion confirmation, got: %q", out)
	}

	data, _ := os.ReadFile(logFile)
	content := string(data)
	if !strings.Contains(content, "First entry") {
		t.Errorf("first entry should remain, got: %q", content)
	}
	if strings.Contains(content, "Second entry") {
		t.Errorf("second entry should be deleted, got: %q", content)
	}
}

func TestDeleteCmd_InteractiveConfirmY(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - First entry\n"+
		"03/03/26 09:30 -0500 - Second entry\n")

	out, _, err := executeCmdWithInput(t, "Y\n", "--file", logFile, "delete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "Second entry") {
		t.Errorf("should show preview of last entry, got: %q", out)
	}
	if !strings.Contains(out, "Do you wish to continue?") {
		t.Errorf("should show prompt, got: %q", out)
	}
	if !strings.Contains(out, "Deleted last line from file") {
		t.Errorf("should confirm deletion, got: %q", out)
	}

	data, _ := os.ReadFile(logFile)
	if strings.Contains(string(data), "Second entry") {
		t.Errorf("second entry should be deleted")
	}
}

func TestDeleteCmd_InteractiveDenyN(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - First entry\n"+
		"03/03/26 09:30 -0500 - Second entry\n")

	out, _, err := executeCmdWithInput(t, "n\n", "--file", logFile, "delete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "No line deleted") {
		t.Errorf("should show cancellation message, got: %q", out)
	}

	// File should be unchanged
	data, _ := os.ReadFile(logFile)
	if !strings.Contains(string(data), "Second entry") {
		t.Errorf("second entry should still exist")
	}
}

func TestDeleteCmd_InteractiveDenyRandom(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - Entry\n")

	out, _, err := executeCmdWithInput(t, "maybe\n", "--file", logFile, "delete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "No line deleted") {
		t.Errorf("any answer other than Y should cancel, got: %q", out)
	}
}

func TestDeleteCmd_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "")

	_, _, err := executeCmd(t, "--file", logFile, "delete", "--yes")
	if err == nil {
		t.Fatal("expected error for empty file")
	}
}

func TestDeleteCmd_DeleteLastError(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - Entry\n")

	// Make file read-only so DeleteLast (which writes) fails
	if err := os.Chmod(logFile, 0444); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(logFile, 0644) })

	_, _, err := executeCmd(t, "--file", logFile, "delete", "--yes")
	if err == nil {
		t.Fatal("expected error when file is read-only")
	}
}

func TestDeleteCmd_PathWithSpaces(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "path with spaces")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - First entry\n"+
		"03/03/26 09:30 -0500 - Second entry\n")

	out, _, err := executeCmd(t, "--file", logFile, "delete", "--yes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "Deleted last line from file") {
		t.Errorf("expected deletion confirmation, got: %q", out)
	}

	data, _ := os.ReadFile(logFile)
	content := string(data)
	if !strings.Contains(content, "First entry") {
		t.Errorf("first entry should remain, got: %q", content)
	}
	if strings.Contains(content, "Second entry") {
		t.Errorf("second entry should be deleted, got: %q", content)
	}
}

func TestDeleteCmd_FileNotFound(t *testing.T) {
	_, _, err := executeCmd(t, "--file", "/nonexistent/path/log.txt", "delete", "--yes")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

// writeOnRead appends line to the log when the prompt reads its answer,
// standing in for another writer that runs while the user is deciding.
type writeOnRead struct {
	t      *testing.T
	path   string
	line   string
	answer io.Reader
}

func (r *writeOnRead) Read(p []byte) (int, error) {
	if r.line != "" {
		f, err := os.OpenFile(r.path, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			r.t.Fatal(err)
		}
		_, _ = f.WriteString(r.line + "\n")
		_ = f.Close()
		r.line = ""
	}
	return r.answer.Read(p)
}

func TestDeleteCmd_RefusesWhenLastLineChangesDuringPrompt(t *testing.T) {
	const content = "03/03/26 09:00 -0500 - First entry\n03/03/26 09:30 -0500 - Typo entry\n"
	logFile := writeLogFile(t, t.TempDir(), content)
	const added = "03/03/26 10:00 -0500 - Added while prompting"

	cmd := NewRootCmd()
	out := new(bytes.Buffer)
	cmd.SetOut(out)
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetIn(&writeOnRead{t: t, path: logFile, line: added, answer: strings.NewReader("Y\n")})
	cmd.SetArgs([]string{"--file", logFile, "delete"})

	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "changed") {
		t.Fatalf("error = %v, want one saying the last line changed", err)
	}
	if data, _ := os.ReadFile(logFile); string(data) != content+added+"\n" {
		t.Errorf("file content = %q, want nothing deleted", data)
	}
}

func TestDeleteCmd_PrintsDeletedLine(t *testing.T) {
	logFile := writeLogFile(t, t.TempDir(), "03/03/26 09:00 -0500 - First entry\n03/03/26 09:30 -0500 - Second entry\n")

	out, _, err := executeCmd(t, "--file", logFile, "delete", "--yes")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "Deleted last line from file:\n03/03/26 09:30 -0500 - Second entry") {
		t.Errorf("output = %q, want the deleted line", out)
	}
}
