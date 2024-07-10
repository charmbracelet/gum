package ansi

import "github.com/charmbracelet/x/ansi"

// Strip strips a string of any of it's ansi sequences.
func Strip(text string) string {
	return ansi.Strip(text)
}
