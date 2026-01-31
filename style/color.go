package style

import (
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// ParseColor parses a color string and returns a normalized color value.
// It supports:
// - Hex colors: "#FF0000", "FF0000"
// - ANSI color codes: "1", "2", etc.
// - Named colors: "red", "green", "blue", etc. (SVG/CSS color names)
//
// If the input is empty or cannot be parsed, it returns the original value
// to maintain backward compatibility with lipgloss.Color behavior.
func ParseColor(s string) string {
	if s == "" {
		return s
	}

	// If it looks like a hex code or ANSI code, return as-is
	if strings.HasPrefix(s, "#") {
		return s
	}

	// If it's a pure number (ANSI code), return as-is
	if isNumeric(s) {
		return s
	}

	// Try to parse as a named color using go-colorful
	c, err := colorful.Hex(s)
	if err == nil {
		return c.Hex()
	}

	// Try named color (SVG/CSS color names)
	// go-colorful doesn't have built-in named color support,
	// so we'll add a lookup table for common colors
	if hex, ok := namedColors[strings.ToLower(s)]; ok {
		return hex
	}

	// Return original value if we can't parse it
	// This maintains backward compatibility
	return s
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

// namedColors maps CSS/SVG color names to hex values
var namedColors = map[string]string{
	// Basic colors
	"black":   "#000000",
	"white":   "#FFFFFF",
	"red":     "#FF0000",
	"green":   "#008000",
	"blue":    "#0000FF",
	"yellow":  "#FFFF00",
	"cyan":    "#00FFFF",
	"magenta": "#FF00FF",

	// Extended colors
	"gray":       "#808080",
	"grey":       "#808080",
	"silver":     "#C0C0C0",
	"maroon":     "#800000",
	"olive":      "#808000",
	"lime":       "#00FF00",
	"aqua":       "#00FFFF",
	"teal":       "#008080",
	"navy":       "#000080",
	"fuchsia":    "#FF00FF",
	"purple":     "#800080",
	"orange":     "#FFA500",
	"pink":       "#FFC0CB",
	"brown":      "#A52A2A",
	"coral":      "#FF7F50",
	"crimson":    "#DC143C",
	"gold":       "#FFD700",
	"indigo":     "#4B0082",
	"ivory":      "#FFFFF0",
	"khaki":      "#F0E68C",
	"lavender":   "#E6E6FA",
	"salmon":     "#FA8072",
	"tan":        "#D2B48C",
	"tomato":     "#FF6347",
	"turquoise":  "#40E0D0",
	"violet":     "#EE82EE",
	"wheat":      "#F5DEB3",

	// Dark variants
	"darkblue":   "#00008B",
	"darkcyan":   "#008B8B",
	"darkgray":   "#A9A9A9",
	"darkgrey":   "#A9A9A9",
	"darkgreen":  "#006400",
	"darkred":    "#8B0000",
	"darkorange": "#FF8C00",

	// Light variants
	"lightblue":  "#ADD8E6",
	"lightcyan":  "#E0FFFF",
	"lightgray":  "#D3D3D3",
	"lightgrey":  "#D3D3D3",
	"lightgreen": "#90EE90",
	"lightpink":  "#FFB6C1",
}
