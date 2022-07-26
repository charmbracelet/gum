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
	Background string `help:"Background Color" default:"${defaultBackground}" group:"Style Flags"`
	Foreground string `help:"Foreground Color" default:"${defaultForeground}" group:"Style Flags"`

	// Border
	Border           string `help:"Border Style" enum:"none,hidden,normal,rounded,thick,double" default:"none" group:"Style Flags"`
	BorderBackground string `help:"Border Background Color" group:"Style Flags"`
	BorderForeground string `help:"Border Foreground Color" group:"Style Flags"`

	// Layout
	Align   string `help:"Text Alignment" enum:"left,center,right,bottom,middle,top" default:"left" group:"Style Flags"`
	Height  int    `help:"Text height" group:"Style Flags"`
	Width   int    `help:"Text width" group:"Style Flags"`
	Margin  string `help:"Text margin" default:"${defaultMargin}" group:"Style Flags"`
	Padding string `help:"Text padding" default:"${defaultPadding}" group:"Style Flags"`

	// Format
	Bold          bool `help:"Bold text" group:"Style Flags"`
	Faint         bool `help:"Faint text" group:"Style Flags"`
	Italic        bool `help:"Italicize text" group:"Style Flags"`
	Strikethrough bool `help:"Strikethrough text" group:"Style Flags"`
	Underline     bool `help:"Underline text" default:"${defaultUnderline}" group:"Style Flags"`
}
