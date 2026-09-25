// Package spin provides a shell script interface for the spinner bubble.
// https://github.com/charmbracelet/bubbles/tree/master/spinner
//
// It is useful for displaying that some task is running in the background
// while consuming it's output so that it is not shown to the user.
//
// For example, let's do a long running task: $ sleep 5
//
// We can simply prepend a spinner to this task to show it to the user, while
// performing the task / command in the background.
//
// $ gum spin -t "Taking a nap..." -- sleep 5
//
// The spinner will automatically exit when the task is complete.
package spin

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/charmbracelet/x/xpty"
)

type model struct {
	spinner    spinner.Model
	title      string
	padding    []int
	align      string
	command    []string
	quitting   bool
	isTTY      bool
	status     int
	stdout     string
	stderr     string
	output     string
	showStdout bool
	showStderr bool
	showError  bool
	err        error
}

var (
	bothbuf bytes.Buffer
	outbuf  bytes.Buffer
	errbuf  bytes.Buffer

	executing *exec.Cmd
)

type errorMsg error

type finishCommandMsg struct {
	stdout string
	stderr string
	output string
	status int
}

func commandStart(command []string) tea.Cmd {
	return func() tea.Msg {
		var args []string
		if len(command) > 1 {
			args = command[1:]
		}

		executing = exec.CommandContext(context.Background(), command[0], args...) //nolint:gosec
		executing.Stdin = os.Stdin

		isTerminal := term.IsTerminal(os.Stdout.Fd())

		// NOTE(@andreynering): We had issues with Git Bash on Windows
		// when it comes to handling PTYs, so we're falling back to
		// to redirecting stdout/stderr as usual to avoid issues.
		//nolint:nestif
		if isTerminal && runtime.GOOS == "windows" {
			executing.Stdout = io.MultiWriter(&bothbuf, &outbuf)
			executing.Stderr = io.MultiWriter(&bothbuf, &errbuf)
			_ = executing.Run()
		} else if isTerminal {
			stdoutPty, err := openPty(os.Stdout)
			if err != nil {
				return errorMsg(err)
			}

			stderrPty, err := openPty(os.Stderr)
			if err != nil {
				_ = stdoutPty.Close()
				return errorMsg(err)
			}

			if outUnixPty, isOutUnixPty := stdoutPty.(*xpty.UnixPty); isOutUnixPty {
				executing.Stdout = outUnixPty.Slave()
			}
			if errUnixPty, isErrUnixPty := stderrPty.(*xpty.UnixPty); isErrUnixPty {
				executing.Stderr = errUnixPty.Slave()
			}

			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				_, _ = io.Copy(io.MultiWriter(&bothbuf, &outbuf), stdoutPty)
			}()
			go func() {
				defer wg.Done()
				_, _ = io.Copy(io.MultiWriter(&bothbuf, &errbuf), stderrPty)
			}()

			if err = stdoutPty.Start(executing); err != nil {
				_ = stdoutPty.Close()
				_ = stderrPty.Close()
				return errorMsg(err)
			}
			_ = xpty.WaitProcess(context.Background(), executing)

			// Close the PTYs so the copy goroutines hit EOF, then wait
			// for them to drain before we read from the buffers. Without
			// this, short-lived commands (e.g. cp that fails with
			// "Permission denied") can race the goroutines and we end up
			// returning empty stdout/stderr.
			_ = stdoutPty.Close()
			_ = stderrPty.Close()
			wg.Wait()
		} else {
			executing.Stdout = os.Stdout
			executing.Stderr = os.Stderr
			_ = executing.Run()
		}

		status := executing.ProcessState.ExitCode()
		if status == -1 {
			status = 1
		}

		return finishCommandMsg{
			stdout: outbuf.String(),
			stderr: errbuf.String(),
			output: bothbuf.String(),
			status: status,
		}
	}
}

func commandAbort() tea.Msg {
	if executing != nil && executing.Process != nil {
		_ = executing.Process.Signal(syscall.SIGINT)
	}
	return tea.InterruptMsg{}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		commandStart(m.command),
	)
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var out string
	if m.showStderr {
		out += errbuf.String()
	}
	if m.showStdout {
		out += outbuf.String()
	}

	if !m.isTTY {
		return tea.NewView(m.title)
	}

	var header string
	if m.align == "left" {
		header = m.spinner.View() + " " + m.title
	} else {
		header = m.title + " " + m.spinner.View()
	}
	return tea.NewView(lipgloss.NewStyle().
		Padding(m.padding...).
		Render(header, "", out))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case finishCommandMsg:
		m.stdout = msg.stdout
		m.stderr = msg.stderr
		m.output = msg.output
		m.status = msg.status
		m.quitting = true
		return m, tea.Quit
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, commandAbort
		}
	case errorMsg:
		m.err = msg
		m.quitting = true
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}
