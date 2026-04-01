package confirm

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for prompting a user to confirm an
// action with an affirmative or negative answer.
func (o Options) Run() error {
	line, err := stdin.Read(stdin.SingleLine(true))
	if err == nil {
		switch line {
		case "yes", "y":
			return nil
		default:
			return exit.ErrExit(1)
		}
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	// Filter out empty extra options
	var extra []string
	for _, e := range o.Extra {
		if e != "" {
			extra = append(extra, e)
		}
	}

	defaultSelected := 0
	if !o.Default {
		defaultSelected = 1
	}

	top, right, bottom, left := style.ParsePadding(o.Padding)
	m := model{
		affirmative:      o.Affirmative,
		negative:         o.Negative,
		extra:            extra,
		showOutput:       o.ShowOutput,
		confirmation:     o.Default,
		selected:         defaultSelected,
		defaultSelection: o.Default,
		keys:             defaultKeymap(o.Affirmative, o.Negative),
		help:             help.New(),
		showHelp:         o.ShowHelp,
		prompt:           o.Prompt,
		selectedStyle:    o.SelectedStyle.ToLipgloss(),
		unselectedStyle:  o.UnselectedStyle.ToLipgloss(),
		promptStyle:      o.PromptStyle.ToLipgloss(),
		padding:          []int{top, right, bottom, left},
	}
	tm, err := tea.NewProgram(
		m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	).Run()
	if err != nil && ctx.Err() != context.DeadlineExceeded {
		return fmt.Errorf("unable to confirm: %w", err)
	}
	m = tm.(model)

	if o.ShowOutput {
		var confirmationText string
		switch {
		case m.selected == 0:
			confirmationText = m.affirmative
		case m.selected == 1 && m.negative != "":
			confirmationText = m.negative
		default:
			baseIdx := 1
			if m.negative != "" {
				baseIdx = 2
			}
			extraIdx := m.selected - baseIdx
			if extraIdx >= 0 && extraIdx < len(extra) {
				confirmationText = extra[extraIdx]
			}
		}
		fmt.Println(m.prompt, confirmationText)
	}

	// Exit code: 0 = affirmative, 1 = negative, 2+ = extra options
	if m.selected == 0 {
		return nil
	}

	return exit.ErrExit(m.selected)
}
