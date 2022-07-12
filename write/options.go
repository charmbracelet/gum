package write

import "github.com/charmbracelet/gum/style"

// Options are the customization options for the textarea.
type Options struct {
	Width           int    `help:"Text area width" default:"50"`
	Height          int    `help:"Text area height" default:"5"`
	Placeholder     string `help:"Placeholder value" default:"Write something..."`
	Prompt          string `help:"Prompt to display" default:"┃ "`
	ShowCursorLine  bool   `help:"Show cursor line" default:"false"`
	ShowLineNumbers bool   `help:"Show line numbers" default:"false"`
	Value           string `help:"Initial value (can be passed via stdin)" default:""`

	BaseStyle             style.Styles `embed:"" prefix:"base." set:"name=base"`
	CursorLineNumberStyle style.Styles `embed:"" prefix:"cursor-line-number." set:"defaultForeground=7" set:"name=cursor line number"`
	CursorLineStyle       style.Styles `embed:"" prefix:"cursor-line." set:"name=cursor line"`
	CursorStyle           style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" set:"name=cursor"`
	EndOfBufferStyle      style.Styles `embed:"" prefix:"end-of-buffer." set:"defaultForeground=0" set:"name=end of buffer"`
	LineNumberStyle       style.Styles `embed:"" prefix:"line-number." set:"defaultForeground=7" set:"name=line number"`
	PlaceholderStyle      style.Styles `embed:"" prefix:"placeholder." set:"defaultForeground=240" set:"name=placeholder"`
	PromptStyle           style.Styles `embed:"" prefix:"prompt." set:"defaultForeground=7" set:"name=prompt"`
}
