package main

import (
	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func main() {
	lipgloss.SetColorProfile(termenv.ANSI256)
	gum := &Gum{}
	ctx := kong.Parse(
		gum,
		kong.Name("gum"),
		kong.Description("Tasty Bubble Gum for your shell."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}),
		kong.Vars{"defaultBackground": "", "defaultForeground": ""},
	)
	ctx.Run()
}
