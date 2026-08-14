package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStructuredLogPreservesTimeContextKey(t *testing.T) {
	out := filepath.Join(t.TempDir(), "gum.log")

	err := (Options{
		File:       out,
		Formatter:  "text",
		Level:      "none",
		Structured: true,
		Text: []string{
			"Test log message",
			"time", "0",
			"normal_var_name", "0",
			"time", "0",
		},
	}).Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("expected log file to be written: %v", err)
	}

	got := string(data)
	if count := strings.Count(got, "time=0"); count != 2 {
		t.Fatalf("expected two time context values in %q, got %d", got, count)
	}

	if !strings.Contains(got, "normal_var_name=0") {
		t.Fatalf("expected normal context value in %q", got)
	}
}
