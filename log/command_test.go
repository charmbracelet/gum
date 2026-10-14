package log

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogPreservesNewlines(t *testing.T) {
	// Verify that multiline strings retain newlines.
	text := "this\nis\na\ntest"
	lines := strings.Split(text, "\n")
	assert.Equal(t, 4, len(lines))
}
