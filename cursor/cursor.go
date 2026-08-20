// Package cursor provides cursor modes.
package cursor

import (
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	"charm.land/gum/v2/style"
)

// TextInput applies the given cursor mode and style to the text input.
func TextInput(i *textinput.Model, mode string, s style.Styles) {
	styles := i.Styles()
	styles.Cursor.Color = s.ToLipgloss().GetForeground()
	styles.Cursor.Blink = mode == "blink"
	i.SetStyles(styles)
	i.SetVirtualCursor(mode != "hide")
}

// TextArea applies the given cursor mode and style to the text area.
func TextArea(a *textarea.Model, mode string, s style.Styles) {
	styles := a.Styles()
	styles.Cursor.Color = s.ToLipgloss().GetForeground()
	styles.Cursor.Blink = mode == "blink"
	a.SetStyles(styles)
	a.SetVirtualCursor(mode != "hide")
}
