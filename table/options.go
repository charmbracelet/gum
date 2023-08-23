package table

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the table command.
type Options struct {
	Separator     string       `short:"s" help:"Row separator" default:","`
	Columns       []string     `short:"c" help:"Column names"`
	Widths        []int        `short:"w" help:"Column widths"`
	Height        int          `help:"Table height" default:"10"`
	CellStyle     style.Styles `embed:"" prefix:"cell." envprefix:"GUM_TABLE_CELL_"`
	HeaderStyle   style.Styles `embed:"" prefix:"header." envprefix:"GUM_TABLE_HEADER_"`
	SelectedStyle style.Styles `embed:"" prefix:"selected." set:"defaultForeground=212" envprefix:"GUM_TABLE_SELECTED_"`
	File          string       `short:"f" help:"file path" default:""`
	ReturnColumn  int          `short:"r" help:"Which column number should be returned instead of whole row as string. Default=0 returns whole Row" default:"0"`
}
