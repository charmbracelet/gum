package table

import (
	"encoding/csv"
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

// Run provides a shell script interface for rendering tabular data (CSV).
func (o Options) Run() error {
	var reader *csv.Reader
	if o.File != "" {
		file, err := os.Open(o.File)
		if err != nil {
			return fmt.Errorf("could not find file at path %s", o.File)
		}
		reader = csv.NewReader(file)
	} else {
		if stdin.IsEmpty() {
			return fmt.Errorf("no data provided")
		}
		reader = csv.NewReader(os.Stdin)
	}

	separatorRunes := []rune(o.Separator)
	if len(separatorRunes) != 1 {
		return fmt.Errorf("separator must be single character")
	}
	reader.Comma = separatorRunes[0]

	var columnNames []string
	var err error
	// If no columns are provided we'll use the first row of the CSV as the
	// column names.
	if len(o.Columns) <= 0 {
		columnNames, err = reader.Read()
		if err != nil {
			return fmt.Errorf("unable to parse columns")
		}
	} else {
		columnNames = o.Columns
	}

	data, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("invalid data provided")
	}
	var columns = make([]table.Column, 0, len(columnNames))

	for i, title := range columnNames {
		width := runewidth.StringWidth(title)
		if len(o.Widths) > i {
			width = o.Widths[i]
		}
		columns = append(columns, table.Column{
			Title: title,
			Width: width,
		})
	}

	defaultStyles := table.DefaultStyles()

	styles := table.Styles{
		Cell:     defaultStyles.Cell.Inherit(o.CellStyle.ToLipgloss()),
		Header:   defaultStyles.Header.Inherit(o.HeaderStyle.ToLipgloss()),
		Selected: defaultStyles.Selected.Inherit(o.SelectedStyle.ToLipgloss()),
	}

	var rows = make([]table.Row, 0, len(data))
	for _, row := range data {
		rows = append(rows, table.Row(row))
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

	if tm == nil {
		return fmt.Errorf("failed to get selection")
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
