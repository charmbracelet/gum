package write

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options are the customization options for the textarea.
type Options struct {
	Width           int           `help:"Text area width (0 for terminal width)" default:"0" env:"GUM_WRITE_WIDTH"`
	Height          int           `help:"Text area height" default:"5" env:"GUM_WRITE_HEIGHT"`
	Header          string        `help:"Header value" default:"" env:"GUM_WRITE_HEADER"`
	Placeholder     string        `help:"Placeholder value" default:"Write something..." env:"GUM_WRITE_PLACEHOLDER"`
	Prompt          string        `help:"Prompt to display" default:"┃ " env:"GUM_WRITE_PROMPT"`
	ShowCursorLine  bool          `help:"Show cursor line" default:"false" env:"GUM_WRITE_SHOW_CURSOR_LINE"`
	ShowLineNumbers bool          `help:"Show line numbers" default:"false" env:"GUM_WRITE_SHOW_LINE_NUMBERS"`
	Value           string        `help:"Initial value (can be passed via stdin)" default:"" env:"GUM_WRITE_VALUE"`
	CharLimit       int           `help:"Maximum value length (0 for no limit)" default:"0"`
	MaxLines        int           `help:"Maximum number of lines (0 for no limit)" default:"0"`
	ShowHelp        bool          `help:"Show help key binds" negatable:"" default:"true" env:"GUM_WRITE_SHOW_HELP"`
	CursorMode      string        `prefix:"cursor." name:"mode" help:"Cursor mode" default:"blink" enum:"blink,hide,static" env:"GUM_WRITE_CURSOR_MODE"`
	Timeout         time.Duration `help:"Timeout until choose returns selected element" default:"0s" env:"GUM_WRITE_TIMEOUT"`
	StripANSI       bool          `help:"Strip ANSI sequences when reading from STDIN" default:"true" negatable:"" env:"GUM_WRITE_STRIP_ANSI"`

	BaseStyle             style.Styles `embed:"" prefix:"base." envprefix:"GUM_WRITE_BASE_"`
	CursorLineNumberStyle style.Styles `embed:"" prefix:"cursor-line-number." set:"defaultForeground=7" envprefix:"GUM_WRITE_CURSOR_LINE_NUMBER_"`
	CursorLineStyle       style.Styles `embed:"" prefix:"cursor-line." envprefix:"GUM_WRITE_CURSOR_LINE_"`
	CursorStyle           style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" envprefix:"GUM_WRITE_CURSOR_"`
	EndOfBufferStyle      style.Styles `embed:"" prefix:"end-of-buffer." set:"defaultForeground=0" envprefix:"GUM_WRITE_END_OF_BUFFER_"`
	LineNumberStyle       style.Styles `embed:"" prefix:"line-number." set:"defaultForeground=7" envprefix:"GUM_WRITE_LINE_NUMBER_"`
	HeaderStyle           style.Styles `embed:"" prefix:"header." set:"defaultForeground=240" envprefix:"GUM_WRITE_HEADER_"`
	PlaceholderStyle      style.Styles `embed:"" prefix:"placeholder." set:"defaultForeground=240" envprefix:"GUM_WRITE_PLACEHOLDER_"`
	PromptStyle           style.Styles `embed:"" prefix:"prompt." set:"defaultForeground=7" envprefix:"GUM_WRITE_PROMPT_"`
	Padding               string       `help:"Padding" default:"${defaultPadding}" group:"Style Flags" env:"GUM_WRITE_PADDING"`
}
