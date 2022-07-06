package spin

import "fmt"

// Run provides a shell script interface for the spinner bubble.
// https://github.com/charmbracelet/bubbles/spinner
func (o Options) Run() {
	fmt.Println(o)
}
