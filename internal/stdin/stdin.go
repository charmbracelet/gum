package stdin

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// Read reads input from an stdin pipe.
func Read() (string, error) {
	if IsEmpty() {
		return "", fmt.Errorf("stdin is empty")
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

	return strings.TrimSuffix(ansi.Strip(b.String()), "\n"), nil
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
