//go:build unix

package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/charmbracelet/x/xpty"
)

func TestConfirmExitStatus(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "gum")
	build := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, ".")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build gum: %v\n%s", err, output)
	}

	tests := []struct {
		name       string
		args       []string
		input      string
		signal     os.Signal
		wantStatus int
		wantOutput string
	}{
		{name: "sigterm-default-yes", signal: syscall.SIGTERM, wantStatus: 1},
		{name: "sigterm-default-no", args: []string{"--default=false"}, signal: syscall.SIGTERM, wantStatus: 1},
		{name: "sigterm-show-output", args: []string{"--show-output"}, signal: syscall.SIGTERM, wantStatus: 1},
		{name: "sigterm-custom-labels", args: []string{"--affirmative=Reboot", "--negative=Cancel"}, signal: syscall.SIGTERM, wantStatus: 1},
		{name: "enter-yes", input: "\r", wantStatus: 0},
		{name: "enter-no", args: []string{"--default=false"}, input: "\r", wantStatus: 1},
		{name: "explicit-yes", args: []string{"--default=false", "--show-output"}, input: "y", wantStatus: 0, wantOutput: "Synthetic confirmation? Yes"},
		{name: "explicit-no", args: []string{"--show-output"}, input: "n", wantStatus: 1, wantOutput: "Synthetic confirmation? No"},
		{name: "escape", input: "\x1b", wantStatus: 1},
		{name: "timeout-yes", args: []string{"--timeout=500ms", "--show-output"}, wantStatus: 0, wantOutput: "Synthetic confirmation? Yes"},
		{name: "timeout-no", args: []string{"--default=false", "--timeout=500ms"}, wantStatus: 1},
		{name: "sigint", signal: syscall.SIGINT, wantStatus: 130},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terminal, err := xpty.NewUnixPty(80, 24)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = terminal.Close() })
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			args := append([]string{"confirm", "Synthetic confirmation?"}, tt.args...)
			command := exec.CommandContext(ctx, binary, args...)
			command.Env = append(os.Environ(), "TERM=xterm-256color")
			if err := terminal.Start(command); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if command.ProcessState == nil {
					_ = command.Process.Kill()
					_ = command.Wait()
				}
			})
			// Only the child keeps the slave open, so reads finish when it exits.
			_ = terminal.Slave().Close()
			ready := make(chan struct{})
			output := make(chan string, 1)
			go func() {
				var text strings.Builder
				announced := false
				buffer := make([]byte, 4096)
				for {
					n, err := terminal.Read(buffer)
					text.Write(buffer[:n])
					if !announced && strings.Contains(text.String(), "Synthetic confirmation?") {
						close(ready)
						announced = true
					}
					if err != nil {
						output <- text.String()
						return
					}
				}
			}()
			if tt.signal != nil || tt.input != "" {
				select {
				case <-ready:
				case text := <-output:
					_ = command.Wait()
					t.Fatalf("confirmation exited before showing the prompt: %q", text)
				case <-ctx.Done():
					_ = command.Wait()
					t.Fatal("confirmation prompt did not appear")
				}
			}
			if tt.signal != nil {
				if err := command.Process.Signal(tt.signal); err != nil {
					t.Fatal(err)
				}
			} else if tt.input != "" {
				if _, err := terminal.Write([]byte(tt.input)); err != nil {
					t.Fatal(err)
				}
			}
			err = command.Wait()
			status := 0
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				status = exitErr.ExitCode()
			} else if err != nil {
				t.Fatal(err)
			}
			text := <-output
			if status != tt.wantStatus {
				t.Errorf("exit status = %d, want %d; output: %q", status, tt.wantStatus, text)
			}
			if tt.wantOutput != "" && !strings.Contains(text, tt.wantOutput) {
				t.Errorf("missing chosen answer %q in output: %q", tt.wantOutput, text)
			}
			if tt.signal == syscall.SIGTERM && strings.Contains(text, "Synthetic confirmation? Yes") {
				t.Errorf("interrupted confirmation printed an affirmative answer: %q", text)
			}
		})
	}
}
