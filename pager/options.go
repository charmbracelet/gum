package pager

import (
	"github.com/charmbracelet/gum/style"
	"time"
)

// Options are the options for the pager.
type Options struct {
	//nolint:staticcheck
	Style           style.Styles  `embed:"" help:"Style the pager" set:"defaultBorder=rounded" set:"defaultPadding=0 1" set:"defaultBorderForeground=212" envprefix:"GUM_PAGER_"`
	HelpStyle       style.Styles  `embed:"" prefix:"help." help:"Style the help text" set:"defaultForeground=241" envprefix:"GUM_PAGER_HELP_"`
	Content         string        `arg:"" optional:"" help:"Display content to scroll"`
	ShowLineNumbers bool          `help:"Show line numbers" default:"true"`
	LineNumberStyle style.Styles  `embed:"" prefix:"line-number." help:"Style the line numbers" set:"defaultForeground=237" envprefix:"GUM_PAGER_LINE_NUMBER_"`
	SoftWrap        bool          `help:"Soft wrap lines" default:"false"`
	Timeout         time.Duration `help:"Timeout until command exits" default:"0" env:"GUM_CONFIRM_TIMEOUT"`
}
