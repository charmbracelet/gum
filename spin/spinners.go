package spin

import "github.com/charmbracelet/bubbles/spinner"

var spinnerMap = map[string]spinner.Spinner{
	"line":      spinner.Line,
	"dot":       spinner.Dot,
	"minidot":   spinner.MiniDot,
	"jump":      spinner.Jump,
	"pulse":     spinner.Pulse,
	"points":    spinner.Points,
	"globe":     spinner.Globe,
	"moon":      spinner.Moon,
	"monkey":    spinner.Monkey,
	"meter":     spinner.Meter,
	"hamburger": spinner.Hamburger,
}
