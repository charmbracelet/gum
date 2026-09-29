package filter

import (
	"reflect"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/sahilm/fuzzy"
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

func matchStrs(matches []fuzzy.Match) []string {
	out := make([]string, len(matches))
	for i, m := range matches {
		out[i] = m.Str
	}
	return out
}

func TestInitialChoicesNoStrictValue(t *testing.T) {
	options := []string{"alice", "carol", "bob"}

	t.Run("prefilled prefix used to print the first option", func(t *testing.T) {
		_, _, matches := initialChoices(Options{
			Options:   options,
			Value:     "al",
			Fuzzy:     true,
			Sort:      true,
			FuzzySort: true,
			Strict:    true,
		})
		if len(matches) == 0 || matches[0].Str != "alice" {
			t.Fatalf("strict mode should still prefer alice, got %v", matchStrs(matches))
		}
	})

	t.Run("accepting --value al prints al not alice", func(t *testing.T) {
		choices, _, matches := initialChoices(Options{
			Options:   options,
			Value:     "al",
			Fuzzy:     true,
			Sort:      true,
			FuzzySort: true,
			Strict:    false,
		})
		if choices["al"] != "al" {
			t.Fatalf("expected the typed value in the choice map, got %#v", choices["al"])
		}
		if len(matches) == 0 || matches[0].Str != "al" {
			t.Fatalf("Enter would print %v, want al", matchStrs(matches))
		}
	})

	t.Run("accepting --value dan prints dan not empty", func(t *testing.T) {
		_, _, matches := initialChoices(Options{
			Options:   options,
			Value:     "dan",
			Fuzzy:     true,
			Sort:      true,
			FuzzySort: true,
			Strict:    false,
		})
		if len(matches) == 0 {
			t.Fatal("Enter would print nothing")
		}
		if matches[0].Str != "dan" {
			t.Fatalf("Enter would print %q, want dan", matches[0].Str)
		}
	})
}

func TestBuildMatchesNoStrictValue(t *testing.T) {
	options := []string{"alice", "carol", "bob"}

	t.Run("strict value still selects the first option match", func(t *testing.T) {
		matches := buildMatches("al", options, true, true, true)
		if len(matches) == 0 {
			t.Fatal("expected option matches")
		}
		if matches[0].Str != "alice" {
			t.Fatalf("expected alice, got %q", matches[0].Str)
		}
	})

	t.Run("no-strict value that prefixes an option is accepted as-is", func(t *testing.T) {
		// Reproduces https://github.com/charmbracelet/gum/issues/957
		// `gum filter --no-strict --value al alice carol bob` then Enter
		// used to print "alice".
		matches := buildMatches("al", options, true, true, false)
		if len(matches) == 0 || matches[0].Str != "al" {
			t.Fatalf("expected first match %q, got %v", "al", matchStrs(matches))
		}
	})

	t.Run("no-strict value that matches nothing is still selectable", func(t *testing.T) {
		// `gum filter --no-strict --value dan alice carol bob` then Enter
		// used to print nothing.
		matches := buildMatches("dan", options, true, true, false)
		if len(matches) == 0 || matches[0].Str != "dan" {
			t.Fatalf("expected first match %q, got %v", "dan", matchStrs(matches))
		}
	})

	t.Run("empty value lists all options even when not strict", func(t *testing.T) {
		matches := buildMatches("", options, true, true, false)
		if got, want := matchStrs(matches), options; !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})

	t.Run("no-strict exact match includes the typed value", func(t *testing.T) {
		matches := buildMatches("al", options, false, true, false)
		if len(matches) == 0 || matches[0].Str != "al" {
			t.Fatalf("expected first match %q, got %v", "al", matchStrs(matches))
		}
	})
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
