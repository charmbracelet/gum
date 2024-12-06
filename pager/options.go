package pager

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options are the options for the pager.
type Options struct {
	//nolint:staticcheck
	Style               style.Styles  `embed:"" help:"Style the pager" set:"defaultBorder=rounded" set:"defaultPadding=0 1" set:"defaultBorderForeground=212" envprefix:"GUM_PAGER_"`
	Content             string        `arg:"" optional:"" help:"Display content to scroll"`
	ShowLineNumbers     bool          `help:"Show line numbers" default:"true"`
	LineNumberStyle     style.Styles  `embed:"" prefix:"line-number." help:"Style the line numbers" set:"defaultForeground=237" envprefix:"GUM_PAGER_LINE_NUMBER_"`
	SoftWrap            bool          `help:"Soft wrap lines" default:"false"`
	MatchStyle          style.Styles  `embed:"" prefix:"match." help:"Style the matched text" set:"defaultForeground=212" set:"defaultBold=true" envprefix:"GUM_PAGER_MATCH_"`                                                      //nolint:staticcheck
	MatchHighlightStyle style.Styles  `embed:"" prefix:"match-highlight." help:"Style the matched highlight text" set:"defaultForeground=235" set:"defaultBackground=225" set:"defaultBold=true" envprefix:"GUM_PAGER_MATCH_HIGH_"` //nolint:staticcheck
	Timeout             time.Duration `help:"Timeout until command exits" default:"0s" env:"GUM_PAGER_TIMEOUT"`

	// Deprecated: this has no effect anymore.
	HelpStyle style.Styles `embed:"" prefix:"help." help:"Style the help text" set:"defaultForeground=241" envprefix:"GUM_PAGER_HELP_" hidden:""`
}
