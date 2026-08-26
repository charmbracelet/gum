//go:build unix

package spin

import (
	"os"

	"golang.org/x/sys/unix"
)

// hushStdin stops the terminal from echoing anything that arrives on stdin
// while the spinner is running, and returns a function that discards any
// pending input (terminal capability report replies included) and restores
// the original terminal state.
//
// Bubble Tea queries the terminal for synchronized output (mode 2026) and
// unicode core (mode 2027) support when a program starts. Because gum spin
// runs its tea.Program with a nil input reader (the wrapped command owns
// stdin), nobody consumes the terminal's replies, and the tty line discipline
// echoes them into the spinner output (e.g. "^[[?2026;2$y^[[?2027;3$y").
// Silencing echo and flushing the input queue on exit keeps those replies
// invisible without stealing stdin from the wrapped command.
func hushStdin(f *os.File) (restore func(), err error) {
	fd := int(f.Fd())

	saved, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	hushed := *saved
	hushed.Lflag &^= unix.ECHO
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &hushed); err != nil {
		return nil, err
	}

	return func() {
		// Discard any unread terminal report replies so they are not
		// echoed or read back after the terminal state is restored.
		_ = flushInput(fd)
		_ = unix.IoctlSetTermios(fd, ioctlWriteTermios, saved)
	}, nil
}
