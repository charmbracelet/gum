package progress

import "time"

// Options is the available options for the progress command.
type Options struct {
	ColorEnd   string        `help:"The end color of the progress bar gradient." default:"#EE6FF8"`
	ColorStart string        `help:"The start color of the progress bar gradient." default:"#5A56E0"`
	Damping    float64       `help:"Damping of the spring animation." default:"0.9"`
	Frequency  float64       `help:"Frequency of the spring animation." default:"10.0"`
	Increment  float64       `help:"The percentage to increment the progress bar per tick." default:"0.1"`
	Interval   time.Duration `help:"The interval of time to wait before incrementing." default:"100ms"`
}
