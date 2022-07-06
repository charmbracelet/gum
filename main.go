package main

import (
	"fmt"

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
		fmt.Println(pop.Input)
	case "search":
		fmt.Println(pop.Search)
	case "spin":
		fmt.Println(pop.Spin)
	case "style":
		fmt.Println(pop.Style)
	case "layout":
	default:
	}
}
