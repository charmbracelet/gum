package filter

// Options is the customization options for the filter command.
type Options struct {
	HighlightColor string `help:"Color for highlighted matches" default:"212"`
	Indicator      string `help:"Character for selection" default:"•"`
	IndicatorColor string `help:"Color for indicator" default:"212"`
	Placeholder    string `help:"Placeholder value" default:"Search..."`
	Prompt         string `help:"Prompt to display" default:"> "`
	PromptColor    string `help:"Color for prompt" default:"240"`
	Width          int    `help:"Input width" default:"20"`
}
