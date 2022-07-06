package main

import (
	"github.com/alecthomas/kong"
)

func main() {
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
	case "spin":
		pop.Spin.Run()
	case "style":
		pop.Style.Run()
	case "layout":
	default:
	}
}
