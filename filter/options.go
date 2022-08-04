package filter

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the filter command.
type Options struct {
	Indicator      string       `help:"Character for selection" default:"•" env:"GUM_FILTER_INDICATOR"`
	IndicatorStyle style.Styles `embed:"" prefix:"indicator." set:"defaultForeground=212" envprefix:"GUM_FILTER_INDICATOR_"`
	TextStyle      style.Styles `embed:"" prefix:"text." envprefix:"GUM_FILTER_TEXT_"`
	MatchStyle     style.Styles `embed:"" prefix:"match." set:"defaultForeground=212" envprefix:"GUM_FILTER_MATCH_"`
	Placeholder    string       `help:"Placeholder value" default:"Filter..." env:"GUM_FILTER_PLACEHOLDER"`
	Prompt         string       `help:"Prompt to display" default:"> " env:"GUM_FILTER_PROMPT"`
	PromptStyle    style.Styles `embed:"" prefix:"prompt." set:"defaultForeground=240" envprefix:"GUM_FILTER_PROMPT_"`
	Width          int          `help:"Input width" default:"20" env:"GUM_FILTER_WIDTH"`
	Height         int          `help:"Input height" default:"0" env:"GUM_FILTER_HEIGHT"`
}
