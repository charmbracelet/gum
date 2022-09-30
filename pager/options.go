package pager

import "github.com/charmbracelet/gum/style"

// Options are the options for the pager.
type Options struct {
	Style     style.Styles `embed:"" help:"Style the pager" envprefix:"GUM_PAGER_"`
	HelpStyle style.Styles `embed:"" prefix:"help." help:"Style the help text" set:"defaultForeground=241" envprefix:"GUM_PAGER_HELP_"`
	Content   string       `arg:"" optional:"" help:"Display content to scroll"`
}
