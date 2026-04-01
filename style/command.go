// Package style provides a shell script interface for Lip Gloss.
// https://github.com/charmbracelet/lipgloss
//
// It allows you to use Lip Gloss to style text without needing to use Go. All
// of the styling options are available as flags.
package style

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/x/term"
)

// Run provides a shell script interface for the Lip Gloss styling.
// https://github.com/charmbracelet/lipgloss
func (o Options) Run() error {
	var text string
	if len(o.Text) > 0 {
		text = strings.Join(o.Text, "\n")
	} else {
		text, _ = stdin.Read(stdin.StripANSI(o.StripANSI))
		if text == "" {
			return errors.New("no input provided, see `gum style --help`")
		}
	}
	if o.Trim {
		var lines []string
		for _, line := range strings.Split(text, "\n") {
			lines = append(lines, strings.TrimSpace(line))
		}
		text = strings.Join(lines, "\n")
	}
	s := o.Style.ToLipgloss()
	if o.AutoWrap && o.Style.Width == 0 {
		if w, _, err := term.GetSize(os.Stdout.Fd()); err == nil && w > 0 {
			s = s.Width(w)
		}
	}
	fmt.Println(s.Render(text))
	return nil
}
