package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

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
	lipgloss.SetColorProfile(termenv.ANSI256)

	if Version == "" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
			Version = info.Main.Version
		} else {
			Version = "unknown (built from source)"
		}
	}
	version := fmt.Sprintf("gum version %s", Version)
	if len(CommitSHA) >= 7 {
		version += " (" + CommitSHA[:7] + ")"
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
			"version":           version,
			"defaultBackground": "",
			"defaultForeground": "",
			"defaultMargin":     "0 0",
			"defaultPadding":    "0 0",
			"defaultUnderline":  "false",
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
