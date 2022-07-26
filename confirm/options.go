package confirm

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the confirm command.
type Options struct {
	Affirmative     string       `help:"The title of the affirmative action" default:"Yes"`
	Negative        string       `help:"The title of the negative action" default:"No"`
	Prompt          string       `arg:"" help:"Prompt to display." default:"Are you sure?"`
	Vertical        bool         `help:"Whether to display the options vertically" default:"false"`
	PromptStyle     style.Styles `embed:"" prefix:"prompt." help:"The style of the prompt" set:"defaultMargin=1 0 0 0"`
	UnselectedStyle style.Styles `embed:"" prefix:"unselected." help:"The style of the unselected action" set:"defaultBackground=237" set:"defaultPadding=0 3" set:"defaultMargin=1 1"`
	SelectedStyle   style.Styles `embed:"" prefix:"selected." help:"The style of the selected action" set:"defaultBackground=240" set:"defaultForeground=212" set:"defaultPadding=0 3" set:"defaultMargin=1 1"`
}
