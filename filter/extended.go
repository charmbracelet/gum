package filter

import (
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/rivo/uniseg"
	"github.com/sahilm/fuzzy"
)

type searchTerm struct {
	text  string
	exact bool
}

// parseSearchTerms accepts incomplete input while the user is typing. An
// unescaped leading apostrophe makes a term exact; whitespace separates terms.
func parseSearchTerms(query string) []searchTerm {
	var terms []searchTerm
	var word strings.Builder
	exact := false
	flush := func() {
		if word.Len() > 0 {
			terms = append(terms, searchTerm{text: word.String(), exact: exact})
		}
		word.Reset()
		exact = false
	}
	runes := []rune(query)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if unicode.IsSpace(r) {
			flush()
			continue
		}
		if r == '\\' && i+1 < len(runes) {
			next := runes[i+1]
			if unicode.IsSpace(next) || next == '\\' || next == '\'' {
				word.WriteRune(next)
				i++
				continue
			}
		}
		if r == '\'' && word.Len() == 0 && !exact {
			exact = true
			continue
		}
		word.WriteRune(r)
	}
	flush()
	return terms
}

type matchSource []fuzzy.Match

func (s matchSource) Len() int            { return len(s) }
func (s matchSource) String(i int) string { return s[i].Str }

// extendedMatches intersects terms without imposing their order within an
// option. Each pass searches only the surviving options. Scores and highlights
// are combined, while original option indexes and stable ties are preserved.
func extendedMatches(query string, choices []string, fuzzyMode, sorted bool) []fuzzy.Match {
	matches := make([]fuzzy.Match, len(choices))
	for i, choice := range choices {
		matches[i] = fuzzy.Match{Str: choice, Index: i}
	}
	for _, term := range parseSearchTerms(query) {
		var next []fuzzy.Match
		if term.exact || !fuzzyMode {
			next = extendedExactMatches(term.text, matchSource(matches))
		} else {
			next = fuzzy.FindFromNoSort(term.text, matchSource(matches))
		}
		for i := range next {
			previous := matches[next[i].Index]
			next[i].Index = previous.Index
			next[i].Score += previous.Score
			next[i].MatchedIndexes = append(next[i].MatchedIndexes, previous.MatchedIndexes...)
			slices.Sort(next[i].MatchedIndexes)
			next[i].MatchedIndexes = slices.Compact(next[i].MatchedIndexes)
		}
		matches = next
		if len(matches) == 0 {
			break
		}
	}
	if sorted && fuzzyMode {
		sort.Stable(fuzzy.Matches(matches))
	}
	return matches
}

// extendedExactMatches folds runes rather than lowercasing the whole option:
// folding can change UTF-8 byte lengths, but rendering needs original offsets.
func extendedExactMatches(pattern string, choices matchSource) []fuzzy.Match {
	var matches []fuzzy.Match
	runes := []rune(pattern)
	indexes := make([]int, len(runes))
	for i, choice := range choices {
		for start := range choice.Str {
			pos := start
			matched := 0
			for _, want := range runes {
				if pos >= len(choice.Str) {
					break
				}
				got, size := utf8.DecodeRuneInString(choice.Str[pos:])
				if !equalFoldRune(got, want) {
					break
				}
				indexes[matched] = pos
				matched++
				pos += size
			}
			if matched == len(runes) {
				matches = append(matches, fuzzy.Match{Str: choice.Str, Index: i, MatchedIndexes: slices.Clone(indexes)})
				break
			}
		}
	}
	return matches
}

func equalFoldRune(a, b rune) bool {
	if a == b {
		return true
	}
	for folded := unicode.SimpleFold(a); folded != a; folded = unicode.SimpleFold(folded) {
		if folded == b {
			return true
		}
	}
	return false
}

// extendedMatchedRanges converts byte indexes to half-open terminal-cell
// ranges, highlighting complete graphemes even when a term matches a combining
// mark or part of an emoji sequence.
func extendedMatchedRanges(text string, indexes []int) [][2]int {
	var ranges [][2]int
	graphemes := uniseg.NewGraphemes(text)
	cell, match := 0, 0
	for match < len(indexes) && graphemes.Next() {
		_, end := graphemes.Positions()
		width := graphemes.Width()
		if indexes[match] < end {
			if len(ranges) > 0 && ranges[len(ranges)-1][1] == cell {
				ranges[len(ranges)-1][1] = cell + width
			} else {
				ranges = append(ranges, [2]int{cell, cell + width})
			}
			for match < len(indexes) && indexes[match] < end {
				match++
			}
		}
		cell += width
	}
	return ranges
}
