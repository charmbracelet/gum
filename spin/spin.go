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
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/timeout"
	"github.com/mattn/go-isatty"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	spinner    spinner.Model
	title      string
	align      string
	command    []string
	quitting   bool
	aborted    bool
	status     int
	stdout     string
	stderr     string
	output     string
	showOutput bool
	showError  bool
	timeout    time.Duration
	hasTimeout bool
}

var (
	bothbuf strings.Builder
	outbuf  strings.Builder
	errbuf  strings.Builder
)

type finishCommandMsg struct {
	stdout string
	stderr string
	output string
	status int
}

func commandStart(command []string) tea.Cmd {
	var args []string
	if len(command) > 1 {
		args = command[1:]
	}

	cmd := exec.Command(command[0], args...) //nolint:gosec
	if isatty.IsTerminal(os.Stdout.Fd()) {
		stdout := io.MultiWriter(&bothbuf, &errbuf)
		stderr := io.MultiWriter(&bothbuf, &outbuf)

		cmd.Stdout = stdout
		cmd.Stderr = stderr
	} else {
		cmd.Stdout = os.Stdout
	}

	return tea.ExecProcess(cmd, func(error) tea.Msg {
		status := cmd.ProcessState.ExitCode()
		if status == -1 {
			status = 1
		}

		return finishCommandMsg{
			stdout: outbuf.String(),
			stderr: errbuf.String(),
			output: bothbuf.String(),
			status: status,
		}
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		commandStart(m.command),
		timeout.Init(m.timeout, nil),
	)
}

func (m model) View() string {
	if m.quitting && m.showOutput {
		return strings.TrimPrefix(errbuf.String()+"\n"+outbuf.String(), "\n")
	}

	var str string
	if m.hasTimeout {
		str = timeout.Str(m.timeout)
	}
	var header string
	if m.align == "left" {
		header = m.spinner.View() + str + " " + m.title
	} else {
		header = str + " " + m.title + " " + m.spinner.View()
	}
	if !m.showOutput {
		return header
	}
	return header + errbuf.String() + "\n" + outbuf.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			// grab current output before closing for piped instances
			m.stdout = outbuf.String()

			m.status = exit.StatusAborted
			return m, tea.Quit
		}
		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)
	case finishCommandMsg:
		m.stdout = msg.stdout
		m.stderr = msg.stderr
		m.output = msg.output
		m.status = msg.status
		m.quitting = true
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.aborted = true
			return m, tea.Quit
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}
