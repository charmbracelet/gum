package spin

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options is the customization options for the spin command.
type Options struct {
	Command []string `arg:"" help:"Command to run"`

	ShowOutput    bool         `help:"Show output of command" default:"false" env:"GUM_SPIN_SHOW_OUTPUT"`
	OutputOnError bool         `help:"Show output of command if it fails" default:"false" env:"GUM_SPIN_OUTPUT_ON_ERROR"`
	Spinner       string       `help:"Spinner type" short:"s" type:"spinner" enum:"line,dot,minidot,jump,pulse,points,globe,moon,monkey,meter,hamburger" default:"dot" env:"GUM_SPIN_SPINNER"`
	SpinnerStyle  style.Styles `embed:"" prefix:"spinner." set:"defaultForeground=212" envprefix:"GUM_SPIN_SPINNER_"`
	Title         string       `help:"Text to display to user while spinning" default:"Loading..." env:"GUM_SPIN_TITLE"`
	TitleStyle    style.Styles `embed:"" prefix:"title." envprefix:"GUM_SPIN_TITLE_"`
	Align         string       `help:"Alignment of spinner with regard to the title" short:"a" type:"align" enum:"left,right" default:"left" env:"GUM_SPIN_ALIGN"`
  Timeout      time.Duration `help:"Timeout until spin command aborts" default:"0" env:"GUM_SPIN_TIMEOUT"`
}
