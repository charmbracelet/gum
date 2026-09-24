package filter

import (
	"testing"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/sahilm/fuzzy"
)

// newTestModel builds a model with just the fields the cursor and selection
// helpers need, so they can be exercised without starting a bubbletea program.
func newTestModel(options []string, limit int) model {
	v := viewport.New(80, len(options))
	return model{
		viewport: &v,
		matches:  matchAll(options),
		selected: make(map[string]struct{}),
		limit:    limit,
	}
}

func TestMatchAll(t *testing.T) {
	options := []string{"foo", "bar", "baz"}
	matches := matchAll(options)

	if len(matches) != len(options) {
		t.Fatalf("expected %d matches but got %d", len(options), len(matches))
	}
	for i, option := range options {
		if matches[i].Str != option {
			t.Errorf("expected match %d to be %q but got %q", i, option, matches[i].Str)
		}
		if len(matches[i].MatchedIndexes) != 0 {
			t.Errorf("expected match %d to have no matched indexes", i)
		}
	}

	if got := matchAll(nil); len(got) != 0 {
		t.Errorf("expected no matches for a nil slice but got %d", len(got))
	}
}

func TestExactMatches(t *testing.T) {
	choices := []string{"Foo", "bar", "foobar", "BAZ"}

	for name, tt := range map[string]struct {
		search   string
		expected []string
	}{
		"matches are case insensitive": {"foo", []string{"Foo", "foobar"}},
		"uppercase search":             {"FOO", []string{"Foo", "foobar"}},
		"substring in the middle":      {"oba", []string{"foobar"}},
		"matches the tail":             {"bar", []string{"bar", "foobar"}},
		"no match":                     {"zzz", nil},
		"empty search matches all":     {"", []string{"Foo", "bar", "foobar", "BAZ"}},
	} {
		t.Run(name, func(t *testing.T) {
			matches := exactMatches(tt.search, choices)

			if len(matches) != len(tt.expected) {
				t.Fatalf("expected %d matches but got %d", len(tt.expected), len(matches))
			}
			for i, want := range tt.expected {
				if matches[i].Str != want {
					t.Errorf("expected match %d to be %q but got %q", i, want, matches[i].Str)
				}
			}
		})
	}
}

func TestExactMatchesIndexes(t *testing.T) {
	matches := exactMatches("oob", []string{"foobar"})

	if len(matches) != 1 {
		t.Fatalf("expected one match but got %d", len(matches))
	}
	// the original index into the choices slice is preserved
	if matches[0].Index != 0 {
		t.Errorf("expected the choice index to be 0 but got %d", matches[0].Index)
	}
	// "oob" starts at offset 1 of "foobar"
	expected := []int{1, 2, 3}
	if len(matches[0].MatchedIndexes) != len(expected) {
		t.Fatalf("expected %d matched indexes but got %v", len(expected), matches[0].MatchedIndexes)
	}
	for i, want := range expected {
		if matches[0].MatchedIndexes[i] != want {
			t.Errorf("expected matched index %d to be %d but got %d", i, want, matches[0].MatchedIndexes[i])
		}
	}
}

func TestExactMatchesUsesFirstOccurrence(t *testing.T) {
	// "an" occurs at offsets 1 and 3 of "banana"; the first one is reported
	matches := exactMatches("an", []string{"banana"})

	if len(matches) != 1 {
		t.Fatalf("expected one match but got %d", len(matches))
	}
	expected := []int{1, 2}
	if len(matches[0].MatchedIndexes) != len(expected) {
		t.Fatalf("expected %d matched indexes but got %v", len(expected), matches[0].MatchedIndexes)
	}
	for i, want := range expected {
		if matches[0].MatchedIndexes[i] != want {
			t.Errorf("expected matched index %d to be %d but got %d", i, want, matches[0].MatchedIndexes[i])
		}
	}
}

func TestExactMatchesPreservesChoiceIndex(t *testing.T) {
	matches := exactMatches("bar", []string{"foo", "bar", "baz", "foobar"})

	if len(matches) != 2 {
		t.Fatalf("expected two matches but got %d", len(matches))
	}
	if matches[0].Index != 1 {
		t.Errorf("expected the first match to carry index 1 but got %d", matches[0].Index)
	}
	if matches[1].Index != 3 {
		t.Errorf("expected the second match to carry index 3 but got %d", matches[1].Index)
	}
}

func TestCursorDownWraps(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 1)

	m.CursorDown()
	if m.cursor != 1 {
		t.Errorf("expected the cursor to move to 1 but it is %d", m.cursor)
	}
	m.CursorDown()
	if m.cursor != 2 {
		t.Errorf("expected the cursor to move to 2 but it is %d", m.cursor)
	}
	// moving past the last match wraps around to the first
	m.CursorDown()
	if m.cursor != 0 {
		t.Errorf("expected the cursor to wrap to 0 but it is %d", m.cursor)
	}
}

