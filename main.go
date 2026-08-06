// Package main is Gum: a tool for glamorous shell scripts.
package main

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alecthomas/kong"
	"github.com/charmbracelet/gum/internal/exit"
)

const shaLen = 7

const defaultBool = "false"

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
			Compact:             true,
			Summary:             false,
			NoExpandSubcommands: true,
		}),
		kong.Vars{
			"version":                 version,
			"versionNumber":           Version,
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
			"defaultUnderline":        defaultBool,
			"defaultBold":             defaultBool,
			"defaultFaint":            defaultBool,
			"defaultItalic":           defaultBool,
			"defaultStrikethrough":    defaultBool,
		},
	)
	if err := ctx.Run(); err != nil {
		var ex exit.ErrExit
		if errors.As(err, &ex) {
			os.Exit(int(ex))
		}
		if errors.Is(err, tea.ErrInterrupted) {
			os.Exit(exit.StatusAborted)
		}
		if errors.Is(err, tea.ErrProgramKilled) {
			fmt.Fprintln(os.Stderr, "timed out")
			os.Exit(exit.StatusTimeout)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
