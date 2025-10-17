package format

import (
	"strings"
	"testing"
)

func TestEmojiFormat_NoTrailingPadding(t *testing.T) {
	input := "I :heart: Bubble Gum :candy:"
	output, err := emoji(input)
	if err != nil {
		t.Fatalf("emoji() returned error: %v", err)
	}

	trimmed := strings.TrimRight(output, " \t\r\n")

	if !strings.HasSuffix(trimmed, "🍬") {
		t.Errorf("expected output to end with candy emoji; got trimmed=%q", trimmed)
	}

	if !strings.Contains(trimmed, "I") || !strings.Contains(trimmed, "Bubble Gum") {
		t.Errorf("output missing expected content; got trimmed=%q", trimmed)
	}

	if !strings.Contains(trimmed, "❤") {
		t.Errorf("output missing heart emoji expansion; got trimmed=%q", trimmed)
	}

	lines := strings.Split(output, "\n")
	if len(lines) > 0 {
		lastLine := lines[0]
		trailingSpaces := len(lastLine) - len(strings.TrimRight(lastLine, " "))

		if trailingSpaces > 5 {
			t.Errorf("excessive trailing spaces detected (%d spaces); output=%q", trailingSpaces, output)
		}
	}
}

func TestEmojiFormat_MultipleEmojis(t *testing.T) {
	input := ":smile: :rocket: :heart:"
	output, err := emoji(input)
	if err != nil {
		t.Fatalf("emoji() returned error: %v", err)
	}

	trimmed := strings.TrimSpace(output)

	expectedEmojis := []string{"😄", "🚀", "❤"}
	for _, emoji := range expectedEmojis {
		if !strings.Contains(trimmed, emoji) {
			t.Errorf("expected emoji %s in output; got trimmed=%q", emoji, trimmed)
		}
	}
}

func TestEmojiFormat_NoEmojis(t *testing.T) {
	input := "Plain text without emoji codes"
	output, err := emoji(input)
	if err != nil {
		t.Fatalf("emoji() returned error: %v", err)
	}

	trimmed := strings.TrimSpace(output)

	if !strings.Contains(trimmed, "Plain text") {
		t.Errorf("expected plain text to be preserved; got trimmed=%q", trimmed)
	}
}
