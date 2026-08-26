//go:build linux || solaris

package spin

import "golang.org/x/sys/unix"

// System V style termios ioctls. Used on Linux (and Android, which satisfies
// the linux build tag) and on Solaris (and illumos, which satisfies solaris).
//
// NOTE: this file must not be named tty_linux.go — the _linux filename suffix
// is an implicit build constraint that is ANDed with the //go:build line, which
// would silently drop the solaris half of the tag.
const (
	ioctlReadTermios = unix.TCGETS
	// ioctlWriteTermios applies the new settings immediately.
	ioctlWriteTermios = unix.TCSETS
	// ioctlWriteTermiosFlush applies the new settings and discards pending
	// input in the same operation, i.e. tcsetattr(3) with TCSAFLUSH.
	ioctlWriteTermiosFlush = unix.TCSETSF
)
