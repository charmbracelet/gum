package style

import "testing"

func TestParsePadding(t *testing.T) {
	for name, tt := range map[string]struct {
		in                       string
		top, right, bottom, left int
	}{
		"single token": {
			in:     "2",
			top:    2,
			right:  2,
			bottom: 2,
			left:   2,
		},
		"two tokens": {
			in:     "1 3",
			top:    1,
			right:  3,
			bottom: 1,
			left:   3,
		},
		"four tokens": {
			in:     "1 2 3 4",
			top:    1,
			right:  2,
			bottom: 3,
			left:   4,
		},
		"three tokens (unsupported, returns zeros)": {
			in: "1 2 3",
		},
		"five tokens (over max)": {
			in: "1 2 3 4 5",
		},
		"non-integer token": {
			in: "abc",
		},
		"empty string": {
			in: "",
		},
		"zero values": {
			in: "0 0 0 0",
		},
		"negative values": {
			in:     "-1 -2",
			top:    -1,
			right:  -2,
			bottom: -1,
			left:   -2,
		},
	} {
		t.Run(name, func(t *testing.T) {
			top, right, bottom, left := ParsePadding(tt.in)
			if top != tt.top || right != tt.right || bottom != tt.bottom || left != tt.left {
				t.Errorf("expected (%d, %d, %d, %d), got (%d, %d, %d, %d)",
					tt.top, tt.right, tt.bottom, tt.left,
					top, right, bottom, left)
			}
		})
	}
}

func TestParseMargin(t *testing.T) {
	// parseMargin is an alias of ParsePadding; verify aliasing holds.
	top, right, bottom, left := parseMargin("2 4")
	if top != 2 || right != 4 || bottom != 2 || left != 4 {
		t.Errorf("expected (2, 4, 2, 4), got (%d, %d, %d, %d)", top, right, bottom, left)
	}
}
