package confirm

import (
	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/gum/timeout"
)

// Options is the customization timeout for the confirm command.
type Options struct {
	Default     bool         `help:"Default confirmation action" default:"true"`
	Affirmative string       `help:"The title of the affirmative action" default:"Yes"`
	Negative    string       `help:"The title of the negative action" default:"No"`
	Prompt      string       `arg:"" help:"Prompt to display." default:"Are you sure?"`
	PromptStyle style.Styles `embed:"" prefix:"prompt." help:"The style of the prompt" set:"defaultMargin=1 0 0 0" envprefix:"GUM_CONFIRM_PROMPT_"`
	//nolint:staticcheck
	SelectedStyle style.Styles `embed:"" prefix:"selected." help:"The style of the selected action" set:"defaultBackground=212" set:"defaultForeground=230" set:"defaultPadding=0 3" set:"defaultMargin=1 1" envprefix:"GUM_CONFIRM_SELECTED_"`
	//nolint:staticcheck
	UnselectedStyle style.Styles `embed:"" prefix:"unselected." help:"The style of the unselected action" set:"defaultBackground=235" set:"defaultForeground=254" set:"defaultPadding=0 3" set:"defaultMargin=1 1" envprefix:"GUM_CONFIRM_UNSELECTED_"`
	timeout.Options              // including timeout command options [Timeout,...]
}
