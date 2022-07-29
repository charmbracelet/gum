package input

import "github.com/charmbracelet/gum/style"

// Options are the customization options for the input.
type Options struct {
	Placeholder string       `help:"Placeholder value" default:"Type something..."`
	Prompt      string       `help:"Prompt to display" default:"> "`
	PromptStyle style.Styles `embed:"" prefix:"prompt." set:"name=prompt"`
	CursorStyle style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" set:"name=cursor"`
	Value       string       `help:"Initial value (can also be passed via stdin)" default:""`
	Width       int          `help:"Input width" default:"40"`
	Password    bool         `help:"Mask input characters" default:"false"`
}
