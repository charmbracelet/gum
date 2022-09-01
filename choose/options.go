package choose

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the choose command.
type Options struct {
	Options []string `arg:"" optional:"" help:"Options to choose from."`

	Limit             int          `help:"Maximum number of options to pick" default:"1" group:"Selection"`
	NoLimit           bool         `help:"Pick unlimited number of options (ignores limit)" group:"Selection"`
	Optional          bool         `help:"Do not force highlighted option if none is selected" group:"Selection"`
	Height            int          `help:"Height of the list" default:"10" env:"GUM_CHOOSE_HEIGHT"`
	Cursor            string       `help:"Prefix to show on item that corresponds to the cursor position" default:"> " env:"GUM_CHOOSE_CURSOR"`
	CursorPrefix      string       `help:"Prefix to show on the cursor item (hidden if limit is 1)" default:"[•] " env:"GUM_CHOOSE_CURSOR_PREFIX"`
	SelectedPrefix    string       `help:"Prefix to show on selected items (hidden if limit is 1)" default:"[✕] " env:"GUM_CHOOSE_SELECTED_PREFIX"`
	UnselectedPrefix  string       `help:"Prefix to show on selected items (hidden if limit is 1)" default:"[ ] " env:"GUM_CHOOSE_UNSELECTED_PREFIX"`
	CursorStyle       style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" envprefix:"GUM_CHOOSE_CURSOR_"`
	ItemStyle         style.Styles `embed:"" prefix:"item." hidden:"" envprefix:"GUM_CHOOSE_ITEM_"`
	SelectedItemStyle style.Styles `embed:"" prefix:"selected." set:"defaultForeground=212" envprefix:"GUM_CHOOSE_SELECTED_"`
}
