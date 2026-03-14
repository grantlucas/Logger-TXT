package cmd

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Pin time.Local to a fixed zone matching the -0500 offsets used in test
	// log entries. This prevents CI environments (typically UTC) from
	// interpreting date-only range boundaries differently than local dev
	// machines, which caused entries near end-of-day to fall outside the
	// range when parsed in UTC.
	time.Local = time.FixedZone("EST", -5*3600)

	os.Exit(m.Run())
}
