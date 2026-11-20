package pager

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

func testModel(content string, width, height int) model {
	vp := viewport.New(viewport.WithWidth(width), viewport.WithHeight(height))
	vp.Style = vp.Style.BorderStyle(lipgloss.NormalBorder()).Margin(1)

	return model{
		content:         content,
		origContent:     content,
		viewport:        vp,
		help:            help.New(),
		showLineNumbers: true,
		softWrap:        true,
	}
}

func TestLayoutContentFillsInteriorHeight(t *testing.T) {
	m := testModel("line\nline\nline", 30, 10)
	interior := m.viewport.Height() - m.viewport.Style.GetVerticalFrameSize()

	got := m.layoutContent(m.viewport.Width(), m.viewport.Height())

	if h := lipgloss.Height(got); h != interior {
		t.Errorf("expected content to fill %d rows inside the frame, got %d", interior, h)
	}
}

func TestLayoutContentNoFillerForExactFit(t *testing.T) {
	m := testModel("", 30, 10)
	interior := m.viewport.Height() - m.viewport.Style.GetVerticalFrameSize()

	lines := make([]string, 0, interior)
	for i := 1; i <= interior; i++ {
		lines = append(lines, fmt.Sprintf("line %d", i))
	}
	m.content = strings.Join(lines, "\n")

	got := m.layoutContent(m.viewport.Width(), m.viewport.Height())

	if h := lipgloss.Height(got); h != interior {
		t.Errorf("expected exact fit of %d rows, got %d", interior, h)
	}
	if strings.Contains(ansi.Strip(got), "~") {
		t.Error("exact fit should not add filler lines")
	}
}

func TestLayoutContentTrailingNewlineInput(t *testing.T) {
	m := testModel("first\nsecond\n", 30, 10)
	interior := m.viewport.Height() - m.viewport.Style.GetVerticalFrameSize()

	got := m.layoutContent(m.viewport.Width(), m.viewport.Height())

	if h := lipgloss.Height(got); h != interior {
		t.Errorf("expected %d rows for trailing-newline input, got %d", interior, h)
	}
	stripped := ansi.Strip(got)
	if strings.Contains(stripped, "\n\n") {
		t.Error("a trailing newline in the input should not produce a phantom empty line")
	}
}

func TestLayoutContentLongContentNotPadded(t *testing.T) {
	var sb strings.Builder
	for i := 1; i <= 50; i++ {
		sb.WriteString(fmt.Sprintf("line %d", i))
		if i < 50 {
			sb.WriteString("\n")
		}
	}
	m := testModel(sb.String(), 30, 10)

	got := m.layoutContent(m.viewport.Width(), m.viewport.Height())

	if h := lipgloss.Height(got); h != 50 {
		t.Errorf("expected long content to keep its own height of 50 rows, got %d", h)
	}
	if strings.Contains(ansi.Strip(got), "~") {
		t.Error("long content should not be padded with filler lines")
	}
}

func TestProcessTextFittedContentNotScrollable(t *testing.T) {
	for _, n := range []int{1, 2, 5} {
		t.Run(fmt.Sprintf("content_lines=%d", n), func(t *testing.T) {
			lines := make([]string, n)
			for i := range lines {
				lines[i] = fmt.Sprintf("L%d", i+1)
			}
			m := testModel(strings.Join(lines, "\n"), 40, 12)
			m.processText(tea.WindowSizeMsg{Width: 40, Height: 12})

			interior := m.viewport.Height() - m.viewport.Style.GetVerticalFrameSize()

			// Content that fits must not leave phantom scroll range:
			// starting at offset 0 we are already at the bottom.
			if !m.viewport.AtBottom() {
				t.Error("fitted content should not be scrollable")
			}
			if got := m.viewport.TotalLineCount(); got != max(n, interior) {
				t.Errorf("expected %d total lines, got %d", max(n, interior), got)
			}
		})
	}
}
