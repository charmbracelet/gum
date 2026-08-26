//go:build linux || solaris || darwin || dragonfly || freebsd || netbsd || openbsd

package spin

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// hushStdin stops the terminal from echoing anything that arrives on stdin
// while the spinner is running, and returns a function that restores the
// original terminal state while discarding any input that queued up in the
// meantime (terminal capability report replies included).
//
// Bubble Tea queries the terminal for synchronized output (mode 2026) and
// unicode core (mode 2027) support when a program starts. Because gum spin
// runs its tea.Program with a nil input reader (the wrapped command owns
// stdin), nobody consumes the terminal's replies, and the tty line discipline
// echoes them into the spinner output (e.g. "^[[?2026;2$y^[[?2027;3$y").
// Silencing echo hides them while the spinner runs; restoring with a flushing
// tcsetattr then drops them from the input queue so the shell does not read
// them back as if the user had typed them.
//
// NOTE: only ECHO is cleared, and stdin stays in cooked mode, so the wrapped
// command keeps working. It does mean a wrapped command that prompts for
// non-secret input (a "yes/no" host key confirmation, say) will not echo what
// the user types for as long as the spinner is up. Password prompts are
// unaffected, as they clear ECHO themselves.
func hushStdin(f *os.File) (restore func(), err error) {
	fd := int(f.Fd()) //nolint:gosec // a file descriptor always fits in an int

	saved, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, fmt.Errorf("unable to read terminal state: %w", err)
	}

	hushed := *saved
	hushed.Lflag &^= unix.ECHO
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &hushed); err != nil {
		return nil, fmt.Errorf("unable to silence terminal echo: %w", err)
	}

	return func() {
		// Restore the original state and discard any unread terminal
		// report replies in one atomic operation, so there is no window
		// in which echo is back on but the replies are still queued.
		if err := unix.IoctlSetTermios(fd, ioctlWriteTermiosFlush, saved); err != nil {
			// Never leave the terminal without echo: fall back to a
			// plain, non-flushing restore.
			_ = unix.IoctlSetTermios(fd, ioctlWriteTermios, saved)
		}
	}, nil
}
