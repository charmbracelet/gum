package ansi

import "regexp"

var ansiEscape = regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)

func Strip(text string) string {
	return ansiEscape.ReplaceAllString(text, "")
}
