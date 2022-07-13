package style

import (
	"fmt"
	"strings"
)

// Run provides a shell script interface for the Lip Gloss styling.
// https://github.com/charmbracelet/lipgloss
func (o Options) Run() error {
	text := strings.Join(o.Text, "\n")
	fmt.Println(o.Style.ToLipgloss().Render(text))
	return nil
}
