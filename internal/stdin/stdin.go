package stdin

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

var ErrEmpty = errors.New("stdin is empty")

// Read reads input from an stdin pipe.
func Read() (string, error) {
	if IsEmpty() {
		return "", ErrEmpty
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return "", fmt.Errorf("failed to write rune: %w", err)
		}
	}

	return strings.TrimSpace(b.String()), nil
}

// ReadStrip reads input from an stdin pipe and strips ansi sequences.
func ReadStrip() (string, error) {
	s, err := Read()
	return ansi.Strip(s), err
}

// ReadLine reads only one line and returns it.
func ReadLine() (string, error) {
	if IsEmpty() {
		return "", ErrEmpty
	}
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	return strings.TrimSpace(string(line)), err
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
