package style

import "fmt"

// Run provides a shell script interface for the Lip Gloss styling.
// https://github.com/charmbracelet/lipgloss
func (o Options) Run() {
	fmt.Println(o)
}
