// Package cursor provides cursor modes.
package cursor

import (
	"github.com/charmbracelet/bubbles/cursor"
)

// Modes maps strings to cursor modes.
var Modes = map[string]cursor.Mode{
	"blink":  cursor.CursorBlink,
	"hide":   cursor.CursorHide,
	"static": cursor.CursorStatic,
}
