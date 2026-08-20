package panel

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/internal/tty"
)

// Run provides a shell script interface for running multiple panels
// side by side, allowing users to navigate between them.
func (o Options) Run() error {
	if len(o.Panel) == 0 {
		return os.ErrInvalid
	}

	panels, err := parsePanelsFromArgs(o.Panel, o.InputDelimiter, o.StripANSI)
	if err != nil {
		return err
	}

	height := o.Height
	if height <= 0 {
		height = 10
	}

	m := orchestrator{
		panels:    panels,
		activeIdx: o.Active,
		height:    height,
		borderStyle: o.BorderStyle(),
		showHelp:    o.ShowHelp,
		customHelp:  o.CustomHelp,
		gap:         o.Gap,
		vertical:    o.Vertical,
		reverse:     o.Reverse,
		stacked:     o.Stacked,
		delimiter:   o.Delimiter,
		debug:       o.Debug,
		all:         o.All,
		single:      o.Single,

		// Border styles
		activeBorderStyle:   o.ActiveBorderStyle.ToLipgloss(),
		inactiveBorderStyle: o.InactiveBorderStyle.ToLipgloss(),

		// Common styles
		matchStyle:        o.MatchStyle.ToLipgloss(),
		cursorStyle:       o.CursorStyle.ToLipgloss(),
		headerStyle:       o.HeaderStyle.ToLipgloss(),
		itemStyle:         o.ItemStyle.ToLipgloss(),
		selectedItemStyle: o.SelectedItemStyle.ToLipgloss(),
		indicatorStyle:    o.IndicatorStyle.ToLipgloss(),

		// Choose-specific styles
		selectedPrefixStyle:   o.SelectedPrefixStyle.ToLipgloss(),
		unselectedPrefixStyle: o.UnselectedPrefixStyle.ToLipgloss(),

		// Filter-specific styles
		textStyle:        o.TextStyle.ToLipgloss(),
		cursorTextStyle:  o.CursorTextStyle.ToLipgloss(),
		promptStyle:      o.PromptStyle.ToLipgloss(),
		placeholderStyle: o.PlaceholderStyle.ToLipgloss(),

		keymap: defaultPanelKeymap(),
		help:   help.New(),
	}

	if initErr := m.initModels(o); initErr != nil {
		return initErr
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	tm, err := tea.NewProgram(
		m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return err
	}

	result := tm.(orchestrator)

	if !result.submitted {
		return os.ErrInvalid
	}

	var finalOutput string
	if o.Single {
		finalOutput = result.getSingleResult()
	} else {
		output := result.getResults(o.OutputDelimiter)
		if o.Stacked {
			finalOutput = strings.Join(output, "\n"+o.Delimiter+"\n")
		} else {
			finalOutput = strings.Join(output, "\n")
		}
	}

	tty.Println(finalOutput)
	return nil
}
