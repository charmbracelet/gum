package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/gum/internal/log"
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
	case "cat":
		input, err := stdin.Read()
		if err != nil || input == "" {
			log.Error("No input provided.")
			log.Info("Please provide input with stdin or pass a file name as an argument.")
			return
		}
		gum.Cat.Text = input
		gum.Cat.Run()
	case "cat <file>":
		bytes, err := os.ReadFile(gum.Cat.File)
		if err != nil {
			log.Error(err.Error())
			return
		}
		gum.Cat.Text = string(bytes)
		gum.Cat.Run()
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
