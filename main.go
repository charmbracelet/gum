package main

import (
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
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}))
	switch ctx.Command() {
	case "input":
		gum.Input.Run()
	case "write":
		gum.Write.Run()
	case "filter":
		gum.Filter.Run()
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
