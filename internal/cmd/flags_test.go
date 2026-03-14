package cmd

import (
	"testing"
	"time"

	"github.com/grantlucas/Logger-TXT/internal/entry"
)

func testEntry(typ, project string) entry.Entry {
	return entry.Entry{
		Time:    time.Date(2026, 3, 14, 10, 0, 0, 0, time.UTC),
		Type:    typ,
		Project: project,
		Message: "test message",
	}
}

func TestBuildEntryFilter_NilWhenBothEmpty(t *testing.T) {
	filter := buildEntryFilter("", "")
	if filter != nil {
		t.Fatal("buildEntryFilter(\"\", \"\") should return nil")
	}
}

func TestBuildEntryFilter_TypeOnly(t *testing.T) {
	filter := buildEntryFilter("WORK", "")
	if filter == nil {
		t.Fatal("expected non-nil filter")
	}
	if !filter(testEntry("WORK", "API")) {
		t.Error("should match WORK entry")
	}
	if filter(testEntry("MEETING", "API")) {
		t.Error("should not match MEETING entry")
	}
}

func TestBuildEntryFilter_ProjectOnly(t *testing.T) {
	filter := buildEntryFilter("", "API")
	if filter == nil {
		t.Fatal("expected non-nil filter")
	}
	if !filter(testEntry("WORK", "API")) {
		t.Error("should match API project")
	}
	if filter(testEntry("WORK", "WEB")) {
		t.Error("should not match WEB project")
	}
}

func TestBuildEntryFilter_BothANDLogic(t *testing.T) {
	filter := buildEntryFilter("WORK", "API")
	if filter == nil {
		t.Fatal("expected non-nil filter")
	}
	if !filter(testEntry("WORK", "API")) {
		t.Error("should match WORK+API")
	}
	if filter(testEntry("WORK", "WEB")) {
		t.Error("should not match WORK+WEB")
	}
	if filter(testEntry("MEETING", "API")) {
		t.Error("should not match MEETING+API")
	}
}

func TestBuildEntryFilter_CaseInsensitive(t *testing.T) {
	filter := buildEntryFilter("WORK", "API")
	if !filter(testEntry("work", "api")) {
		t.Error("should match case-insensitively")
	}
	if !filter(testEntry("Work", "Api")) {
		t.Error("should match mixed case")
	}
}
