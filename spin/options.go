package spin

// Options is the customization options for the spin command.
type Options struct {
	Color   string `help:"Spinner color" default:"#FF06B7"`
	Display string `help:"Text to display to user while spinning" default:"Loading..."`
	Spinner string `help:"Spinner type" enum:"line,dot,minidot,jump,pulse,points,globe,moon,monkey,meter,hamburger" default:"dot"`
}
