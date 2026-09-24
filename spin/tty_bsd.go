//go:build darwin || dragonfly || freebsd || netbsd || openbsd

package spin

import "golang.org/x/sys/unix"

// BSD style termios ioctls. Used on macOS (and iOS, which satisfies the darwin
// build tag) and on the BSDs.
const (
	ioctlReadTermios = unix.TIOCGETA
	// ioctlWriteTermios applies the new settings immediately.
	ioctlWriteTermios = unix.TIOCSETA
	// ioctlWriteTermiosFlush applies the new settings and discards pending
	// input in the same operation, i.e. tcsetattr(3) with TCSAFLUSH.
	ioctlWriteTermiosFlush = unix.TIOCSETAF
)
