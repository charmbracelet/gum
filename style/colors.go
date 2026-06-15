package style

import "strings"

// namedColors maps common color names to hex values for gum style flags.
var namedColors = map[string]string{
	"black":   "#000000",
	"red":     "#FF0000",
	"green":   "#00FF00",
	"yellow":  "#FFFF00",
	"blue":    "#0000FF",
	"magenta": "#FF00FF",
	"cyan":    "#00FFFF",
	"white":   "#FFFFFF",
}

func resolveColor(value string) string {
	if value == "" {
		return value
	}
	if hex, ok := namedColors[strings.ToLower(value)]; ok {
		return hex
	}
	return value
}
