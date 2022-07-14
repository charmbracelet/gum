package choose

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the choose command.
type Options struct {
	Options []string `arg:"" optional:"" help:"Options to choose from."`

	Limit             int          `help:"Maximum number of options to pick" default:"1" group:"Selection"`
	NoLimit           bool         `help:"Pick unlimited number of options (ignores limit)" group:"Selection"`
	Height            int          `help:"Height of the list" default:"10"`
	Cursor            string       `help:"Prefix to show on item that corresponds to the cursor position" default:"> "`
	CursorPrefix      string       `help:"Prefix to show on the cursor item (hidden if limit is 1)" default:"[•] "`
	SelectedPrefix    string       `help:"Prefix to show on selected items (hidden if limit is 1)" default:"[✕] "`
	UnselectedPrefix  string       `help:"Prefix to show on selected items (hidden if limit is 1)" default:"[ ] "`
	CursorStyle       style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" set:"name=indicator"`
	ItemStyle         style.Styles `embed:"" prefix:"item." hidden:"" set:"defaultForeground=255" set:"name=item"`
	SelectedItemStyle style.Styles `embed:"" prefix:"selected." set:"defaultForeground=212" set:"name=selected item"`
}
