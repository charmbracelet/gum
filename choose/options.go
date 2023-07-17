package choose

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options is the customization options for the choose command.
type Options struct {
	Options           []string      `arg:"" optional:"" help:"Options to choose from."`
	Limit             int           `help:"Maximum number of options to pick" default:"1" group:"Selection"`
	NoLimit           bool          `help:"Pick unlimited number of options (ignores limit)" group:"Selection"`
	Ordered           bool          `help:"Maintain the order of the selected options" env:"GUM_CHOOSE_ORDERED"`
	Height            int           `help:"Height of the list" default:"10" env:"GUM_CHOOSE_HEIGHT"`
	Cursor            string        `help:"Prefix to show on item that corresponds to the cursor position" default:"> " env:"GUM_CHOOSE_CURSOR"`
	Header            string        `help:"Header value" default:"" env:"GUM_CHOOSE_HEADER"`
	CursorPrefix      string        `help:"Prefix to show on the cursor item (hidden if limit is 1)" default:"○ " env:"GUM_CHOOSE_CURSOR_PREFIX"`
	SelectedPrefix    string        `help:"Prefix to show on selected items (hidden if limit is 1)" default:"◉ " env:"GUM_CHOOSE_SELECTED_PREFIX"`
	UnselectedPrefix  string        `help:"Prefix to show on unselected items (hidden if limit is 1)" default:"○ " env:"GUM_CHOOSE_UNSELECTED_PREFIX"`
	Selected          []string      `help:"Options that should start as selected" default:"" env:"GUM_CHOOSE_SELECTED"`
	SelectIfOne       bool          `help:"Select the given option if there is only one" group:"Selection"`
	CursorStyle       style.Styles  `embed:"" prefix:"cursor." set:"defaultForeground=212" envprefix:"GUM_CHOOSE_CURSOR_"`
	HeaderStyle       style.Styles  `embed:"" prefix:"header." set:"defaultForeground=240" envprefix:"GUM_CHOOSE_HEADER_"`
	ItemStyle         style.Styles  `embed:"" prefix:"item." hidden:"" envprefix:"GUM_CHOOSE_ITEM_"`
	SelectedItemStyle style.Styles  `embed:"" prefix:"selected." set:"defaultForeground=212" envprefix:"GUM_CHOOSE_SELECTED_"`
	Timeout           time.Duration `help:"Timeout until choose returns selected element" default:"0" env:"GUM_CCHOOSE_TIMEOUT"` // including timeout command options [Timeout,...]
}
