package cmd

import (
	"bytes"
	"os"
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
	t.Cleanup(func() { os.Chmod(logFile, 0644) })

	_, _, err := executeCmd(t, "--file", logFile, "delete", "--yes")
	if err == nil {
		t.Fatal("expected error when file is read-only")
	}
}

func TestDeleteCmd_FileNotFound(t *testing.T) {
	_, _, err := executeCmd(t, "--file", "/nonexistent/path/log.txt", "delete", "--yes")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}
