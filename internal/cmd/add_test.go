package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestAddCmd_WithType(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "")

	out, _, err := executeCmd(t, "--file", logFile, "add", "-t", "work", "Fixed login bug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `"Fixed login bug" logged under the type WORK`) {
		t.Errorf("expected type confirmation, got: %q", out)
	}

	data, _ := os.ReadFile(logFile)
	content := string(data)
	if !strings.Contains(content, " - WORK - Fixed login bug") {
		t.Errorf("expected WORK type in file, got: %q", content)
	}
}

func TestAddCmd_WithProject(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "")

	out, _, err := executeCmd(t, "--file", logFile, "add", "-p", "api", "Deployed v1.3.2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `"Deployed v1.3.2" logged in the project API`) {
		t.Errorf("expected project confirmation, got: %q", out)
	}

	data, _ := os.ReadFile(logFile)
	content := string(data)
	if !strings.Contains(content, " - (API) - Deployed v1.3.2") {
		t.Errorf("expected (API) project in file, got: %q", content)
	}
}

func TestAddCmd_WithTypeAndProject(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "")

	out, _, err := executeCmd(t, "--file", logFile, "add", "-t", "work", "-p", "api", "Reviewed pull request")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `"Reviewed pull request" logged under the type WORK in the project API`) {
		t.Errorf("expected type+project confirmation, got: %q", out)
	}

	data, _ := os.ReadFile(logFile)
	content := string(data)
	if !strings.Contains(content, " - WORK (API) - Reviewed pull request") {
		t.Errorf("expected WORK (API) in file, got: %q", content)
	}
}

func TestAddCmd_MultiWordArgs(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "")

	out, _, err := executeCmd(t, "--file", logFile, "add", "-t", "personal", "Picked", "up", "groceries")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `"Picked up groceries" logged under the type PERSONAL`) {
		t.Errorf("expected joined message, got: %q", out)
	}

	data, _ := os.ReadFile(logFile)
	content := string(data)
	if !strings.Contains(content, "Picked up groceries") {
		t.Errorf("expected joined message in file, got: %q", content)
	}
}

func TestAddCmd_AppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "03/03/26 09:00 -0500 - Existing entry\n")

	_, _, err := executeCmd(t, "--file", logFile, "add", "New entry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(logFile)
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), string(data))
	}
	if !strings.Contains(lines[0], "Existing entry") {
		t.Errorf("first line should be original, got: %q", lines[0])
	}
	if !strings.Contains(lines[1], "New entry") {
		t.Errorf("second line should be new, got: %q", lines[1])
	}
}

func TestAddCmd_SimpleMessage(t *testing.T) {
	dir := t.TempDir()
	logFile := writeLogFile(t, dir, "")

	out, _, err := executeCmd(t, "--file", logFile, "add", "Grabbed a coffee")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check confirmation output
	if !strings.Contains(out, `"Grabbed a coffee" logged`) {
		t.Errorf("expected confirmation message, got: %q", out)
	}

	// Check file was written
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.HasSuffix(strings.TrimRight(content, "\n"), " - Grabbed a coffee") {
		t.Errorf("unexpected file content: %q", content)
	}
}
