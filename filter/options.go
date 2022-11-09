package filter

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the filter command.
type Options struct {
	Indicator             string       `help:"Character for selection" default:"•" env:"GUM_FILTER_INDICATOR"`
	IndicatorStyle        style.Styles `embed:"" prefix:"indicator." set:"defaultForeground=212" envprefix:"GUM_FILTER_INDICATOR_"`
	Limit                 int          `help:"Maximum number of options to pick" default:"1" group:"Selection"`
	NoLimit               bool         `help:"Pick unlimited number of options (ignores limit)" group:"Selection"`
	NoMatchNeeded         bool         `help:"Only returns if anything matched. Otherwise return Filter" group:"Selection"`
	SelectedPrefix        string       `help:"Character to indicate selected items (hidden if limit is 1)" default:" ◉ " env:"GUM_FILTER_SELECTED_PREFIX"`
	SelectedPrefixStyle   style.Styles `embed:"" prefix:"selected-indicator." set:"defaultForeground=212" envprefix:"GUM_FILTER_SELECTED_PREFIX_"`
	UnselectedPrefix      string       `help:"Character to indicate unselected items (hidden if limit is 1)" default:" ○ " env:"GUM_FILTER_UNSELECTED_PREFIX"`
	UnselectedPrefixStyle style.Styles `embed:"" prefix:"unselected-prefix." set:"defaultForeground=240" envprefix:"GUM_FILTER_UNSELECTED_PREFIX_"`
	TextStyle             style.Styles `embed:"" prefix:"text." envprefix:"GUM_FILTER_TEXT_"`
	MatchStyle            style.Styles `embed:"" prefix:"match." set:"defaultForeground=212" envprefix:"GUM_FILTER_MATCH_"`
	Placeholder           string       `help:"Placeholder value" default:"Filter..." env:"GUM_FILTER_PLACEHOLDER"`
	Prompt                string       `help:"Prompt to display" default:"> " env:"GUM_FILTER_PROMPT"`
	PromptStyle           style.Styles `embed:"" prefix:"prompt." set:"defaultForeground=240" envprefix:"GUM_FILTER_PROMPT_"`
	Width                 int          `help:"Input width" default:"20" env:"GUM_FILTER_WIDTH"`
	Height                int          `help:"Input height" default:"0" env:"GUM_FILTER_HEIGHT"`
	Value                 string       `help:"Initial filter value" default:"" env:"GUM_FILTER_VALUE"`
	Reverse               bool         `help:"Display from the bottom of the screen" env:"GUM_FILTER_REVERSE"`
	Fuzzy                 bool         `help:"Enable fuzzy matching" default:"true" env:"GUM_FILTER_FUZZY" negatable:""`
}
