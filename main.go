package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var bubbleGumPink = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

func main() {
	lipgloss.SetColorProfile(termenv.ANSI256)

	gum := &Gum{}
	ctx := kong.Parse(
		gum,
		kong.Description(fmt.Sprintf("Tasty bubble %s for your shell.", bubbleGumPink.Render("gum"))),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}),
		kong.Vars{
			"defaultBackground": "",
			"defaultForeground": "",
			"defaultMargin":     "0 0",
			"defaultPadding":    "0 0",
			"defaultUnderline":  "false",
		},
	)
	ctx.Run()
}
