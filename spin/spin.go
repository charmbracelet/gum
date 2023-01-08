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
	"github.com/charmbracelet/gum/timeout"
	"os/exec"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	spinner spinner.Model
	title   string
	align   string
	command []string
	aborted bool

	status     int
	stdout     string
	stderr     string
	timeout    time.Duration
	hasTimeout bool
}

type finishCommandMsg struct {
	stdout string
	stderr string
	status int
}

func commandStart(command []string) tea.Cmd {
	return func() tea.Msg {
		var args []string
		if len(command) > 1 {
			args = command[1:]
		}
		cmd := exec.Command(command[0], args...) //nolint:gosec

		var outbuf, errbuf strings.Builder
		cmd.Stdout = &outbuf
		cmd.Stderr = &errbuf

		_ = cmd.Run()

		status := cmd.ProcessState.ExitCode()

		if status == -1 {
			status = 1
		}

		return finishCommandMsg{
			stdout: outbuf.String(),
			stderr: errbuf.String(),
			status: status,
		}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		commandStart(m.command),
		timeout.Init(m.timeout, nil),
	)
}
func (m model) View() string {
	var str string
	if m.hasTimeout {
		str = timeout.TimeoutStr(m.timeout)
	}

	if m.align == "left" {
		return m.spinner.View() + " " + m.title + str
	}

	return str + " " + m.title + " " + m.spinner.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case timeout.TickTimeoutMsg:
		if msg.TimeoutValue <= 0 {
			m.status = 130
			return m, tea.Quit
		}
		m.timeout = msg.TimeoutValue
		return m, timeout.Tick(msg.TimeoutValue, msg.Data)
	case finishCommandMsg:
		m.stdout = msg.stdout
		m.stderr = msg.stderr
		m.status = msg.status
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
