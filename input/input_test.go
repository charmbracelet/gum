package input

import (
	"testing"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

func TestViewNoPanicOnZeroWidth(t *testing.T) {
	const placeholder = "Type something..."

	i := textinput.New()
	i.Focus()
	i.Prompt = "?"
	i.Placeholder = placeholder

	tm := model{
		autoWidth: true,
		padding:   []int{0, 0, 0, 0},
		textinput: i,
		showHelp:  true,
		help:      help.New(),
		keymap:    defaultKeymap(),
	}

	// script -qec with no controlling terminal reports 0x0, which used to
	// make autoWidth compute a negative width and panic in bubbles'
	// textinput placeholderView (make([]rune, Width()+1)).
	updated, _ := tm.Update(tea.WindowSizeMsg{Width: 0, Height: 0})
	_ = updated.View()
}
