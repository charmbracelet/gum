package spin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/gum/v2/internal/exit"
)

func TestTimeoutTerminatesChildProcess(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping Unix-specific command on Windows")
	}

	start := time.Now()
	opts := Options{
		Command: []string{"sleep", "5"},
		Timeout: 200 * time.Millisecond,
	}

	err := opts.Run()
	elapsed := time.Since(start)

	if elapsed > 3*time.Second {
		t.Fatalf("command ran for %v, expected termination near 200ms", elapsed)
	}

	if err == nil {
		t.Fatal("expected error due to timeout, got nil")
	}

	if !errors.Is(err, tea.ErrProgramKilled) {
		t.Fatalf("expected ErrProgramKilled, got: %v", err)
	}
}

func TestTimeoutTerminatesProcessTree(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping Unix-specific command on Windows")
	}

	pidFile := filepath.Join(t.TempDir(), "child.pid")
	// Write subshell PID to pidFile and sleep in child process
	cmdStr := fmt.Sprintf("echo $$ > %s && sleep 10", pidFile)

	start := time.Now()
	opts := Options{
		Command: []string{"sh", "-c", cmdStr},
		Timeout: 300 * time.Millisecond,
	}

	err := opts.Run()
	elapsed := time.Since(start)

	if elapsed > 3*time.Second {
		t.Fatalf("command ran for %v, expected termination near 300ms", elapsed)
	}

	if err == nil {
		t.Fatal("expected error due to timeout, got nil")
	}

	if !errors.Is(err, tea.ErrProgramKilled) {
		t.Fatalf("expected ErrProgramKilled, got: %v", err)
	}

	// Read recorded PID
	data, readErr := os.ReadFile(pidFile)
	if readErr != nil {
		t.Fatalf("failed to read PID file: %v", readErr)
	}
	pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
	if parseErr != nil {
		t.Fatalf("failed to parse PID: %v", parseErr)
	}

	// Verify the process is no longer running as an active process
	time.Sleep(100 * time.Millisecond)
	if isProcessRunning(pid) {
		t.Fatalf("expected process %d to be terminated, but it is still running", pid)
	}
}

func isProcessRunning(pid int) bool {
	// Check /proc/<pid>/status on Linux if available
	procStatus := fmt.Sprintf("/proc/%d/status", pid)
	if content, err := os.ReadFile(procStatus); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "State:") {
				// State 'Z' is a zombie process whose execution has ended
				if strings.Contains(line, "Z") {
					return false
				}
				// Any other state (like 'S' sleeping or 'R' running) means it is alive
				return true
			}
		}
	}

	// Fallback to syscall.Kill(pid, 0)
	err := syscall.Kill(pid, 0)
	return err == nil
}

func TestCommandSuccess(t *testing.T) {
	cmd := []string{"echo", "hello"}
	if runtime.GOOS == "windows" {
		cmd = []string{"cmd", "/c", "echo", "hello"}
	}

	opts := Options{
		Command: cmd,
		Timeout: 5 * time.Second,
	}

	err := opts.Run()
	var ex exit.ErrExit
	if errors.As(err, &ex) && int(ex) != 0 {
		t.Fatalf("expected exit code 0, got: %d", int(ex))
	}
}

func TestCommandExitCode(t *testing.T) {
	cmd := []string{"sh", "-c", "exit 42"}
	if runtime.GOOS == "windows" {
		cmd = []string{"cmd", "/c", "exit", "42"}
	}

	opts := Options{
		Command: cmd,
		Timeout: 5 * time.Second,
	}

	err := opts.Run()
	var ex exit.ErrExit
	if !errors.As(err, &ex) {
		t.Fatalf("expected ErrExit, got: %v", err)
	}
	if int(ex) != 42 {
		t.Fatalf("expected exit code 42, got %d", int(ex))
	}
}
