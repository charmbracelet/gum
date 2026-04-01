package style

import "strings"

// namedColors maps common color names to their ANSI 256 equivalents.
var namedColors = map[string]string{
	"black":   "0",
	"red":     "1",
	"green":   "2",
	"yellow":  "3",
	"blue":    "4",
	"magenta": "5",
	"cyan":    "6",
	"white":   "7",

	// Bright variants
	"brightblack":   "8",
	"brightred":     "9",
	"brightgreen":   "10",
	"brightyellow":  "11",
	"brightblue":    "12",
	"brightmagenta": "13",
	"brightcyan":    "14",
	"brightwhite":   "15",

	// Common aliases
	"gray":   "8",
	"grey":   "8",
	"orange": "208",
	"pink":   "212",
	"purple": "99",
}

// ResolveColor converts a named color to its ANSI value, or returns the
// input unchanged if it's already a valid color spec (hex, ANSI number).
func ResolveColor(color string) string {
	if color == "" {
		return color
	}
	if v, ok := namedColors[strings.ToLower(color)]; ok {
		return v
	}
	return color
}
