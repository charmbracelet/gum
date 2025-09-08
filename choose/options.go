package choose

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options is the customization options for the choose command.
type Options struct {
	Options          []string      `arg:"" optional:"" help:"Options to choose from."`
	Limit            int           `help:"Maximum number of options to pick" default:"1" group:"Selection"`
	NoLimit          bool          `help:"Pick unlimited number of options (ignores limit)" group:"Selection"`
	Ordered          bool          `help:"Maintain the order of the selected options" env:"GUM_CHOOSE_ORDERED"`
	Height           int           `help:"Height of the list" default:"10" env:"GUM_CHOOSE_HEIGHT"`
	Cursor           string        `help:"Prefix to show on item that corresponds to the cursor position" default:"> " env:"GUM_CHOOSE_CURSOR"`
	ShowHelp         bool          `help:"Show help keybinds" default:"true" negatable:"" env:"GUM_CHOOSE_SHOW_HELP"`
	Timeout          time.Duration `help:"Timeout until choose returns selected element" default:"0s" env:"GUM_CHOOSE_TIMEOUT"` // including timeout command options [Timeout,...]
	Header           string        `help:"Header value" default:"Choose:" env:"GUM_CHOOSE_HEADER"`
	CursorPrefix     string        `help:"Prefix to show on the cursor item (hidden if limit is 1)" default:"• " env:"GUM_CHOOSE_CURSOR_PREFIX"`
	SelectedPrefix   string        `help:"Prefix to show on selected items (hidden if limit is 1)" default:"✓ " env:"GUM_CHOOSE_SELECTED_PREFIX"`
	UnselectedPrefix string        `help:"Prefix to show on unselected items (hidden if limit is 1)" default:"• " env:"GUM_CHOOSE_UNSELECTED_PREFIX"`
	Selected         []string      `help:"Options that should start as selected (selects all if given *)" default:"" env:"GUM_CHOOSE_SELECTED"`
	SelectIfOne      bool          `help:"Select the given option if there is only one" group:"Selection"`
	InputDelimiter   string        `help:"Option delimiter when reading from STDIN" default:"\n" env:"GUM_CHOOSE_INPUT_DELIMITER"`
	OutputDelimiter  string        `help:"Option delimiter when writing to STDOUT" default:"\n" env:"GUM_CHOOSE_OUTPUT_DELIMITER"`
	LabelDelimiter   string        `help:"Allows to set a delimiter, so options can be set as label:value" default:"" env:"GUM_CHOOSE_LABEL_DELIMITER"`
	StripANSI        bool          `help:"Strip ANSI sequences when reading from STDIN" default:"true" negatable:"" env:"GUM_CHOOSE_STRIP_ANSI"`
	Padding          string        `help:"Padding" default:"${defaultPadding}" group:"Style Flags" env:"GUM_CHOOSE_PADDING"`

	CursorStyle       style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" envprefix:"GUM_CHOOSE_CURSOR_"`
	HeaderStyle       style.Styles `embed:"" prefix:"header." set:"defaultForeground=99" envprefix:"GUM_CHOOSE_HEADER_"`
	ItemStyle         style.Styles `embed:"" prefix:"item." hidden:"" envprefix:"GUM_CHOOSE_ITEM_"`
	SelectedItemStyle style.Styles `embed:"" prefix:"selected." set:"defaultForeground=212" envprefix:"GUM_CHOOSE_SELECTED_"`
}
