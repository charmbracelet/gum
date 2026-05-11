package filter

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestMatchedRanges(t *testing.T) {
	for name, tt := range map[string]struct {
		in  []int
		out [][2]int
	}{
		"empty": {
			in:  []int{},
			out: [][2]int{},
		},
		"one char": {
			in:  []int{1},
			out: [][2]int{{1, 1}},
		},
		"2 char range": {
			in:  []int{1, 2},
			out: [][2]int{{1, 2}},
		},
		"multiple char range": {
			in:  []int{1, 2, 3, 4, 5, 6},
			out: [][2]int{{1, 6}},
		},
		"multiple char ranges": {
			in:  []int{1, 2, 3, 5, 6, 10, 11, 12, 13, 23, 24, 40, 42, 43, 45, 52},
			out: [][2]int{{1, 3}, {5, 6}, {10, 13}, {23, 24}, {40, 40}, {42, 43}, {45, 45}, {52, 52}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			match := matchedRanges(tt.in)
			if !reflect.DeepEqual(match, tt.out) {
				t.Errorf("expected %v, got %v", tt.out, match)
			}
		})
	}
}

func TestByteToChar(t *testing.T) {
	stStr := "\x1b[90m\ue615\x1b[39m \x1b[3m\x1b[32mDow\x1b[0m\x1b[90m\x1b[39m\x1b[3wnloads"
	str := " Downloads"
	rng := [2]int{4, 7}
	expect := "Dow"

	if got := str[rng[0]:rng[1]]; got != expect {
		t.Errorf("expected %q, got %q", expect, got)
	}

	start, stop := bytePosToVisibleCharPos(str, rng)
	if got := ansi.Strip(ansi.Cut(stStr, start, stop)); got != expect {
		t.Errorf("expected %+q, got %+q", expect, got)
	}
}
