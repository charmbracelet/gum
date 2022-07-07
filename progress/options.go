package progress

import "time"

// Options is the available options for the progress command.
type Options struct {
	Damping   float64       `help:"Damping of the spring animation." default:"0.9"`
	Frequency float64       `help:"Frequency of the spring animation." default:"10.0"`
	Increment float64       `help:"The percentage to increment the progress bar per tick." default:"0.1"`
	Interval  time.Duration `help:"The interval of time to wait before incrementing." default:"100ms"`
}
