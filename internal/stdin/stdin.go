// Package stdin handles processing input from stdin.
package stdin

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

type options struct {
	ansiStrip  bool
	singleLine bool
}

// Option is a read option.
type Option func(*options)

// StripANSI optionally strips ansi sequences.
func StripANSI(b bool) Option {
	return func(o *options) {
		o.ansiStrip = b
	}
}

// SingleLine reads a single line.
func SingleLine(b bool) Option {
	return func(o *options) {
		o.singleLine = b
	}
}

// Read reads input from an stdin pipe.
func Read(opts ...Option) (string, error) {
	if IsEmpty() {
		return "", fmt.Errorf("stdin is empty")
	}

	options := options{}
	for _, opt := range opts {
		opt(&options)
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	if options.singleLine {
		line, _, err := reader.ReadLine()
		if err != nil {
			return "", fmt.Errorf("failed to read line: %w", err)
		}
		_, err = b.Write(line)
		if err != nil {
			return "", fmt.Errorf("failed to write: %w", err)
		}
	}

	for !options.singleLine {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return "", fmt.Errorf("failed to write rune: %w", err)
		}
	}

	s := strings.TrimSpace(b.String())
	if options.ansiStrip {
		return ansi.Strip(s), nil
	}
	return s, nil
}

// IsEmpty returns whether stdin is empty.
func IsEmpty() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return true
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return true
	}

	return false
}