func TestCursorUpWraps(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 1)

	// moving up from the first match wraps around to the last
	m.CursorUp()
	if m.cursor != 2 {
		t.Errorf("expected the cursor to wrap to 2 but it is %d", m.cursor)
	}
	m.CursorUp()
	if m.cursor != 1 {
		t.Errorf("expected the cursor to move to 1 but it is %d", m.cursor)
	}
}

func TestCursorReversed(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 1)
	m.reverse = true

	// with a reversed list the directions are swapped
	m.CursorUp()
	if m.cursor != 1 {
		t.Errorf("expected the cursor to move to 1 but it is %d", m.cursor)
	}
	m.CursorDown()
	if m.cursor != 0 {
		t.Errorf("expected the cursor to move back to 0 but it is %d", m.cursor)
	}
}

func TestCursorWithNoMatches(t *testing.T) {
	m := newTestModel(nil, 1)

	// with nothing to move through the cursor must stay put
	m.CursorUp()
	m.CursorDown()
	if m.cursor != 0 {
		t.Errorf("expected the cursor to stay at 0 but it is %d", m.cursor)
	}
}

func TestToggleSelection(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 2)

	m.ToggleSelection()
	if _, ok := m.selected["a"]; !ok {
		t.Fatalf("expected %q to be selected", "a")
	}
	if m.numSelected != 1 {
		t.Fatalf("expected one selected item but got %d", m.numSelected)
	}

	// toggling the same item again deselects it
	m.ToggleSelection()
	if _, ok := m.selected["a"]; ok {
		t.Fatalf("expected %q to be deselected", "a")
	}
	if m.numSelected != 0 {
		t.Fatalf("expected no selected items but got %d", m.numSelected)
	}
}

func TestToggleSelectionRespectsLimit(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 2)

	m.ToggleSelection()
	m.CursorDown()
	m.ToggleSelection()
	m.CursorDown()
	m.ToggleSelection() // over the limit, must be ignored

	if m.numSelected != 2 {
		t.Fatalf("expected the limit of 2 to be respected but got %d", m.numSelected)
	}
	if _, ok := m.selected["c"]; ok {
		t.Fatalf("expected %q not to be selected once the limit was reached", "c")
	}
}

func TestSelectAll(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 10)

	m = m.selectAll()

	if m.numSelected != 3 {
		t.Fatalf("expected all three items to be selected but got %d", m.numSelected)
	}
	for _, want := range []string{"a", "b", "c"} {
		if _, ok := m.selected[want]; !ok {
			t.Errorf("expected %q to be selected", want)
		}
	}
}

func TestSelectAllStopsAtLimit(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 2)

	m = m.selectAll()

	if m.numSelected != 2 {
		t.Fatalf("expected the limit of 2 to be respected but got %d", m.numSelected)
	}
}

func TestSelectAllDoesNotDoubleCount(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 10)

	m.ToggleSelection() // selects "a"
	m = m.selectAll()

	if m.numSelected != 3 {
		t.Fatalf("expected three selected items but got %d", m.numSelected)
	}
}

func TestDeselectAll(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, 10)

	m = m.selectAll()
	m = m.deselectAll()

	if m.numSelected != 0 {
		t.Fatalf("expected no selected items but got %d", m.numSelected)
	}
	if len(m.selected) != 0 {
		t.Fatalf("expected the selection to be empty but it holds %d items", len(m.selected))
	}
}

func TestDefaultKeymap(t *testing.T) {
	k := defaultKeymap()

	if len(k.ShortHelp()) == 0 {
		t.Errorf("expected the short help to list some bindings")
	}
	if k.FullHelp() != nil {
		t.Errorf("expected the full help to be nil")
	}
	if !k.Submit.Enabled() {
		t.Errorf("expected the submit binding to be enabled")
	}
}

func TestMatchesAreUsableAsFuzzyMatches(t *testing.T) {
	// matchAll and exactMatches must both produce values the rest of the
	// package can treat uniformly as fuzzy.Match
	var all []fuzzy.Match
	all = append(all, matchAll([]string{"foo"})...)
	all = append(all, exactMatches("foo", []string{"foobar"})...)

	if len(all) != 2 {
		t.Fatalf("expected two matches but got %d", len(all))
	}
	if all[0].Str != "foo" || all[1].Str != "foobar" {
		t.Errorf("unexpected match strings: %q, %q", all[0].Str, all[1].Str)
	}
}
