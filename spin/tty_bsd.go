//go:build darwin || freebsd || netbsd || openbsd || dragonfly

package spin

import "golang.org/x/sys/unix"

const (
	ioctlReadTermios  = unix.TIOCGETA
	ioctlWriteTermios = unix.TIOCSETA
)

// flushInput discards pending input on the given file descriptor.
func flushInput(fd int) error {
	return unix.IoctlSetInt(fd, unix.TIOCFLUSH, unix.TCIFLUSH)
}
