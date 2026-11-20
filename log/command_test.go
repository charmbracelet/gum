package log

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogPreservesNewlines(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	msg := "this\nis\na\ntest"
	opts := Options{
		File: logFile,
		Text: []string{msg},
	}

	if err := opts.Run(); err != nil {
		t.Fatalf("unexpected error running log command: %v", err)
	}

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed reading log output: %v", err)
	}

	output := string(data)
	if !strings.Contains(output, "this\nis\na\ntest") {
		t.Errorf("expected output to contain preserved newlines %q, got %q", msg, output)
	}
}

func TestLogMultilineScenarios(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		opts     Options
		expected []string
	}{
		{
			name: "single line message",
			opts: Options{
				Text: []string{"simple message"},
			},
			expected: []string{"simple message"},
		},
		{
			name: "multiline with trailing newline",
			opts: Options{
				Text: []string{"line1\nline2\nline3\n"},
			},
			expected: []string{"line1\nline2\nline3"},
		},
		{
			name: "formatted multiline message",
			opts: Options{
				Format: true,
				Text:   []string{"first line\n%s\nthird line", "second line"},
			},
			expected: []string{"first line\nsecond line\nthird line"},
		},
		{
			name: "structured logging with multiline message",
			opts: Options{
				Structured: true,
				Text:       []string{"event occurred:\ndetail 1\ndetail 2", "code", "123"},
			},
			expected: []string{"event occurred:\ndetail 1\ndetail 2"},
		},
		{
			name: "multiline message with info level and prefix",
			opts: Options{
				Level:  "info",
				Prefix: "TEST",
				Text:   []string{"line A\nline B"},
			},
			expected: []string{"line A\nline B"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			logFile := filepath.Join(tmpDir, "multiline.log")
			opts := tc.opts
			opts.File = logFile

			if err := opts.Run(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			data, err := os.ReadFile(logFile)
			if err != nil {
				t.Fatalf("failed reading log file: %v", err)
			}

			out := string(data)
			for _, exp := range tc.expected {
				if !strings.Contains(out, exp) {
					t.Errorf("expected log output to contain %q, but got: %q", exp, out)
				}
			}
		})
	}
}

func TestLogEmptyText(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "empty.log")

	opts := Options{
		File: logFile,
		Text: []string{},
	}

	if err := opts.Run(); err != nil {
		t.Fatalf("expected nil error on empty text, got: %v", err)
	}
}
