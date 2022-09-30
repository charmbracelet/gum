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
	Background string `help:"Background Color" default:"${defaultBackground}" group:"Style Flags" env:"BACKGROUND"`
	Foreground string `help:"Foreground Color" default:"${defaultForeground}" group:"Style Flags" env:"FOREGROUND"`

	// Border
	Border           string `help:"Border Style" enum:"none,hidden,normal,rounded,thick,double" default:"${defaultBorder}" group:"Style Flags" env:"BORDER"`
	BorderBackground string `help:"Border Background Color" group:"Style Flags" env:"BORDER_BACKGROUND"`
	BorderForeground string `help:"Border Foreground Color" group:"Style Flags" default:"${defaultBorderForeground}" env:"BORDER_FOREGROUND"`

	// Layout
	Align   string `help:"Text Alignment" enum:"left,center,right,bottom,middle,top" default:"left" group:"Style Flags" env:"ALIGN"`
	Height  int    `help:"Text height" group:"Style Flags" env:"HEIGHT"`
	Width   int    `help:"Text width" group:"Style Flags" env:"WIDTH"`
	Margin  string `help:"Text margin" default:"${defaultMargin}" group:"Style Flags" env:"MARGIN"`
	Padding string `help:"Text padding" default:"${defaultPadding}" group:"Style Flags" env:"PADDING"`

	// Format
	Bold          bool `help:"Bold text" group:"Style Flags" env:"BOLD"`
	Faint         bool `help:"Faint text" group:"Style Flags" env:"FAINT"`
	Italic        bool `help:"Italicize text" group:"Style Flags" env:"ITALIC"`
	Strikethrough bool `help:"Strikethrough text" group:"Style Flags" env:"STRIKETHROUGH"`
	Underline     bool `help:"Underline text" default:"${defaultUnderline}" group:"Style Flags" env:"UNDERLINE"`
}
