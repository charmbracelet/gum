package table

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"

	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/style"
)

// this solves scenarios where a comma may be between quotes
func splitAtDelim(s, delim string) []string {
	var res []string
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if string(s[i]) == delim && !inString {
			res = append(res, s[beg:i])
			beg = i+1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

// Run provides a shell script interface for rendering tabular data (CSV)
func (o Options) Run() error {
	csv, err := stdin.Read()
	if err != nil {
		return fmt.Errorf("no data provided: %w", err)
	}

	if csv == "" {
		return fmt.Errorf("no data provided")
	}

	// If no columns are provided we'll use the first row of the CSV as the
	// column names.
	lines := strings.Split(csv, "\n")
	if len(o.Columns) <= 0 {
		if len(lines) > 0 {
			o.Columns = splitAtDelim(lines[0], o.Separator)
			lines = lines[1:]
		} else {
			return fmt.Errorf("no columns provided")
		}
	}

	var columns []table.Column

	for i, title := range o.Columns {
		width := runewidth.StringWidth(title)
		if len(o.Widths) > i {
			width = o.Widths[i]
		}
		columns = append(columns, table.Column{
			Title: title,
			Width: width,
		})
	}

	var rows []table.Row

	for _, line := range lines {
		if line == "" {
			continue
		}
		row := splitAtDelim(line, o.Separator)
		if len(row) != len(columns) {
			return fmt.Errorf("row %q has %d columns, expected %d", line, len(row), len(columns))
		}
		rows = append(rows, row)
	}

	defaultStyles := table.DefaultStyles()

	styles := table.Styles{
		Cell:     defaultStyles.Cell.Inherit(o.CellStyle.ToLipgloss()),
		Header:   defaultStyles.Header.Inherit(o.HeaderStyle.ToLipgloss()),
		Selected: defaultStyles.Selected.Inherit(o.SelectedStyle.ToLipgloss()),
	}

	table := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(o.Height),
		table.WithRows(rows),
		table.WithStyles(styles),
	)

	tm, err := tea.NewProgram(model{table: table}, tea.WithOutput(os.Stderr)).StartReturningModel()

	if err != nil {
		return fmt.Errorf("failed to start tea program: %w", err)
	}

	m := tm.(model)
	fmt.Println(strings.Join([]string(m.selected), o.Separator))

	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}
