// Package style provides a shell script interface for Lip Gloss.
// https://github.com/charmbracelet/lipgloss
//
// It allows you to use Lip Gloss to style text without needing to use Go. All
// of the styling options are available as flags.
package style

import (
	"fmt"
	"strings"

	"github.com/alecthomas/kong"
)

// Run provides a shell script interface for the Lip Gloss styling.
// https://github.com/charmbracelet/lipgloss
func (o Options) Run() error {
	text := strings.Join(o.Text, "\n")
	fmt.Println(o.Style.ToLipgloss().Render(text))
	return nil
}

// HideFlags hides the flags from the usage output. This is used in conjunction
// with BeforeReset hook.
func HideFlags(ctx *kong.Context) {
	n := ctx.Selected()
	if n == nil {
		return
	}
	for _, f := range n.Flags {
		if g := f.Group; g != nil && g.Key == groupName {
			if !strings.HasSuffix(f.Name, ".foreground") {
				f.Hidden = true
			}
		}
	}
}
