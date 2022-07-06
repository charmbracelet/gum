package write

// Options are the customization options for the textarea.
type Options struct {
	ShowCursorLine  bool   `help:"Show cursor line" default:"true"`
	ShowLineNumbers bool   `help:"Show line numbers" default:"true"`
	Placeholder     string `help:"Placeholder value" default:"Write something..."`
	Prompt          string `help:"Prompt to display" default:"┃ "`
	Width           int    `help:"Text area width" default:"50"`
	Height          int    `help:"Text area height" default:"5"`
}
