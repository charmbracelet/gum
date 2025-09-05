package table

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Run provides a shell script interface for rendering tabular data (CSV).
func (o Options) Run() error {
	var input *os.File
	if o.File != "" {
		var err error
		input, err = os.Open(o.File)
		if err != nil {
			return fmt.Errorf("could not render file: %w", err)
		}
	} else {
		if stdin.IsEmpty() {
			return fmt.Errorf("no data provided")
		}
		input = os.Stdin
	}
	defer input.Close() //nolint: errcheck

	transformer := unicode.BOMOverride(encoding.Nop.NewDecoder())
	reader := csv.NewReader(transform.NewReader(input, transformer))
	reader.LazyQuotes = o.LazyQuotes
	reader.FieldsPerRecord = o.FieldsPerRecord
	separatorRunes := []rune(o.Separator)
	if len(separatorRunes) != 1 {
		return fmt.Errorf("separator must be single character")
	}
	reader.Comma = separatorRunes[0]

	writer := csv.NewWriter(os.Stdout)
	writer.Comma = separatorRunes[0]

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
	columns := make([]table.Column, 0, len(columnNames))

	for i, title := range columnNames {
		width := lipgloss.Width(title)
		if len(o.Widths) > i {
			width = o.Widths[i]
		}
		columns = append(columns, table.Column{
			Title: title,
			Width: width,
		})
	}

	defaultStyles := table.DefaultStyles()
	top, right, bottom, left := style.ParsePadding(o.Padding)

	styles := table.Styles{
		Cell:     defaultStyles.Cell.Inherit(o.CellStyle.ToLipgloss()),
		Header:   defaultStyles.Header.Inherit(o.HeaderStyle.ToLipgloss()),
		Selected: o.SelectedStyle.ToLipgloss(),
	}

	rows := make([]table.Row, 0, len(data))
	for row := range data {
		if len(data[row]) > len(columns) {
			return fmt.Errorf("invalid number of columns")
		}

		// fixes the data in case we have more columns than rows:
		for len(data[row]) < len(columns) {
			data[row] = append(data[row], "")
		}

		for i, col := range data[row] {
			if len(o.Widths) == 0 {
				width := lipgloss.Width(col)
				if width > columns[i].Width {
					columns[i].Width = width
				}
			}
		}

		rows = append(rows, table.Row(data[row]))
	}

	if o.Print {
		table := ltable.New().
			Headers(columnNames...).
			Rows(data...).
			BorderStyle(o.BorderStyle.ToLipgloss()).
			Border(style.Border[o.Border]).
			StyleFunc(func(row, _ int) lipgloss.Style {
				if row == 0 {
					return styles.Header
				}
				return styles.Cell
			})

		fmt.Println(table.Render())
		return nil
	}

	opts := []table.Option{
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithRows(rows),
		table.WithStyles(styles),
	}
	if o.Height > 0 {
		opts = append(opts, table.WithHeight(o.Height-top-bottom))
	}

	table := table.New(opts...)

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	m := model{
		table:     table,
		showHelp:  o.ShowHelp,
		hideCount: o.HideCount,
		help:      help.New(),
		keymap:    defaultKeymap(),
		padding:   []int{top, right, bottom, left},
	}
	tm, err := tea.NewProgram(
		m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return fmt.Errorf("failed to start tea program: %w", err)
	}

	if tm == nil {
		return fmt.Errorf("failed to get selection")
	}

	m = tm.(model)
	if o.ReturnColumn > 0 && o.ReturnColumn <= len(m.selected) {
		if err = writer.Write([]string{m.selected[o.ReturnColumn-1]}); err != nil {
			return fmt.Errorf("failed to write col %d of selected row: %w", o.ReturnColumn, err)
		}
	} else {
		if err = writer.Write([]string(m.selected)); err != nil {
			return fmt.Errorf("failed to write selected row: %w", err)
		}
	}

	writer.Flush()

	return nil
}
