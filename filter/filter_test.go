package filter

import (
	"reflect"
	"testing"
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
