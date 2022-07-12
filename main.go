package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func main() {
	lipgloss.SetColorProfile(termenv.ANSI256)
	gum := &Gum{}
	ctx := kong.Parse(gum,
		kong.Name("gum"),
		kong.Description("Tasty Bubble Gum for your shell."),
		kong.UsageOnError(),
		kong.Vars{
			"defaultBackground": "",
			"defaultForeground": "",
		},
	)
	switch ctx.Command() {
	case "input":
		gum.Input.Run()
	case "write":
		gum.Write.Run()
	case "filter":
		gum.Filter.Run()
	case "choose":
		input, _ := stdin.Read()
		gum.Choose.Options = strings.Split(input, "\n")
		gum.Choose.Run()
	case "choose <options>":
		gum.Choose.Run()
	case "spin <command>":
		gum.Spin.Run()
	case "progress":
		gum.Progress.Run()
	case "style":
		input, _ := stdin.Read()
		gum.Style.Text = []string{input}
		gum.Style.Run()
	case "style <text>":
		gum.Style.Run()
	case "join <text>":
		gum.Join.Run()
	}
}
