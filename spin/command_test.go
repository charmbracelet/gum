package spin

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/gum/internal/exit"
)

func TestRunDoesNotPrintSuccessfulCommandOutputWithoutTTY(t *testing.T) {
	withRedirectedStandardStreams(t, func(stdoutPath, stderrPath string) {
		err := Options{
			Command:   []string{"sh", "-c", "printf 'stdout'; printf 'stderr' >&2"},
			ShowError: true,
		}.Run()
		if err != nil && !isExitCode(err, 0) {
			t.Fatalf("Run() error = %v", err)
		}

		assertFileContents(t, stdoutPath, "")
		assertFileContents(t, stderrPath, "")
	})
}

func isExitCode(err error, code int) bool {
	var ex exit.ErrExit
	return errors.As(err, &ex) && int(ex) == code
}

func TestRunPrintsBufferedOutputOnFailureWithoutTTY(t *testing.T) {
	withRedirectedStandardStreams(t, func(stdoutPath, stderrPath string) {
		err := Options{
			Command:   []string{"sh", "-c", "printf 'stdout'; printf 'stderr' >&2; exit 7"},
			ShowError: true,
		}.Run()

		var ex exit.ErrExit
		if !errors.As(err, &ex) || int(ex) != 7 {
			t.Fatalf("Run() error = %v, want exit status 7", err)
		}

		assertFileContainsAll(t, stdoutPath, "stdout", "stderr")
		assertFileContents(t, stderrPath, "")
	})
}

func withRedirectedStandardStreams(t *testing.T, fn func(stdoutPath, stderrPath string)) {
	t.Helper()

	tmpDir := t.TempDir()
	stdoutFile, err := os.Create(filepath.Join(tmpDir, "stdout"))
	if err != nil {
		t.Fatalf("Create(stdout) error = %v", err)
	}
	defer stdoutFile.Close() //nolint:errcheck

	stderrFile, err := os.Create(filepath.Join(tmpDir, "stderr"))
	if err != nil {
		t.Fatalf("Create(stderr) error = %v", err)
	}
	defer stderrFile.Close() //nolint:errcheck

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	os.Stdout = stdoutFile
	os.Stderr = stderrFile
	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
	}()

	fn(stdoutFile.Name(), stderrFile.Name())
}

func assertFileContents(t *testing.T, path string, want string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	if got := string(data); got != want {
		t.Fatalf("contents(%q) = %q, want %q", path, got, want)
	}
}

func assertFileContainsAll(t *testing.T, path string, want ...string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}

	got := string(data)
	for _, part := range want {
		if !strings.Contains(got, part) {
			t.Fatalf("contents(%q) = %q, missing %q", path, got, part)
		}
	}
}
