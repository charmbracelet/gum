package input

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options are the customization options for the input.
type Options struct {
	Placeholder      string        `help:"Placeholder value" default:"Type something..." env:"GUM_INPUT_PLACEHOLDER"`
	Prompt           string        `help:"Prompt to display" default:"> " env:"GUM_INPUT_PROMPT"`
	PromptStyle      style.Styles  `embed:"" prefix:"prompt." envprefix:"GUM_INPUT_PROMPT_"`
	PlaceholderStyle style.Styles  `embed:"" prefix:"placeholder." set:"defaultForeground=240" envprefix:"GUM_INPUT_PLACEHOLDER_"`
	CursorStyle      style.Styles  `embed:"" prefix:"cursor." set:"defaultForeground=212" envprefix:"GUM_INPUT_CURSOR_"`
	CursorMode       string        `prefix:"cursor." name:"mode" help:"Cursor mode" default:"blink" enum:"blink,hide,static" env:"GUM_INPUT_CURSOR_MODE"`
	Value            string        `help:"Initial value (can also be passed via stdin)" default:""`
	CharLimit        int           `help:"Maximum value length (0 for no limit)" default:"400"`
	Width            int           `help:"Input width (0 for terminal width)" default:"0" env:"GUM_INPUT_WIDTH"`
	Password         bool          `help:"Mask input characters" default:"false"`
	ShowHelp         bool          `help:"Show help keybinds" default:"true" negatable:"" env:"GUM_INPUT_SHOW_HELP"`
	Header           string        `help:"Header value" default:"" env:"GUM_INPUT_HEADER"`
	HeaderStyle      style.Styles  `embed:"" prefix:"header." set:"defaultForeground=240" envprefix:"GUM_INPUT_HEADER_"`
	Timeout          time.Duration `help:"Timeout until input aborts" default:"0s" env:"GUM_INPUT_TIMEOUT"`
	StripANSI        bool          `help:"Strip ANSI sequences when reading from STDIN" default:"true" negatable:"" env:"GUM_INPUT_STRIP_ANSI"`
	Padding          string        `help:"Padding" default:"${defaultPadding}" group:"Style Flags" env:"GUM_INPUT_PADDING"`
}
