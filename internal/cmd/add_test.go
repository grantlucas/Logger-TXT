package cmd

import (
	"os"
	"path/filepath"
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

func TestAddCmd_AppendError(t *testing.T) {
	// Point to a path inside a non-existent read-only location
	_, _, err := executeCmd(t, "--file", "/dev/null/impossible/log.txt", "add", "test")
	if err == nil {
		t.Fatal("expected error when append fails")
	}
}

func TestAddCmd_PathWithSpaces(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "path with spaces")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	logFile := writeLogFile(t, dir, "")

	out, _, err := executeCmd(t, "--file", logFile, "add", "-t", "work", "Spaced path entry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, `"Spaced path entry" logged under the type WORK`) {
		t.Errorf("expected confirmation, got: %q", out)
	}

	data, _ := os.ReadFile(logFile)
	if !strings.Contains(string(data), " - WORK - Spaced path entry") {
		t.Errorf("expected entry in file, got: %q", string(data))
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

func TestAddCmd_At(t *testing.T) {
	const existing = "03/03/26 09:00 -0500 - A\n05/03/26 09:00 -0500 - C\n"
	tests := []struct {
		name     string
		args     []string
		wantOut  string
		wantErr  string
		wantFile string
	}{
		{
			name:     "inserts between existing entries",
			args:     []string{"add", "--at", "04/03/26 14:30", "-t", "work", "Reviewed", "budget"},
			wantOut:  `"Reviewed budget" logged under the type WORK at 04/03/26 14:30 -0500`,
			wantFile: "03/03/26 09:00 -0500 - A\n04/03/26 14:30 -0500 - WORK - Reviewed budget\n05/03/26 09:00 -0500 - C\n",
		},
		{
			name:     "appends when newest",
			args:     []string{"add", "--at", "06/03/26 08:15", "D"},
			wantOut:  `"D" logged at 06/03/26 08:15 -0500`,
			wantFile: existing + "06/03/26 08:15 -0500 - D\n",
		},
		{
			name:     "rejects date without time",
			args:     []string{"add", "--at", "04/03/26", "X"},
			wantErr:  "--at requires a time",
			wantFile: existing,
		},
		{
			name:     "rejects future time",
			args:     []string{"add", "--at", "01/01/68 10:00", "X"},
			wantErr:  "is in the future",
			wantFile: existing,
		},
		{
			// e.g. --at "$WHEN" with WHEN unset must not log at the current time
			name:     "rejects empty value",
			args:     []string{"add", "--at", "", "X"},
			wantErr:  "invalid --at value",
			wantFile: existing,
		},
		{
			name:     "rejects unparseable time",
			args:     []string{"add", "--at", "yesterday 3pm", "X"},
			wantErr:  "invalid --at value",
			wantFile: existing,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logFile := writeLogFile(t, t.TempDir(), existing)

			out, _, err := executeCmd(t, append([]string{"--file", logFile}, tt.args...)...)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("error = %v, want containing %q", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !strings.Contains(out, tt.wantOut) {
				t.Errorf("output = %q, want containing %q", out, tt.wantOut)
			}
			data, _ := os.ReadFile(logFile)
			if string(data) != tt.wantFile {
				t.Errorf("file content = %q, want %q", data, tt.wantFile)
			}
		})
	}
}
