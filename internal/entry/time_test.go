package entry

import (
	"testing"
	"time"
)

func TestParseInputTime_DateOnly(t *testing.T) {
	loc := time.FixedZone("EST", -5*3600)
	got, dateOnly, err := ParseInputTime("22/02/26", loc)
	if err != nil {
		t.Fatalf("ParseInputTime() error = %v", err)
	}
	if !dateOnly {
		t.Error("dateOnly = false, want true")
	}
	want := time.Date(2026, 2, 22, 0, 0, 0, 0, loc)
	if !got.Equal(want) {
		t.Errorf("time = %v, want %v", got, want)
	}
}

func TestParseInputTime_InvalidInput(t *testing.T) {
	loc := time.FixedZone("EST", -5*3600)
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"garbage", "not-a-date"},
		{"wrong format", "2026-02-22"},
		{"partial date", "22/02"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParseInputTime(tt.input, loc)
			if err == nil {
				t.Errorf("ParseInputTime(%q) expected error, got nil", tt.input)
			}
		})
	}
}

func TestParseInputTime_UsesLocation(t *testing.T) {
	est := time.FixedZone("EST", -5*3600)
	ist := time.FixedZone("IST", 5*3600+1800)

	gotEST, _, _ := ParseInputTime("22/02/26", est)
	gotIST, _, _ := ParseInputTime("22/02/26", ist)

	if gotEST.Location().String() != "EST" {
		t.Errorf("EST location = %v, want EST", gotEST.Location())
	}
	if gotIST.Location().String() != "IST" {
		t.Errorf("IST location = %v, want IST", gotIST.Location())
	}
	// Same wall clock, different instants
	if gotEST.Equal(gotIST) {
		t.Error("same date in different zones should not be equal instants")
	}
}

func TestEndOfDay(t *testing.T) {
	loc := time.FixedZone("EST", -5*3600)
	input := time.Date(2026, 2, 22, 0, 0, 0, 0, loc)
	got := EndOfDay(input)
	want := time.Date(2026, 2, 22, 23, 59, 0, 0, loc)
	if !got.Equal(want) {
		t.Errorf("EndOfDay() = %v, want %v", got, want)
	}
	if got.Location().String() != loc.String() {
		t.Errorf("location = %v, want %v", got.Location(), loc)
	}
}

func TestParseInputTime_DateTime(t *testing.T) {
	loc := time.FixedZone("EST", -5*3600)
	got, dateOnly, err := ParseInputTime("22/02/26 14:30", loc)
	if err != nil {
		t.Fatalf("ParseInputTime() error = %v", err)
	}
	if dateOnly {
		t.Error("dateOnly = true, want false")
	}
	want := time.Date(2026, 2, 22, 14, 30, 0, 0, loc)
	if !got.Equal(want) {
		t.Errorf("time = %v, want %v", got, want)
	}
}

func TestParseDateRange_Valid(t *testing.T) {
	loc := time.FixedZone("EST", -5*3600)

	start, end, err := ParseDateRange("22/02/26", "23/02/26", loc)
	if err != nil {
		t.Fatalf("ParseDateRange() error = %v", err)
	}

	wantStart := time.Date(2026, 2, 22, 0, 0, 0, 0, loc)
	if !start.Equal(wantStart) {
		t.Errorf("start = %v, want %v", start, wantStart)
	}

	// end is date-only, so EndOfDay should be applied
	wantEnd := time.Date(2026, 2, 23, 23, 59, 0, 0, loc)
	if !end.Equal(wantEnd) {
		t.Errorf("end = %v, want %v", end, wantEnd)
	}
}

func TestParseDateRange_WithTime(t *testing.T) {
	loc := time.FixedZone("EST", -5*3600)

	start, end, err := ParseDateRange("22/02/26 09:00", "22/02/26 17:00", loc)
	if err != nil {
		t.Fatalf("ParseDateRange() error = %v", err)
	}

	wantStart := time.Date(2026, 2, 22, 9, 0, 0, 0, loc)
	if !start.Equal(wantStart) {
		t.Errorf("start = %v, want %v", start, wantStart)
	}

	// end has time, so EndOfDay should NOT be applied
	wantEnd := time.Date(2026, 2, 22, 17, 0, 0, 0, loc)
	if !end.Equal(wantEnd) {
		t.Errorf("end = %v, want %v", end, wantEnd)
	}
}

func TestParseDateRange_InvalidStart(t *testing.T) {
	loc := time.FixedZone("EST", -5*3600)

	_, _, err := ParseDateRange("bad", "22/02/26", loc)
	if err == nil {
		t.Fatal("expected error for invalid start")
	}
}

func TestParseDateRange_InvalidEnd(t *testing.T) {
	loc := time.FixedZone("EST", -5*3600)

	_, _, err := ParseDateRange("22/02/26", "bad", loc)
	if err == nil {
		t.Fatal("expected error for invalid end")
	}
}
