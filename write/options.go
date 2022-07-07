package write

// Options are the customization options for the textarea.
type Options struct {
	CursorColor     string `help:"Color for cursor" default:"212"`
	Height          int    `help:"Text area height" default:"5"`
	Placeholder     string `help:"Placeholder value" default:"Write something..."`
	Prompt          string `help:"Prompt to display" default:"┃ "`
	PromptColor     string `help:"Color for prompt" default:"238"`
	ShowCursorLine  bool   `help:"Show cursor line" default:"false"`
	ShowLineNumbers bool   `help:"Show line numbers" default:"false"`
	Width           int    `help:"Text area width" default:"50"`
}
