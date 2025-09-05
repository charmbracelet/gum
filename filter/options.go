package filter

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options is the customization options for the filter command.
type Options struct {
	Options []string `arg:"" optional:"" help:"Options to filter."`

	Indicator             string        `help:"Character for selection" default:"•" env:"GUM_FILTER_INDICATOR"`
	IndicatorStyle        style.Styles  `embed:"" prefix:"indicator." set:"defaultForeground=212" envprefix:"GUM_FILTER_INDICATOR_"`
	Limit                 int           `help:"Maximum number of options to pick" default:"1" group:"Selection"`
	NoLimit               bool          `help:"Pick unlimited number of options (ignores limit)" group:"Selection"`
	SelectIfOne           bool          `help:"Select the given option if there is only one" group:"Selection"`
	Selected              []string      `help:"Options that should start as selected (selects all if given *)" default:"" env:"GUM_FILTER_SELECTED"`
	ShowHelp              bool          `help:"Show help keybinds" default:"true" negatable:"" env:"GUM_FILTER_SHOW_HELP"`
	Strict                bool          `help:"Only returns if anything matched. Otherwise return Filter" negatable:"" default:"true" group:"Selection"`
	SelectedPrefix        string        `help:"Character to indicate selected items (hidden if limit is 1)" default:" ◉ " env:"GUM_FILTER_SELECTED_PREFIX"`
	SelectedPrefixStyle   style.Styles  `embed:"" prefix:"selected-indicator." set:"defaultForeground=212" envprefix:"GUM_FILTER_SELECTED_PREFIX_"`
	UnselectedPrefix      string        `help:"Character to indicate unselected items (hidden if limit is 1)" default:" ○ " env:"GUM_FILTER_UNSELECTED_PREFIX"`
	UnselectedPrefixStyle style.Styles  `embed:"" prefix:"unselected-prefix." set:"defaultForeground=240" envprefix:"GUM_FILTER_UNSELECTED_PREFIX_"`
	HeaderStyle           style.Styles  `embed:"" prefix:"header." set:"defaultForeground=99" envprefix:"GUM_FILTER_HEADER_"`
	Header                string        `help:"Header value" default:"" env:"GUM_FILTER_HEADER"`
	TextStyle             style.Styles  `embed:"" prefix:"text." envprefix:"GUM_FILTER_TEXT_"`
	CursorTextStyle       style.Styles  `embed:"" prefix:"cursor-text." envprefix:"GUM_FILTER_CURSOR_TEXT_"`
	MatchStyle            style.Styles  `embed:"" prefix:"match." set:"defaultForeground=212" envprefix:"GUM_FILTER_MATCH_"`
	Placeholder           string        `help:"Placeholder value" default:"Filter..." env:"GUM_FILTER_PLACEHOLDER"`
	Prompt                string        `help:"Prompt to display" default:"> " env:"GUM_FILTER_PROMPT"`
	PromptStyle           style.Styles  `embed:"" prefix:"prompt." set:"defaultForeground=240" envprefix:"GUM_FILTER_PROMPT_"`
	PlaceholderStyle      style.Styles  `embed:"" prefix:"placeholder." set:"defaultForeground=240" envprefix:"GUM_FILTER_PLACEHOLDER_"`
	Width                 int           `help:"Input width" default:"0" env:"GUM_FILTER_WIDTH"`
	Height                int           `help:"Input height" default:"0" env:"GUM_FILTER_HEIGHT"`
	Value                 string        `help:"Initial filter value" default:"" env:"GUM_FILTER_VALUE"`
	Reverse               bool          `help:"Display from the bottom of the screen" env:"GUM_FILTER_REVERSE"`
	Fuzzy                 bool          `help:"Enable fuzzy matching; otherwise match from start of word" default:"true" env:"GUM_FILTER_FUZZY" negatable:""`
	FuzzySort             bool          `help:"Sort fuzzy results by their scores" default:"true" env:"GUM_FILTER_FUZZY_SORT" negatable:""`
	Timeout               time.Duration `help:"Timeout until filter command aborts" default:"0s" env:"GUM_FILTER_TIMEOUT"`
	InputDelimiter        string        `help:"Option delimiter when reading from STDIN" default:"\n" env:"GUM_FILTER_INPUT_DELIMITER"`
	OutputDelimiter       string        `help:"Option delimiter when writing to STDOUT" default:"\n" env:"GUM_FILTER_OUTPUT_DELIMITER"`
	StripANSI             bool          `help:"Strip ANSI sequences when reading from STDIN" default:"true" negatable:"" env:"GUM_FILTER_STRIP_ANSI"`
	Padding               string        `help:"Padding" default:"${defaultPadding}" group:"Style Flags" env:"GUM_FILTER_PADDING"`

	// Deprecated: use [FuzzySort]. This will be removed at some point.
	Sort bool `help:"Sort fuzzy results by their scores" default:"true" env:"GUM_FILTER_FUZZY_SORT" negatable:"" hidden:""`
}
