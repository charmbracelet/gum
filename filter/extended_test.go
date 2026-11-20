package filter

import (
	"fmt"
	"reflect"
	"testing"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"github.com/sahilm/fuzzy"
)

func TestExtendedQuery(t *testing.T) {
	for _, tc := range []struct {
		query string
		want  []searchTerm
	}{
		{"", nil},
		{"  \t\u2003", nil},
		{"'", nil},
		{"src model", []searchTerm{{text: "src"}, {text: "model"}}},
		{"model src", []searchTerm{{text: "model"}, {text: "src"}}},
		{"src\u2003'model", []searchTerm{{text: "src"}, {text: "model", exact: true}}},
		{`'user\ model`, []searchTerm{{text: "user model", exact: true}}},
		{`user\ model`, []searchTerm{{text: "user model"}}},
		{`\'model`, []searchTerm{{text: "'model"}}},
		{"don't", []searchTerm{{text: "don't"}}},
		{`C:\Users\me`, []searchTerm{{text: `C:\Users\me`}}},
		{`path\\name`, []searchTerm{{text: `path\name`}}},
		{`path\`, []searchTerm{{text: `path\`}}},
		{"src ' ", []searchTerm{{text: "src"}}},
	} {
		t.Run(tc.query, func(t *testing.T) {
			if got := parseSearchTerms(tc.query); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("parseSearchTerms(%q) = %#v, want %#v", tc.query, got, tc.want)
			}
		})
	}
}

func TestExtendedMatches(t *testing.T) {
	choices := []string{"docs/user_model.md", "src/user_model.go", "src/m_o_d_e_l.go", "src/main.go"}
	for _, tc := range []struct {
		query string
		fuzzy bool
		want  []string
	}{
		{"src model", true, []string{"src/user_model.go", "src/m_o_d_e_l.go"}},
		{"model src", true, []string{"src/user_model.go", "src/m_o_d_e_l.go"}},
		{"src 'model", true, []string{"src/user_model.go"}},
		{"src model", false, []string{"src/user_model.go"}},
		{"SRC 'MODEL", true, []string{"src/user_model.go"}},
		{"'missing src", true, []string{}},
		{" ", true, choices},
		{"'", true, choices},
	} {
		t.Run(fmt.Sprintf("%s/fuzzy=%v", tc.query, tc.fuzzy), func(t *testing.T) {
			matches := extendedMatches(tc.query, choices, tc.fuzzy, false)
			got := make([]string, len(matches))
			for i, match := range matches {
				got[i] = match.Str
				if choices[match.Index] != match.Str {
					t.Errorf("wrong original index: %+v", match)
				}
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestExtendedMatchingDetails(t *testing.T) {
	t.Run("multiple terms combine rankings or preserve input order", func(t *testing.T) {
		choices := []string{"s_r_c/m_o_d_e_l", "src/model", "src/m_o_d_e_l"}
		for _, sorted := range []bool{false, true} {
			want := []int{0, 1, 2}
			if sorted {
				want = []int{1, 0, 2}
			}
			matches := extendedMatches("src model", choices, true, sorted)
			got := make([]int, len(matches))
			for i, match := range matches {
				got[i] = match.Index
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("sorted=%v: got indexes %v, want %v", sorted, got, want)
			}
		}
	})
	t.Run("single term retains fuzzy ranking", func(t *testing.T) {
		choices := []string{"unrelated", "model", "m_o_d_e_l", "userModel", "model"}
		if got, want := extendedMatches("mod", choices, true, true), fuzzy.Find("mod", choices); !reflect.DeepEqual(got, []fuzzy.Match(want)) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
	t.Run("overlapping terms have sorted unique highlights", func(t *testing.T) {
		matches := extendedMatches("'model 'del", []string{"model"}, true, true)
		if len(matches) != 1 || !reflect.DeepEqual(matches[0].MatchedIndexes, []int{0, 1, 2, 3, 4}) {
			t.Fatalf("unexpected matches: %+v", matches)
		}
	})
	t.Run("duplicates retain their original indexes", func(t *testing.T) {
		matches := extendedMatches("'model", []string{"model", "other", "model"}, true, true)
		if len(matches) != 2 || matches[0].Index != 0 || matches[1].Index != 2 {
			t.Fatalf("unexpected matches: %+v", matches)
		}
	})
	t.Run("case folding retains original UTF-8 byte offsets", func(t *testing.T) {
		matches := extendedMatches("'kelvin 'go", []string{"Kelvin.go"}, true, false)
		want := []int{0, 3, 4, 5, 6, 7, 9, 10}
		if len(matches) != 1 || !reflect.DeepEqual(matches[0].MatchedIndexes, want) {
			t.Fatalf("got %+v, want indexes %v", matches, want)
		}
	})
	t.Run("Unicode simple fold includes final sigma", func(t *testing.T) {
		if matches := extendedMatches("'σ", []string{"ς", "Σ"}, false, false); len(matches) != 2 {
			t.Fatalf("unexpected matches: %+v", matches)
		}
	})
	t.Run("exact escaped whitespace", func(t *testing.T) {
		matches := extendedMatches(`'user\ model`, []string{"user_model", "user model"}, true, false)
		if len(matches) != 1 || matches[0].Str != "user model" {
			t.Fatalf("unexpected matches: %+v", matches)
		}
	})
}

