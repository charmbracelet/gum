package write

import "github.com/charmbracelet/gum/style"

// Options are the customization options for the textarea.
// nolint:staticcheck
type Options struct {
	Width           int    `help:"Text area width" default:"50"`
	Height          int    `help:"Text area height" default:"5"`
	Placeholder     string `help:"Placeholder value" default:"Write something..."`
	Prompt          string `help:"Prompt to display" default:"┃ "`
	ShowCursorLine  bool   `help:"Show cursor line" default:"false"`
	ShowLineNumbers bool   `help:"Show line numbers" default:"false"`
	Value           string `help:"Initial value (can be passed via stdin)" default:""`
	CharLimit       int    `help:"Maximum value length (0 for no limit)" default:"400"`

	BaseStyle             style.Styles `embed:"" prefix:"base."`
	CursorLineNumberStyle style.Styles `embed:"" prefix:"cursor-line-number." set:"defaultForeground=7"`
	CursorLineStyle       style.Styles `embed:"" prefix:"cursor-line."`
	CursorStyle           style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212"`
	EndOfBufferStyle      style.Styles `embed:"" prefix:"end-of-buffer." set:"defaultForeground=0"`
	LineNumberStyle       style.Styles `embed:"" prefix:"line-number." set:"defaultForeground=7"`
	PlaceholderStyle      style.Styles `embed:"" prefix:"placeholder." set:"defaultForeground=240"`
	PromptStyle           style.Styles `embed:"" prefix:"prompt." set:"defaultForeground=7"`
}
