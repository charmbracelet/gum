package filter

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the filter command.
type Options struct {
	Indicator      string       `help:"Character for selection" default:"•"`
	IndicatorStyle style.Styles `embed:"" prefix:"indicator." set:"defaultForeground=212" set:"name=indicator"`
	TextStyle      style.Styles `embed:"" prefix:"text."`
	MatchStyle     style.Styles `embed:"" prefix:"match." set:"defaultForeground=212" set:"name=matched text"`
	Placeholder    string       `help:"Placeholder value" default:"Filter..."`
	Prompt         string       `help:"Prompt to display" default:"> "`
	PromptStyle    style.Styles `embed:"" prefix:"prompt." set:"defaultForeground=240" set:"name=prompt"`
	Width          int          `help:"Input width" default:"20"`
}
