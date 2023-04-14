package ansi

import "regexp"

var ansiEscape = regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)


// Strip strips a string of any of it's ansi sequences.

func Strip(text string) string {
	return ansiEscape.ReplaceAllString(text, "")
}
