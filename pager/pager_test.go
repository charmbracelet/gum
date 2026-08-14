package pager

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
)

func TestProcessTextPreservesViewportHeightAcrossMatchNavigation(t *testing.T) {
	m := model{
		viewport:        viewport.New(80, 24),
		help:            help.New(),
		content:         strings.Repeat("match line\n", 200),
		showLineNumbers: true,
		keymap:          defaultKeymap(),
		width:           80,
		height:          24,
	}

	m.processText()
	initialHeight := m.viewport.Height

	for range 5 {
		m.processText()
	}

	if m.viewport.Height != initialHeight {
		t.Fatalf("viewport height drifted from %d to %d after repeated processText calls", initialHeight, m.viewport.Height)
	}
}

func TestAlignViewportUsesRenderedContentLine(t *testing.T) {
	m := model{
		viewport:        viewport.New(40, 10),
		help:            help.New(),
		content:         strings.Repeat("filler\n", 50) + "match here",
		showLineNumbers: false,
		softWrap:        false,
		keymap:          defaultKeymap(),
		width:           40,
		height:          10,
	}

	rendered := m.processText()
	m.viewport.SetYOffset(40)
	m.search.matchLipglossStr = "match here"
	m.search.alignViewport(&m, rendered)

	if m.viewport.YOffset <= 40 {
		t.Fatalf("expected viewport to scroll toward match below view, got %d", m.viewport.YOffset)
	}
}
