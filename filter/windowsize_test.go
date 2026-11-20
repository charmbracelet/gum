package filter

import (
	"testing"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func windowSizeModel(size tea.WindowSizeMsg, height int) model {
	v := viewport.New(viewport.WithWidth(size.Width), viewport.WithHeight(height))
	return model{
		textinput:   textinput.New(),
		viewport:    &v,
		headerStyle: lipgloss.NewStyle(),
		height:      height,
		padding:     []int{0, 0, 0, 0},
		keymap:      defaultKeymap(),
		help:        help.New(),
	}
}

func TestWindowSizeMsgIdempotent(t *testing.T) {
	tests := []struct {
		name      string
		height    int
		header    string
		showHelp  bool
		padding   []int
		size      tea.WindowSizeMsg
		wantFirst int
	}{
		{
			name:      "fixed height within window",
			height:    12,
			size:      tea.WindowSizeMsg{Width: 80, Height: 40},
			wantFirst: 11,
		},
		{
			name:      "fixed height equal to window",
			height:    40,
			size:      tea.WindowSizeMsg{Width: 80, Height: 40},
			wantFirst: 39,
		},
		{
			name:      "fixed height with padding",
			height:    12,
			padding:   []int{2, 0, 2, 0},
			size:      tea.WindowSizeMsg{Width: 80, Height: 40},
			wantFirst: 7,
		},
		{
			name:      "fixed height with header and padding",
			height:    12,
			header:    "header",
			padding:   []int{2, 0, 2, 0},
			size:      tea.WindowSizeMsg{Width: 80, Height: 40},
			wantFirst: 6,
		},
		{
			name:      "auto height uses window height",
			height:    0,
			size:      tea.WindowSizeMsg{Width: 80, Height: 40},
			wantFirst: 39,
		},
		{
			name:      "fixed height clamps to window",
			height:    60,
			size:      tea.WindowSizeMsg{Width: 80, Height: 40},
			wantFirst: 39,
		},
		{
			name:      "fixed height with show help",
			height:    12,
			showHelp:  true,
			size:      tea.WindowSizeMsg{Width: 80, Height: 40},
			wantFirst: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := windowSizeModel(tt.size, tt.height)
			m.header = tt.header
			m.showHelp = tt.showHelp
			m.padding = tt.padding
			if m.padding == nil {
				m.padding = []int{0, 0, 0, 0}
			}

			first, _ := m.Update(tt.size)
			m1 := first.(model)
			second, _ := m1.Update(tt.size)
			m2 := second.(model)

			if got := m1.viewport.Height(); got != tt.wantFirst {
				t.Errorf("first WindowSizeMsg: got height %d, want %d", got, tt.wantFirst)
			}
			if got := m2.viewport.Height(); got != m1.viewport.Height() {
				t.Errorf("repeated WindowSizeMsg not idempotent: got %d, want %d", got, m1.viewport.Height())
			}
		})
	}
}
