package cmd

import (
	"strings"
	"testing"
)

func TestRootHelp(t *testing.T) {
	out, _, err := executeCmd(t, "--help")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	phrases := []string{
		"plain text file",
		"Log entry format:",
		"File resolution order:",
		"LOGGERTXT_PATH",
		"logger-txt add",
	}
	for _, phrase := range phrases {
		if !strings.Contains(out, phrase) {
			t.Errorf("root help missing phrase %q", phrase)
		}
	}
}
