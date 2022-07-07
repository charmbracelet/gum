package spin

// Options is the customization options for the spin command.
type Options struct {
	Command []string `arg:"" help:"Command to run"`

	Color   string `help:"Spinner color" default:"212"`
	Spinner string `help:"Spinner type" type:"spinner" enum:"line,dot,minidot,jump,pulse,points,globe,moon,monkey,meter,hamburger" default:"dot"`
	Title   string `help:"Text to display to user while spinning" default:"Loading..."`
}
