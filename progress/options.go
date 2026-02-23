package progress

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options is the customization options for the progress command.
type Options struct {
	Title string `help:"Title to display above the progress bar" default:"" env:"GUM_PROGRESS_TITLE"`

	// Gradient colors
	GradientFrom string `help:"Gradient start color" default:"#5A56E0" env:"GUM_PROGRESS_GRADIENT_FROM"`
	GradientTo   string `help:"Gradient end color" default:"#EE6FF8" env:"GUM_PROGRESS_GRADIENT_TO"`
	Solid        string `help:"Solid fill color (overrides gradient)" default:"" env:"GUM_PROGRESS_SOLID"`

	Width          int           `help:"Width of the progress bar" default:"40" env:"GUM_PROGRESS_WIDTH"`
	ShowPercentage bool          `help:"Show percentage next to progress bar" default:"true" negatable:"" env:"GUM_PROGRESS_SHOW_PERCENTAGE"`
	Timeout        time.Duration `help:"Timeout until progress returns" default:"0s" env:"GUM_PROGRESS_TIMEOUT"`
	Padding        string        `help:"Padding" default:"${defaultPadding}" group:"Style Flags" env:"GUM_PROGRESS_PADDING"`

	TitleStyle style.Styles `embed:"" prefix:"title." envprefix:"GUM_PROGRESS_TITLE_"`
}
