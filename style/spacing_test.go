package style

import "testing"

func TestParsePadding(t *testing.T) {
	tests := []struct {
		input            string
		wantTop          int
		wantRight        int
		wantBottom       int
		wantLeft         int
	}{
		{"0", 0, 0, 0, 0},
		{"1", 1, 1, 1, 1},
		{"1 2", 1, 2, 1, 2},
		{"1 2 3 4", 1, 2, 3, 4},
		{"", 0, 0, 0, 0},
		{"foo", 0, 0, 0, 0},
		{"1 2 3 4 5", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			top, right, bottom, left := ParsePadding(tt.input)
			if top != tt.wantTop || right != tt.wantRight || bottom != tt.wantBottom || left != tt.wantLeft {
				t.Errorf("ParsePadding(%q) = (%d,%d,%d,%d), want (%d,%d,%d,%d)",
					tt.input, top, right, bottom, left, tt.wantTop, tt.wantRight, tt.wantBottom, tt.wantLeft)
			}
		})
	}
}
