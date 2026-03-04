package cmd

import (
	"strings"
	"testing"
)

func assertHelpContains(t *testing.T, args []string, phrases []string) {
	t.Helper()
	out, _, err := executeCmd(t, args...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, phrase := range phrases {
		if !strings.Contains(out, phrase) {
			t.Errorf("help output missing phrase %q\nfull output:\n%s", phrase, out)
		}
	}
}

func TestRootHelp(t *testing.T) {
	assertHelpContains(t, []string{"--help"}, []string{
		"plain text file",
		"Log entry format:",
		"File resolution order:",
		"LOGGERTXT_PATH",
		"logger-txt add",
	})
}

func TestAddHelp(t *testing.T) {
	assertHelpContains(t, []string{"add", "--help"}, []string{
		"automatically uppercase",
		"logger-txt add",
	})
}
