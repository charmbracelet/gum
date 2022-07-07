package input

// Options are the customization options for the input.
type Options struct {
	CursorColor string `help:"Color of prompt" default:"212"`
	Placeholder string `help:"Placeholder value" default:"Type something..."`
	Prompt      string `help:"Prompt to display" default:"> "`
	PromptColor string `help:"Color of prompt" default:"7"`
	Width       int    `help:"Input width" default:"20"`
}
