package date

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/exit"

	"github.com/fxtlabs/date"
)

// Run provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
func (o Options) Run() error {
	picker := basePicker()

	picker.prompt = o.Prompt
	picker.promptStyle = o.PromptStyle.ToLipgloss()
	picker.cursorTextStyle = o.CursorTextStyle.ToLipgloss()
	if value, err := date.ParseISO(o.Value); err == nil {
		picker.Date = value
	}
	p := tea.NewProgram(model{
		picker:      picker,
		aborted:     false,
		header:      o.Header,
		headerStyle: o.HeaderStyle.ToLipgloss(),
		timeout:     o.Timeout,
		hasTimeout:  o.Timeout > 0,
	}, tea.WithOutput(os.Stderr))
	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run input: %w", err)
	}
	m := tm.(model)

	if m.aborted {
		return exit.ErrAborted
	}

	fmt.Println(m.picker.Value())
	return nil
}
