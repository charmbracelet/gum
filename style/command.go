// Package style provides a shell script interface for Lip Gloss.
// https://github.com/charmbracelet/lipgloss
//
// It allows you to use Lip Gloss to style text without needing to use Go. All
// of the styling options are available as flags.
//
// Let's make some text glamorous using bash:
//
//   $ gum style \
//  	--foreground 212 --border double --align center \
//  	--width 50 --margin 2 --padding "2 4" \
//  	"Bubble Gum (1¢)" "So sweet and so fresh\!"
//
//
//    ╔══════════════════════════════════════════════════╗
//    ║                                                  ║
//    ║                                                  ║
//    ║                 Bubble Gum (1¢)                  ║
//    ║              So sweet and so fresh!              ║
//    ║                                                  ║
//    ║                                                  ║
//    ╚══════════════════════════════════════════════════╝
//
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

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	return HideFlags(ctx)
}

// HideFlags hides the flags from the usage output. This is used in conjunction
// with BeforeReset hook.
func HideFlags(ctx *kong.Context) error {
	n := ctx.Selected()
	if n == nil {
		return nil
	}
	for _, f := range n.Flags {
		if g := f.Group; g != nil && g.Key == groupName {
			if !strings.HasSuffix(f.Name, ".foreground") {
				f.Hidden = true
			}
		}
	}
	return nil
}
