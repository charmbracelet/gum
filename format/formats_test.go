package format

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// TestMarkdownTableWraps checks that a table with a cell longer than the wrap
// width wraps the cell content instead of letting the whole table overflow.
// Regression test for https://github.com/charmbracelet/gum/issues/681.
func TestMarkdownTableWraps(t *testing.T) {
	const longCell = "This is a very long description that definitely exceeds the width of one terminal line and should wrap within the cell instead of breaking the whole table layout across all cells"

	input := "| Name | Description |\n" +
		"| --- | --- |\n" +
		"| Long | " + longCell + " |\n"

	out, err := markdown(input, "pink")
	if err != nil {
		t.Fatalf("markdown() error: %v", err)
	}

	// With wrapping, no rendered line should be anywhere near as wide as the
	// unwrapped long cell. Without the fix (word wrap 0) the long cell renders
	// on a single line far exceeding this bound.
	limit := defaultMarkdownWidth + 20 // allow slack for borders/padding
	for _, line := range strings.Split(out, "\n") {
		if w := lipgloss.Width(line); w > limit {
			t.Fatalf("line exceeds wrap width: %d > %d\nline: %q", w, limit, line)
		}
	}
}
