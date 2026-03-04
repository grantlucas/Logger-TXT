package entry

import (
	"testing"
	"time"
)

func TestFormatMessageOnly(t *testing.T) {
	tm := time.Date(2026, 2, 22, 9, 15, 0, 0, time.FixedZone("EST", -5*3600))
	e := Entry{
		Time:    tm,
		Message: "Grabbed a coffee",
	}
	got := e.Format()
	want := "22/02/26 09:15 -0500 - Grabbed a coffee"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestFormatTypeAndMessage(t *testing.T) {
	tm := time.Date(2026, 2, 22, 10, 31, 0, 0, time.FixedZone("EST", -5*3600))
	e := Entry{
		Time:    tm,
		Type:    "WORK",
		Message: "Fixed login bug",
	}
	got := e.Format()
	want := "22/02/26 10:31 -0500 - WORK - Fixed login bug"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestFormatProjectAndMessage(t *testing.T) {
	tm := time.Date(2026, 2, 22, 10, 32, 0, 0, time.FixedZone("EST", -5*3600))
	e := Entry{
		Time:    tm,
		Project: "API",
		Message: "Deployed v1.3.2",
	}
	got := e.Format()
	want := "22/02/26 10:32 -0500 - (API) - Deployed v1.3.2"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestFormatTypeProjectAndMessage(t *testing.T) {
	tm := time.Date(2026, 2, 22, 10, 33, 0, 0, time.FixedZone("EST", -5*3600))
	e := Entry{
		Time:    tm,
		Type:    "WORK",
		Project: "API",
		Message: "Reviewed pull request",
	}
	got := e.Format()
	want := "22/02/26 10:33 -0500 - WORK (API) - Reviewed pull request"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}

func TestFormatPositiveTimezone(t *testing.T) {
	tm := time.Date(2026, 3, 1, 14, 30, 0, 0, time.FixedZone("IST", 5*3600+1800))
	e := Entry{
		Time:    tm,
		Message: "Afternoon tea",
	}
	got := e.Format()
	want := "01/03/26 14:30 +0530 - Afternoon tea"
	if got != want {
		t.Errorf("Format() = %q, want %q", got, want)
	}
}
