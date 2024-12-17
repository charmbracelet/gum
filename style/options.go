package style

// Options is the customization options for the style command.
type Options struct {
	Text      []string        `arg:"" optional:"" help:"Text to which to apply the style"`
	Trim      bool            `help:"Trim whitespaces on every input line" default:"false"`
	StripANSI bool            `help:"Strip ANSI sequences when reading from STDIN" default:"true" negatable:"" env:"GUM_STYLE_STRIP_ANSI"`
	Style     StylesNotHidden `embed:""`
}

// Styles is a flag set of possible styles.
//
// It corresponds to the available options in the lipgloss.Style struct.
//
// This flag set is used in other parts of the application to embed styles for
// components, through embedding and prefixing.
type Styles struct {
	// Colors
	Foreground string `help:"Foreground Color" default:"${defaultForeground}" group:"Style Flags" env:"FOREGROUND"`
	Background string `help:"Background Color" default:"${defaultBackground}" group:"Style Flags" env:"BACKGROUND"`

	// Border
	Border           string `help:"Border Style" enum:"none,hidden,normal,rounded,thick,double" default:"${defaultBorder}" group:"Style Flags" env:"BORDER" hidden:"true"`
	BorderBackground string `help:"Border Background Color" group:"Style Flags" default:"${defaultBorderBackground}" env:"BORDER_BACKGROUND" hidden:"true"`
	BorderForeground string `help:"Border Foreground Color" group:"Style Flags" default:"${defaultBorderForeground}" env:"BORDER_FOREGROUND" hidden:"true"`

	// Layout
	Align   string `help:"Text Alignment" enum:"left,center,right,bottom,middle,top" default:"${defaultAlign}" group:"Style Flags" env:"ALIGN" hidden:"true"`
	Height  int    `help:"Text height" default:"${defaultHeight}" group:"Style Flags" env:"HEIGHT" hidden:"true"`
	Width   int    `help:"Text width" default:"${defaultWidth}" group:"Style Flags" env:"WIDTH" hidden:"true"`
	Margin  string `help:"Text margin" default:"${defaultMargin}" group:"Style Flags" env:"MARGIN" hidden:"true"`
	Padding string `help:"Text padding" default:"${defaultPadding}" group:"Style Flags" env:"PADDING" hidden:"true"`

	// Format
	Bold          bool `help:"Bold text" default:"${defaultBold}" group:"Style Flags" env:"BOLD" hidden:"true"`
	Faint         bool `help:"Faint text" default:"${defaultFaint}" group:"Style Flags" env:"FAINT" hidden:"true"`
	Italic        bool `help:"Italicize text" default:"${defaultItalic}" group:"Style Flags" env:"ITALIC" hidden:"true"`
	Strikethrough bool `help:"Strikethrough text" default:"${defaultStrikethrough}" group:"Style Flags" env:"STRIKETHROUGH" hidden:"true"`
	Underline     bool `help:"Underline text" default:"${defaultUnderline}" group:"Style Flags" env:"UNDERLINE" hidden:"true"`
}

// StylesNotHidden allows the style struct to display full help when not-embedded.
//
// NB: We must duplicate this struct to ensure that `gum style` does not hide
// flags when an error pops up. Ideally, we can dynamically hide or show flags
// based on the command run: https://github.com/alecthomas/kong/issues/316
type StylesNotHidden struct {
	// Colors
	Foreground string `help:"Foreground Color" default:"${defaultForeground}" group:"Style Flags" env:"FOREGROUND"`
	Background string `help:"Background Color" default:"${defaultBackground}" group:"Style Flags" env:"BACKGROUND"`

	// Border
	Border           string `help:"Border Style" enum:"none,hidden,normal,rounded,thick,double" default:"${defaultBorder}" group:"Style Flags" env:"BORDER"`
	BorderBackground string `help:"Border Background Color" group:"Style Flags" default:"${defaultBorderBackground}" env:"BORDER_BACKGROUND"`
	BorderForeground string `help:"Border Foreground Color" group:"Style Flags" default:"${defaultBorderForeground}" env:"BORDER_FOREGROUND"`

	// Layout
	Align   string `help:"Text Alignment" enum:"left,center,right,bottom,middle,top" default:"${defaultAlign}" group:"Style Flags" env:"ALIGN"`
	Height  int    `help:"Text height" default:"${defaultHeight}" group:"Style Flags" env:"HEIGHT"`
	Width   int    `help:"Text width" default:"${defaultWidth}" group:"Style Flags" env:"WIDTH"`
	Margin  string `help:"Text margin" default:"${defaultMargin}" group:"Style Flags" env:"MARGIN"`
	Padding string `help:"Text padding" default:"${defaultPadding}" group:"Style Flags" env:"PADDING"`

	// Format
	Bold          bool `help:"Bold text" default:"${defaultBold}" group:"Style Flags" env:"BOLD"`
	Faint         bool `help:"Faint text" default:"${defaultFaint}" group:"Style Flags" env:"FAINT"`
	Italic        bool `help:"Italicize text" default:"${defaultItalic}" group:"Style Flags" env:"ITALIC"`
	Strikethrough bool `help:"Strikethrough text" default:"${defaultStrikethrough}" group:"Style Flags" env:"STRIKETHROUGH"`
	Underline     bool `help:"Underline text" default:"${defaultUnderline}" group:"Style Flags" env:"UNDERLINE"`
}
