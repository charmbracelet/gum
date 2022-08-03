package write

import "github.com/charmbracelet/gum/style"

// Options are the customization options for the textarea.
// nolint:staticcheck
type Options struct {
	Width           int    `help:"Text area width" default:"50" env:"GUM_WRITE_WIDTH"`
	Height          int    `help:"Text area height" default:"5" env:"GUM_WRITE_HEIGHT"`
	Placeholder     string `help:"Placeholder value" default:"Write something..." env:"GUM_WRITE_PLACEHOLDER"`
	Prompt          string `help:"Prompt to display" default:"┃ " env:"GUM_WRITE_PROMPT"`
	ShowCursorLine  bool   `help:"Show cursor line" default:"false" env:"GUM_WRITE_SHOW_CURSOR_LINE"`
	ShowLineNumbers bool   `help:"Show line numbers" default:"false" env:"GUM_WRITE_SHOW_LINE_NUMBERS"`
	Value           string `help:"Initial value (can be passed via stdin)" default:"" env:"GUM_WRITE_VALUE"`
	CharLimit       int    `help:"Maximum value length (0 for no limit)" default:"400"`

	BaseStyle             style.Styles `embed:"" prefix:"base." envprefix:"GUM_WRITE_BASE_"`
	CursorLineNumberStyle style.Styles `embed:"" prefix:"cursor-line-number." set:"defaultForeground=7" envprefix:"GUM_WRITE_CURSOR_LINE_NUMBER_"`
	CursorLineStyle       style.Styles `embed:"" prefix:"cursor-line." envprefix:"GUM_WRITE_CURSOR_LINE_"`
	CursorStyle           style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" envprefix:"GUM_WRITE_CURSOR_"`
	EndOfBufferStyle      style.Styles `embed:"" prefix:"end-of-buffer." set:"defaultForeground=0" envprefix:"GUM_WRITE_END_OF_BUFFER_"`
	LineNumberStyle       style.Styles `embed:"" prefix:"line-number." set:"defaultForeground=7" envprefix:"GUM_WRITE_LINE_NUMBER_"`
	PlaceholderStyle      style.Styles `embed:"" prefix:"placeholder." set:"defaultForeground=240" envprefix:"GUM_WRITE_PLACEHOLDER_"`
	PromptStyle           style.Styles `embed:"" prefix:"prompt." set:"defaultForeground=7" envprefix:"GUM_WRITE_PROMPT_"`
}
