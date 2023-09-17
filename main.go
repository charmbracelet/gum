package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/charmbracelet/gum/internal/exit"
)

const shaLen = 7

var (
	// Version contains the application version number. It's set via ldflags
	// when building.
	Version = ""

	// CommitSHA contains the SHA of the commit that this application was built
	// against. It's set via ldflags when building.
	CommitSHA = ""
)

var bubbleGumPink = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

func main() {
	lipgloss.SetColorProfile(termenv.NewOutput(os.Stderr).EnvColorProfile())

	if Version == "" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
			Version = info.Main.Version
		} else {
			Version = "unknown (built from source)"
		}
	}
	version := fmt.Sprintf("gum version %s", Version)
	if len(CommitSHA) >= shaLen {
		version += " (" + CommitSHA[:shaLen] + ")"
	}

	gum := &Gum{}
	ctx := kong.Parse(
		gum,
		kong.Description(fmt.Sprintf("A tool for %s shell scripts.", bubbleGumPink.Render("glamorous"))),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}),
		kong.Vars{
			"version":                 version,
			"defaultHeight":           "0",
			"defaultWidth":            "0",
			"defaultAlign":            "left",
			"defaultBorder":           "none",
			"defaultBorderForeground": "",
			"defaultBorderBackground": "",
			"defaultBackground":       "",
			"defaultForeground":       "",
			"defaultMargin":           "0 0",
			"defaultPadding":          "0 0",
			"defaultUnderline":        "false",
			"defaultBold":             "false",
			"defaultFaint":            "false",
			"defaultItalic":           "false",
			"defaultStrikethrough":    "false",
		},
	)
	if err := ctx.Run(); err != nil {
		if errors.Is(err, exit.ErrAborted) {
			os.Exit(exit.StatusAborted)
		}
		fmt.Println(err)
		os.Exit(1)
	}
}