func TestExtendedMatchedRanges(t *testing.T) {
	for _, tc := range []struct {
		text    string
		indexes []int
		want    [][2]int
	}{
		{"abc def", []int{0, 1, 2, 4, 6}, [][2]int{{0, 3}, {4, 5}, {6, 7}}},
		{"界model", []int{0, 3, 4}, [][2]int{{0, 4}}},
		{"e\u0301lan", []int{1}, [][2]int{{0, 1}}},
		{"👩‍💻 code", []int{4, 7}, [][2]int{{0, 2}}},
		{"plain", nil, nil},
	} {
		t.Run(tc.text, func(t *testing.T) {
			if got := extendedMatchedRanges(tc.text, tc.indexes); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestExtendedFilteringWhileTyping(t *testing.T) {
	i := textinput.New()
	i.Focus()
	v := viewport.New(viewport.WithWidth(60), viewport.WithHeight(5))
	choices := []string{"src/user_model.go", "src/main.go", "docs/model.md"}
	m := model{
		textinput: i, viewport: &v, filteringChoices: choices,
		extended: true, fuzzy: true, sort: true, strict: true, limit: 2,
		selected: make(map[string]struct{}), keymap: defaultKeymap(),
	}
	for _, r := range "src 'model" {
		updated, _ := m.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
		m = updated.(model)
	}
	if len(m.matches) != 1 || m.matches[0].Str != "src/user_model.go" {
		t.Fatalf("unexpected matches after typing: %+v", m.matches)
	}
	m.ToggleSelection()
	if _, ok := m.selected["src/user_model.go"]; !ok {
		t.Fatal("extended filtering lost the selected option")
	}
}

func BenchmarkExtendedMatches(b *testing.B) {
	for _, size := range []int{1000, 10000} {
		choices := make([]string, size)
		for i := range choices {
			choices[i] = fmt.Sprintf("src/package-%04d/user_model.go", i)
		}
		for _, query := range []string{"model", "src 'model", "'missing model"} {
			b.Run(fmt.Sprintf("%d/%s", size, query), func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					extendedMatches(query, choices, true, true)
				}
			})
		}
	}
}

func FuzzExtendedMatches(f *testing.F) {
	f.Add("src 'model", "src/user_model.go", true)
	f.Add("'kelvin", "Kelvin.go", false)
	f.Add("'\u0301", "e\u0301", true)
	f.Add(`\'`, "'model", false)
	f.Add("", "", true)
	f.Fuzz(func(t *testing.T, query, choice string, fuzzyMode bool) {
		if len(query) > 256 || len(choice) > 2048 {
			t.Skip()
		}
		matches := extendedMatches(query, []string{choice}, fuzzyMode, true)
		if len(matches) == 0 {
			return
		}
		if len(matches) != 1 || matches[0].Str != choice || matches[0].Index != 0 {
			t.Fatalf("invalid result: %+v", matches)
		}
		previous := -1
		for _, index := range matches[0].MatchedIndexes {
			if index <= previous || index >= len(choice) {
				t.Fatalf("invalid highlight offsets: %v for %q", matches[0].MatchedIndexes, choice)
			}
			previous = index
		}
	})
}
