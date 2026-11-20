//go:build unix

package log

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestLogFileCreateMode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.log")

	old := syscall.Umask(0)
	t.Cleanup(func() { syscall.Umask(old) })

	opts := Options{
		Text:  []string{"hello"},
		File:  path,
		Level: "none",
	}
	if err := opts.Run(); err != nil {
		t.Fatalf("run log: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat log file: %v", err)
	}

	perm := info.Mode().Perm()
	if perm&0o077 != 0 {
		t.Fatalf("created file is group/world accessible: %04o", perm)
	}
	if perm != 0o600 {
		t.Fatalf("created file mode = %04o, want 0600", perm)
	}
}
