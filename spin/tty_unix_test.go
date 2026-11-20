//go:build linux || solaris || darwin || dragonfly || freebsd || netbsd || openbsd

package spin

import (
	"testing"
	"time"

	"github.com/charmbracelet/x/xpty"
	"golang.org/x/sys/unix"
)

// newTestPty returns a pty whose slave stands in for the spinner's stdin.
func newTestPty(t *testing.T) *xpty.UnixPty {
	t.Helper()

	p, err := xpty.NewPty(80, 24)
	if err != nil {
		t.Skipf("cannot open a pty here: %v", err)
	}
	t.Cleanup(func() { _ = p.Close() })

	unixPty, ok := p.(*xpty.UnixPty)
	if !ok {
		t.Skip("not a unix pty")
	}
	return unixPty
}

func lflag(t *testing.T, fd int) uint64 {
	t.Helper()

	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		t.Fatalf("IoctlGetTermios: %v", err)
	}
	return uint64(termios.Lflag)
}

// queuedInput reports how many bytes are sitting in fd's input queue.
//
// The queue has to be read in non-canonical mode: the capability replies do not
// end in a newline, so a cooked-mode read would return nothing whether or not
// the queue was drained. Switching ICANON off does not disturb the queue.
func queuedInput(t *testing.T, fd int) int {
	t.Helper()

	saved, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		t.Fatalf("IoctlGetTermios: %v", err)
	}
	raw := *saved
	raw.Lflag &^= unix.ICANON
	raw.Cc[unix.VMIN] = 0
	raw.Cc[unix.VTIME] = 0
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &raw); err != nil {
		t.Fatalf("IoctlSetTermios: %v", err)
	}
	defer func() { _ = unix.IoctlSetTermios(fd, ioctlWriteTermios, saved) }()

	if err := unix.SetNonblock(fd, true); err != nil {
		t.Fatalf("SetNonblock: %v", err)
	}
	defer func() { _ = unix.SetNonblock(fd, false) }()

	buf := make([]byte, 256)
	total := 0
	for {
		n, err := unix.Read(fd, buf)
		if n > 0 {
			total += n
		}
		if n <= 0 || err != nil {
			return total
		}
	}
}

func TestHushStdinSilencesAndRestoresEcho(t *testing.T) {
	pty := newTestPty(t)
	slave := pty.Slave()
	fd := int(slave.Fd())

	before := lflag(t, fd)
	if before&unix.ECHO == 0 {
		t.Skip("pty starts without ECHO; nothing to silence")
	}

	restore, err := hushStdin(slave)
	if err != nil {
		t.Fatalf("hushStdin: %v", err)
	}

	if during := lflag(t, fd); during&unix.ECHO != 0 {
		t.Error("ECHO is still set while the spinner is running")
	}

	restore()

	if after := lflag(t, fd); after != before {
		t.Errorf("terminal state not restored: Lflag = %#x, want %#x", after, before)
	}
}

// TestHushStdinRestoreFlushesInput is the regression test for the terminal
// capability replies leaking out. Silencing echo keeps them off the screen
// while the spinner runs, but unless restore() also drains the input queue the
// shell reads them back afterwards as if the user had typed them.
func TestHushStdinRestoreFlushesInput(t *testing.T) {
	pty := newTestPty(t)
	slave := pty.Slave()
	fd := int(slave.Fd())

	restore, err := hushStdin(slave)
	if err != nil {
		t.Fatalf("hushStdin: %v", err)
	}

	// Stand in for the terminal answering Bubble Tea's mode 2026/2027 queries
	// while nobody is reading stdin.
	if _, err := pty.Master().Write([]byte("\x1b[?2026;2$y\x1b[?2027;3$y")); err != nil {
		t.Fatalf("write to pty master: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	restore()

	if n := queuedInput(t, fd); n != 0 {
		t.Errorf("restore() left %d bytes of terminal replies queued on stdin", n)
	}
}
