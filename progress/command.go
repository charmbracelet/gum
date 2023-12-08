// Package progress provides a simple progress indicator
// for tracking the progress for input provided via stdin.
//
// It shows a progress bar when the limit is known and some simple stats when not.
//
// ------------------------------------
// #!/bin/bash
//
// urls=(
//
//	"http://example.com/file1.txt"
//	"http://example.com/file2.txt"
//	"http://example.com/file3.txt"
//
// )
//
// for url in "${urls[@]}"; do
//
//	wget -q -nc "$url"
//	echo "Downloaded: $url"
//
// done | gum progress --show-output --limit ${#urls[@]}
// ------------------------------------
package progress

import (
	"bufio"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
)

func (o Options) GetFormatString() string {
	if o.Format != "" {
		return o.Format
	}

	switch {
	case o.Limit == 0 && o.Title == "":
		return "[Elapsed ~ {Elapsed}] Iter {Iter}"
	case o.Limit == 0 && o.Title != "":
		return "[Elapsed ~ {Elapsed}] Iter {Iter} ~ {Title}"
	case o.Limit > 0 && o.Title == "":
		return "{Bar} {Pct}"
	case o.Limit > 0 && o.Title != "":
		return "{Title} ~ {Bar} {Pct}"
	default:
		return "{Iter}"
	}
}

func (o Options) Run() error {
	m := &model{
		reader:                bufio.NewReader(os.Stdin),
		output:                o.ShowOutput,
		isTTY:                 isatty.IsTerminal(os.Stdout.Fd()),
		progressIndicator:     o.ProgressIndicator,
		hideProgressIndicator: o.HideProgressIndicator,

		bfmt:  newBarFormatter(o.GetFormatString(), o.ProgressColor),
		binfo: newBarInfo(o.TitleStyle.ToLipgloss().Render(o.Title), o.Limit),
	}
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run progress: %w", err)
	}

	return m.err
}
