package spin

import "github.com/charmbracelet/gum/style"

// Options is the customization options for the spin command.
// nolint:staticcheck
type Options struct {
	Command []string `arg:"" help:"Command to run"`

	ShowOutput   bool         `help:"Show output of command" default:"false"`
	Spinner      string       `help:"Spinner type" short:"s" type:"spinner" enum:"line,dot,minidot,jump,pulse,points,globe,moon,monkey,meter,hamburger" default:"dot"`
	SpinnerStyle style.Styles `embed:"" prefix:"spinner." set:"defaultForeground=212"`
	Title        string       `help:"Text to display to user while spinning" default:"Loading..."`
	TitleStyle   style.Styles `embed:"" prefix:"title."`
}
