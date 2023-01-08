package input

import (
	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/gum/timeout"
)

// Options are the customization options for the input.
type Options struct {
	Placeholder string       `help:"Placeholder value" default:"Type something..." env:"GUM_INPUT_PLACEHOLDER"`
	Prompt      string       `help:"Prompt to display" default:"> " env:"GUM_INPUT_PROMPT"`
	PromptStyle style.Styles `embed:"" prefix:"prompt." envprefix:"GUM_INPUT_PROMPT_"`
	CursorStyle style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" envprefix:"GUM_INPUT_CURSOR_"`
	Value       string       `help:"Initial value (can also be passed via stdin)" default:""`
	CharLimit   int          `help:"Maximum value length (0 for no limit)" default:"400"`
	Width       int          `help:"Input width" default:"40" env:"GUM_INPUT_WIDTH"`
	Password    bool         `help:"Mask input characters" default:"false"`
	Header      string       `help:"Header value" default:"" env:"GUM_INPUT_HEADER"`
	HeaderStyle style.Styles `embed:"" prefix:"header." set:"defaultForeground=240" envprefix:"GUM_INPUT_HEADER_"`
	timeout.CmdOptions
}
