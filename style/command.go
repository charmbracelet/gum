// Package style provides a shell script interface for Lip Gloss.
// https://github.com/charmbracelet/lipgloss
//
// It allows you to use Lip Gloss to style text without needing to use Go. All
// of the styling options are available as flags.
package style

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for the Lip Gloss styling.
// https://github.com/charmbracelet/lipgloss
func (o Options) Run() error {
	var text string
	if len(o.Text) > 0 {
		text = strings.Join(o.Text, "\n")
	} else {
		text, _ = stdin.Read()
		if text == "" {
			return errors.New("no input provided, see `gum style --help`")
		}
	}
	fmt.Println(o.Style.ToLipgloss().Render(text))
	return nil
}
