// Package format allows you to render formatted text from the command line.
//
// It supports the following types:
//
// 1. Markdown
// 2. Code
// 3. Emoji
// 4. Template
//
// For more information, see the format/README.md file.
package format

import (
	"fmt"
	"os"
	"strings"

	"charm.land/gum/v2/internal/stdin"
	"github.com/charmbracelet/colorprofile"
)

// Run runs the format command.
func (o Options) Run() error {
	var input, output string
	var err error
	if len(o.Template) > 0 {
		input = strings.Join(o.Template, "\n")
	} else {
		input, _ = stdin.Read(stdin.StripANSI(o.StripANSI))
	}

	switch o.Type {
	case "code":
		output, err = code(input, o.Language)
	case "emoji":
		output, err = emoji(input)
	case "template":
		output, err = template(input)
	default:
		output, err = markdown(input, o.Theme)
	}
	if err != nil {
		return err
	}

	// Preserve ANSI colors when piped (e.g. `gum format | less -R`).
	// Previous `tty.Writer()` used `colorprofile.Detect` which strips
	// colors when stdout is not a TTY (NoTTY), breaking `format | cat`
	// (v2 regression, see #1131). `format` is a non-interactive converter
	// and should emit colors even when piped, respecting NO_COLOR,
	// CLICOLOR_FORCE and TERM via `colorprofile.Env` (which ignores TTY).
	profile := colorprofile.Env(os.Environ())
	switch o.Color {
	case "always":
		// Explicit flag overrides TERM=dumb and NO_COLOR - force at least ANSI.
		if profile < colorprofile.ANSI {
			profile = colorprofile.ANSI
		}
	case "never":
		// Fully disable colors (and styles) - stronger than NO_COLOR which keeps bold.
		profile = colorprofile.NoTTY
	default: // auto: keep fixed behavior - preserve when piped, respect env.
		if profile == colorprofile.NoTTY {
			// Env returned NoTTY (no TERM or TERM=dumb without CLICOLOR_FORCE).
			// For `format | less -R` we want colors even when piped, matching v1.
			// Only force ANSI when TERM is not dumb; respect TERM=dumb which
			// explicitly requests no color unless CLICOLOR_FORCE (already handled
			// by Env returning >=ANSI).
			if term := os.Getenv("TERM"); term != "dumb" {
				profile = colorprofile.ANSI
			}
		}
	}
	w := &colorprofile.Writer{
		Forward: os.Stdout,
		Profile: profile,
	}
	if _, err := fmt.Fprint(w, output); err != nil {
		return fmt.Errorf("unable to write output: %w", err)
	}
	return nil
}
