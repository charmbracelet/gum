package style

import "testing"

// TestResolveColor asserts that common named colors (as reported in
// https://github.com/charmbracelet/gum/issues/980) resolve to a valid
// ANSI color value that lipgloss.Color understands, while values that are
// already valid (hex codes, ANSI numbers) or unrecognized pass through
// unchanged.
func TestResolveColor(t *testing.T) {
	for name, tt := range map[string]struct {
		in   string
		want string
	}{
		"named color lowercase":                              {in: "red", want: "1"},
		"named color uppercase":                              {in: "RED", want: "1"},
		"named color mixed case with surrounding whitespace": {in: " Blue ", want: "4"},
		"bright variant":                                     {in: "brightcyan", want: "14"},
		"hex passthrough":                                    {in: "#FF0000", want: "#FF0000"},
		"ansi number passthrough":                            {in: "212", want: "212"},
		"unknown name passthrough":                           {in: "notacolor", want: "notacolor"},
		"empty passthrough":                                  {in: "", want: ""},
	} {
		t.Run(name, func(t *testing.T) {
			if got := resolveColor(tt.in); got != tt.want {
				t.Errorf("resolveColor(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
