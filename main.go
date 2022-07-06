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
		pop.InputCmd()
	case "search":
	case "spin":
	case "style":
	case "layout":
	default:
	}
}
