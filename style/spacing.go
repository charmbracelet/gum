package style

import (
	"strconv"
	"strings"
)

const (
	minTokens  = 1
	halfTokens = 2
	maxTokens  = 4
)

// ParsePadding parses 1 - 4 integers from a string and returns them in a top,
// right, bottom, left order for use in the lipgloss.Padding() method.
func ParsePadding(s string) (int, int, int, int) {
	var ints [maxTokens]int

	tokens := strings.Split(s, " ")

	if len(tokens) > maxTokens {
		return 0, 0, 0, 0
	}

	// All tokens must be an integer
	for i, token := range tokens {
		parsed, err := strconv.Atoi(token)
		if err != nil {
			return 0, 0, 0, 0
		}
		ints[i] = parsed
	}

	if len(tokens) == minTokens {
		return ints[0], ints[0], ints[0], ints[0]
	}

	if len(tokens) == halfTokens {
		return ints[0], ints[1], ints[0], ints[1]
	}

	if len(tokens) == maxTokens {
		return ints[0], ints[1], ints[2], ints[3]
	}

	return 0, 0, 0, 0
}

// parseMargin is an alias for parsePadding since they involve the same logic
// to parse integers to the same format.
var parseMargin = ParsePadding
