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
