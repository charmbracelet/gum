package main

import (
	"testing"

	"github.com/alecthomas/kong"
)

func TestSpinCommandParsesSpinnerAndCommandArgs(t *testing.T) {
	gum := &Gum{}
	parser, err := kong.New(gum, kong.Vars{
		"version":                 "gum version test",
		"versionNumber":           "test",
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
	})
	if err != nil {
		t.Fatalf("failed to build parser: %v", err)
	}

	_, err = parser.Parse([]string{"spin", "--spinner", "monkey", "sleep", "5"})
	if err != nil {
		t.Fatalf("failed to parse spin command: %v", err)
	}

	if got, want := gum.Spin.Spinner, "monkey"; got != want {
		t.Fatalf("spinner = %q, want %q", got, want)
	}

	if got, want := gum.Spin.Command, []string{"sleep", "5"}; len(got) != len(want) {
		t.Fatalf("command length = %d, want %d (%v)", len(got), len(want), got)
	} else {
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("command[%d] = %q, want %q (full command: %v)", i, got[i], want[i], got)
			}
		}
	}
}
