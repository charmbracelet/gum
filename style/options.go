package style

const (
	groupName = "Style Flags"
)

// Options is the customization options for the style command.
type Options struct {
	Text  []string `arg:"" optional:"" help:"Text to which to apply the style"`
	Style Styles   `embed:""`
}

// Styles is a flag set of possible styles.
//
// It corresponds to the available options in the lipgloss.Style struct.
//
// This flag set is used in other parts of the application to embed styles for
// components, through embedding and prefixing.
type Styles struct {
	// Colors
	Background string `help:"Background color of the ${name=element}" default:"${defaultBackground}" group:"Style Flags"`
	Foreground string `help:"color of the ${name=element}" default:"${defaultForeground}" group:"Style Flags"`

	// Border
	Border           string `help:"Border style to apply" enum:"none,hidden,normal,rounded,thick,double" default:"none" group:"Style Flags"`
	BorderBackground string `help:"Border background color" group:"Style Flags"`
	BorderForeground string `help:"Border foreground color" group:"Style Flags"`

	// Layout
	Align   string `help:"Text alignment" enum:"left,center,right,bottom,middle,top" default:"left" group:"Style Flags"`
	Height  int    `help:"Height of output" group:"Style Flags"`
	Width   int    `help:"Width of output" group:"Style Flags"`
	Margin  string `help:"Margin to apply around the text." default:"0 0" group:"Style Flags"`
	Padding string `help:"Padding to apply around the text." default:"0 0"`

	// Format
	Bold          bool `help:"Apply bold formatting" group:"Style Flags"`
	Faint         bool `help:"Apply faint formatting" group:"Style Flags"`
	Italic        bool `help:"Apply italic formatting" group:"Style Flags"`
	Strikethrough bool `help:"Apply strikethrough formatting" group:"Style Flags"`
}
