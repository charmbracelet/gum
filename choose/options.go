package choose

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the choose command.
type Options struct {
	Options []string `arg:"" optional:"" help:"Options to choose from."`

	Height         int          `help:"Height of the list" default:"10"`
	HidePagination bool         `help:"Hide pagination" default:"false"`
	Indicator      string       `help:"Prefix to show on selected item" default:"> "`
	IndicatorStyle style.Styles `embed:"" prefix:"indicator." set:"defaultForeground=212" set:"name=indicator"`
	ItemStyle      style.Styles `embed:"" prefix:"item." hidden:"" set:"defaultForeground=255" set:"name=item"`
	SelectedStyle  style.Styles `embed:"" prefix:"selected." set:"defaultForeground=212" set:"name=selected item"`
}
