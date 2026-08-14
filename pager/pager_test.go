package pager

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func TestProcessTextViewportHeightAccountsForHelpAndSeparator(t *testing.T) {
	m := model{
		viewport:    viewport.New(80, 0),
		help:        help.New(),
		content:     strings.Repeat("line\n", 50),
		origContent: strings.Repeat("line\n", 50),
		keymap:      defaultKeymap(),
		softWrap:    false,
	}

	const termHeight = 24
	m.processText(tea.WindowSizeMsg{Height: termHeight, Width: 80})

	helpHeight := lipgloss.Height(m.helpView())
	want := termHeight - helpHeight - 1
	if m.viewport.Height != want {
		t.Fatalf("viewport.Height = %d, want %d (terminal=%d help=%d)", m.viewport.Height, want, termHeight, helpHeight)
	}

	viewHeight := lipgloss.Height(m.View())
	if viewHeight > termHeight {
		t.Fatalf("rendered view height %d exceeds terminal height %d", viewHeight, termHeight)
	}
}
