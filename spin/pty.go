package spin

import (
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/charmbracelet/x/xpty"
)

func openPty(f *os.File) (pty xpty.Pty, err error) {
	width, height, err := term.GetSize(f.Fd())
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	pty, err = xpty.NewPty(width, height)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return pty, nil
}
