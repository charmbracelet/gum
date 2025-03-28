package cursor

import "github.com/charmbracelet/bubbles/v2/cursor"

// Modes maps strings to cursor modes.
var Modes = map[string]cursor.Mode{
	"blink":  cursor.CursorBlink,
	"hide":   cursor.CursorHide,
	"static": cursor.CursorStatic,
}
