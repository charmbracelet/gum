package style

import "strings"

// namedANSIColors maps the 16 standard ANSI terminal color names to the
// numeric ANSI color codes lipgloss.Color already understands (see
// https://github.com/charmbracelet/lipgloss/blob/master/color.go). Using
// the numeric codes (instead of hardcoded hex values) preserves lipgloss's
// automatic color degradation across ANSI, ANSI256, and TrueColor profiles.
var namedANSIColors = map[string]string{
	"black":         "0",
	"red":           "1",
	"green":         "2",
	"yellow":        "3",
	"blue":          "4",
	"magenta":       "5",
	"cyan":          "6",
	"white":         "7",
	"brightblack":   "8",
	"brightred":     "9",
	"brightgreen":   "10",
	"brightyellow":  "11",
	"brightblue":    "12",
	"brightmagenta": "13",
	"brightcyan":    "14",
	"brightwhite":   "15",
}

// resolveColor maps a common named color (e.g. "red") to the ANSI color
// code lipgloss.Color expects, so that flags like --foreground/--background
// accept human readable names in addition to hex codes and ANSI numbers.
//
// See https://github.com/charmbracelet/gum/issues/980: previously, an
// unrecognized name like "red" was passed straight to lipgloss.Color,
// which silently rendered no color at all.
//
// Values that are not a recognized name (hex codes, ANSI numbers, or truly
// unsupported strings) are returned unchanged so lipgloss/termenv continue
// to handle them exactly as before.
func resolveColor(value string) string {
	if code, ok := namedANSIColors[strings.ToLower(strings.TrimSpace(value))]; ok {
		return code
	}
	return value
}
