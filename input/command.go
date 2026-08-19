package input

import (
	"errors"
	"fmt"
	"os"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/gum/v2/cursor"
	"charm.land/gum/v2/internal/stdin"
	"charm.land/gum/v2/internal/timeout"
	"charm.land/gum/v2/style"
)

// Run provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
func (o Options) Run() error {
	if o.Value == "" {
		if in, _ := stdin.Read(stdin.StripANSI(o.StripANSI)); in != "" {
			o.Value = in
		}
	}

	i := textinput.New()
	if o.Value != "" {
		i.SetValue(o.Value)
	} else if in, _ := stdin.Read(stdin.StripANSI(o.StripANSI)); in != "" {
		i.SetValue(in)
	}
	i.Focus()
	i.Prompt = o.Prompt
	i.Placeholder = o.Placeholder
	i.SetWidth(o.Width)
	styles := i.Styles()
	styles.Focused.Prompt = o.PromptStyle.ToLipgloss()
	styles.Blurred.Prompt = o.PromptStyle.ToLipgloss()
	styles.Focused.Placeholder = o.PlaceholderStyle.ToLipgloss()
	styles.Blurred.Placeholder = o.PlaceholderStyle.ToLipgloss()
	i.SetStyles(styles)
	cursor.TextInput(&i, o.CursorMode, o.CursorStyle)
	i.CharLimit = o.CharLimit

	if o.Password {
		i.EchoMode = textinput.EchoPassword
		i.EchoCharacter = '•'
	}

	top, right, bottom, left := style.ParsePadding(o.Padding)
	m := model{
		textinput:   i,
		header:      o.Header,
		headerStyle: o.HeaderStyle.ToLipgloss(),
		padding:     []int{top, right, bottom, left},
		autoWidth:   o.Width < 1,
		showHelp:    o.ShowHelp,
		help:        help.New(),
		keymap:      defaultKeymap(),
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	p := tea.NewProgram(
		m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	)
	tm, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run input: %w", err)
	}

	m = tm.(model)
	if !m.submitted {
		return errors.New("not submitted")
	}
	fmt.Println(m.textinput.Value())
	return nil
}
