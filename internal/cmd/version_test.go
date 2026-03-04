package cmd

import (
	"strings"
	"testing"
)

func TestVersionCmd_DevVersion(t *testing.T) {
	out, _, err := executeCmd(t, "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Default (dev) version should show basic info
	if !strings.Contains(out, "Logger-TXT") {
		t.Errorf("expected Logger-TXT in output, got: %q", out)
	}
	if !strings.Contains(out, "Version dev") {
		t.Errorf("expected dev version, got: %q", out)
	}
	if !strings.Contains(out, "Author: Grant Lucas") {
		t.Errorf("expected author info, got: %q", out)
	}
	if !strings.Contains(out, "License: GPL") {
		t.Errorf("expected license info, got: %q", out)
	}
	// Dev version should NOT show release URL
	if strings.Contains(out, "Release:") {
		t.Errorf("dev version should not show release URL, got: %q", out)
	}
}

func TestVersionCmd_ReleaseVersion(t *testing.T) {
	// Temporarily set version to a release value
	oldVersion := version
	version = "2.0.0"
	defer func() { version = oldVersion }()

	out, _, err := executeCmd(t, "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(out, "Version 2.0.0") {
		t.Errorf("expected release version, got: %q", out)
	}
	if !strings.Contains(out, "Release: https://github.com/grantlucas/Logger-TXT/releases/tag/v2.0.0") {
		t.Errorf("expected release URL, got: %q", out)
	}
}
