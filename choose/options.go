package choose

// Options is the customization options for the choose command.
type Options struct {
	Options []string `arg:"" optional:"" help:"Options to choose from."`

	Height          int    `help:"Height of the list." default:"10"`
	HidePagination  bool   `help:"Hide pagination." default:"false"`
	Indicator       string `help:"Prefix to show on selected item" default:"> "`
	IndicatorColor  string `help:"Indicator foreground color" default:"212"`
	SelectedColor   string `help:"Selected item foreground color" default:"212"`
	UnselectedColor string `help:"Unselected item foreground color" default:""`
}
