package style

import (
	"strconv"
	"strings"
)

var parseMargin = parsePadding

func parsePadding(s string) (int, int, int, int) {
	var ints []int

	tokens := strings.Split(s, " ")

	// All tokens must be an integer
	for _, token := range tokens {
		parsed, err := strconv.Atoi(token)
		if err != nil {
			return 0, 0, 0, 0
		}
		ints = append(ints, parsed)
	}

	if len(tokens) == 1 {
		return ints[0], ints[0], ints[0], ints[0]
	}

	if len(tokens) == 2 {
		return ints[0], ints[1], ints[0], ints[1]
	}

	if len(tokens) == 4 {
		return ints[0], ints[1], ints[2], ints[3]
	}

	return 0, 0, 0, 0
}
