package style

import "testing"

func TestParseColor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Named colors
		{"red", "#FF0000"},
		{"RED", "#FF0000"},
		{"Red", "#FF0000"},
		{"green", "#008000"},
		{"blue", "#0000FF"},
		{"black", "#000000"},
		{"white", "#FFFFFF"},
		
		// Hex colors should pass through
		{"#FF0000", "#FF0000"},
		{"#ff0000", "#ff0000"},
		
		// ANSI codes should pass through
		{"1", "1"},
		{"255", "255"},
		
		// Empty should pass through
		{"", ""},
		
		// Unknown values should pass through
		{"notacolor", "notacolor"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseColor(tt.input)
			if got != tt.expected {
				t.Errorf("ParseColor(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
