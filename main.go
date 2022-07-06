package main

import (
	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/sodapop/internal/stdin"
	"github.com/muesli/termenv"
)

func main() {
	lipgloss.SetColorProfile(termenv.ANSI256)
	pop := &Pop{}
	ctx := kong.Parse(pop,
		kong.Name("pop"),
		kong.Description("Make your shell pop."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}))
	switch ctx.Command() {
	case "input":
		pop.Input.Run()
	case "search":
		pop.Search.Run()
	case "spin <command>":
		pop.Spin.Run()
	case "style":
		pop.Style.Text, _ = stdin.Read()
		pop.Style.Run()
	case "style <text>":
		pop.Style.Run()
	case "layout":
	}
}
