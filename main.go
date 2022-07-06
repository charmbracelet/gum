package main

import (
	"github.com/alecthomas/kong"
)

// Seashell is the top-level command.
var Seashell struct {
	Input struct {
		Prompt      string `help:"Prompt to display" default:"> "`
		Placeholder string `help:"Placeholder value" default:"..."`
		Width       int    `help:"Input width" default:"20"`
	} `cmd:"" help:"Prompt for input."`

	Search struct {
		Prompt      string `help:"Prompt to display" default:"> "`
		Placeholder string `help:"Placeholder value" default:"..."`
		Width       int    `help:"Input width" default:"20"`
		AccentColor string `help:"Accent color for prompt, indicator, and matches." default:"#FF06B7"`
		Indicator   string `help:"Character to use to indicate the selected matches" default:"•"`
	} `cmd:"" help:"Fuzzy search options."`

	Spin struct {
		Spinner string `help:"Type of spinner to use" enum:"line,dot,minidot,jump,pulse,points,globe,moon,monkey,meter,hamburger" default:"dots"`
		Title   string `help:"Action being performed" default:"Loading..."`
		Color   string `help:"Color to use for the spinner" default:"#FF06B7"`
	} `cmd:"" help:"Show spinner while executing a command."`

	Style struct {
		Background       string `help:"The background color to apply"`
		Foreground       string `help:"The foreground color to apply"`
		BorderBackground string `help:"The border background color to apply"`
		BorderForeground string `help:"The border foreground color to apply"`
		Align            string `help:"The text alignment (left, center, right, bottom, middle, top)"`
		Border           string `help:"The border style to apply (none, hidden, normal, rounded, thick, double)"`
		Height           int    `help:"The height the output should take up"`
		Width            int    `help:"The width the output should take up"`
		Margin           string `help:"Margin to apply around the text."`
		Padding          string `help:"Padding to apply around the text."`
		Bold             bool   `help:"Whether to apply bold formatting"`
		Faint            bool   `help:"Whether to apply faint formatting"`
		Italic           bool   `help:"Whether to apply italic formatting"`
		Strikethrough    bool   `help:"Whether to apply strikethrough formatting"`
	} `cmd:"" help:"Style some text."`

	Layout struct {
	} `cmd:"" help:"Layout some text."`
}

func main() {
	ctx := kong.Parse(&Seashell)
	switch ctx.Command() {
	case "input":
	case "search":
	case "spin":
	case "style":
		fmt.Println(Seashell)
	case "layout":
	default:
	}
}
