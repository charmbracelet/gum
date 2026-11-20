//go:build windows

package spin

import (
	"os"
	"os/exec"
)

func setupProcess(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return killProcess(cmd)
	}
}

func killProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return cmd.Process.Kill()
}

func abortProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return nil
	}
	return cmd.Process.Signal(os.Interrupt)
}
