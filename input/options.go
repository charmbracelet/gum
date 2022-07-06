package input

// Options are the customization options for the input.
type Options struct {
	Placeholder string `help:"Placeholder value" default:"Type something..."`
	Prompt      string `help:"Prompt to display" default:"> "`
	Width       int    `help:"Input width" default:"20"`
}
