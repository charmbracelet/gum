package style

import "testing"

func TestParsePadding(t *testing.T) {
	tests := []struct {
		name                    string
		in                      string
		top, right, bottom, left int
	}{
		{"one token applies to all sides", "1", 1, 1, 1, 1},
		{"two tokens vertical horizontal", "1 2", 1, 2, 1, 2},
		{"three tokens top horizontal bottom", "1 2 3", 1, 2, 3, 2},
		{"four tokens clockwise", "1 2 3 4", 1, 2, 3, 4},
		{"too many tokens is zero", "1 2 3 4 5", 0, 0, 0, 0},
		{"non-integer is zero", "1 a", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			top, right, bottom, left := ParsePadding(tt.in)
			if top != tt.top || right != tt.right || bottom != tt.bottom || left != tt.left {
				t.Errorf("ParsePadding(%q) = (%d,%d,%d,%d), want (%d,%d,%d,%d)",
					tt.in, top, right, bottom, left, tt.top, tt.right, tt.bottom, tt.left)
			}
		})
	}
}
