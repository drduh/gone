package config

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
	"time"
)

var testTime = time.Date(
	2026, 12, 31, 12, 0, 0, 0, time.UTC)

// stubNow sets fixed time for tests.
func stubNow(t *testing.T, now time.Time) {
	t.Helper()
	orig := funcNow
	funcNow = func() time.Time { return now }
	t.Cleanup(func() { funcNow = orig })
}

// stubExit records exit codes for tests.
func stubExit(t *testing.T) *int {
	t.Helper()
	orig := funcExit
	code := -1
	funcExit = func(c int) { code = c }
	t.Cleanup(func() { funcExit = orig })
	return &code
}

// TestSetStart tests start time assignment.
func TestSetStart(t *testing.T) {
	stubNow(t, testTime)
	a := &App{}
	a.Start()
	if !a.StartTime.Equal(testTime) {
		t.Fatalf("expected StartTime %v, got %v",
			testTime, a.StartTime)
	}
}

// TestUptime tests various uptime logging.
func TestUptime(t *testing.T) {
	cases := []struct {
		name      string
		startTime time.Time
		now       time.Time
		want      string
	}{
		{
			name:      "zero start time",
			startTime: time.Time{},
			now:       testTime,
			want:      "0s",
		},
		{
			name:      "round duration",
			startTime: testTime,
			now: testTime.Add(
				3*time.Second + 200*time.Millisecond),
			want: "3s",
		},
		{
			name:      "future start time",
			startTime: testTime.Add(180 * time.Second),
			now:       testTime,
			want:      "0s",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			stubNow(t, tt.now)

			a := &App{StartTime: tt.startTime}
			if got := a.Uptime(); got != tt.want {
				t.Fatalf("expected uptime %q, got %q",
					tt.want, got)
			}
		})
	}
}

// TestStop tests stop with log.
func TestStop(t *testing.T) {
	stubNow(t, testTime.Add(2*time.Second))
	exitCode := stubExit(t)

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	a := &App{StartTime: testTime, Log: logger}
	a.Stop("test")
	if *exitCode != 0 {
		t.Fatalf("expected exit 0, got %d",
			*exitCode)
	}

	out := buf.String()
	for _, want := range []string{
		"stopping application",
		"reason=test",
		"uptime=2s",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("output should have %q, got %q",
				want, out)
		}
	}
}

// TestStopNilLog tests stops without log.
func TestStopNilLog(t *testing.T) {
	stubNow(t, testTime.Add(1*time.Second))
	exitCode := stubExit(t)

	a := &App{
		StartTime: testTime,
		Log:       nil,
	}

	a.Stop("test")
	if *exitCode != 0 {
		t.Fatalf("expected exit 0, got %d",
			*exitCode)
	}
}
