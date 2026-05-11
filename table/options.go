package table

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options is the customization options for the table command.
type Options struct {
	Separator       string   `short:"s" help:"Row separator" default:","`
	Columns         []string `short:"c" help:"Column names"`
	Widths          []int    `short:"w" help:"Column widths"`
	Height          int      `help:"Table height" default:"0"`
	Print           bool     `short:"p" help:"static print" default:"false"`
	File            string   `short:"f" help:"file path" default:""`
	Border          string   `short:"b" help:"border style" default:"rounded" enum:"rounded,thick,normal,hidden,double,none"`
	ShowHelp        bool     `help:"Show help keybinds" default:"true" negatable:"" env:"GUM_TABLE_SHOW_HELP"`
	HideCount       bool     `help:"Hide item count on help keybinds" default:"false" negatable:"" env:"GUM_TABLE_HIDE_COUNT"`
	LazyQuotes      bool     `help:"If LazyQuotes is true, a quote may appear in an unquoted field and a non-doubled quote may appear in a quoted field" default:"false" env:"GUM_TABLE_LAZY_QUOTES"`
	FieldsPerRecord int      `help:"Sets the number of expected fields per record" default:"0" env:"GUM_TABLE_FIELDS_PER_RECORD"`

	BorderStyle   style.Styles  `embed:"" prefix:"border." envprefix:"GUM_TABLE_BORDER_"`
	CellStyle     style.Styles  `embed:"" prefix:"cell." envprefix:"GUM_TABLE_CELL_"`
	HeaderStyle   style.Styles  `embed:"" prefix:"header." envprefix:"GUM_TABLE_HEADER_"`
	SelectedStyle style.Styles  `embed:"" prefix:"selected." set:"defaultForeground=212" envprefix:"GUM_TABLE_SELECTED_"`
	ReturnColumn  int           `short:"r" help:"Which column number should be returned instead of whole row as string. Default=0 returns whole Row" default:"0"`
	Timeout       time.Duration `help:"Timeout until choose returns selected element" default:"0s" env:"GUM_TABLE_TIMEOUT"`
	Padding       string        `help:"Padding" default:"${defaultPadding}" group:"Style Flags" env:"GUM_TABLE_PADDING"`
}
