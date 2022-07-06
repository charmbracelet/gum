package search

// Options is the customization options for the search.
type Options struct {
	AccentColor string `help:"Accent color for prompt, indicator, and matches" default:"#FF06B7"`
	Indicator   string `help:"Character for selection" default:"•"`
	Placeholder string `help:"Placeholder value" default:"Search..."`
	Prompt      string `help:"Prompt to display" default:"> "`
	Width       int    `help:"Input width" default:"20"`
}
