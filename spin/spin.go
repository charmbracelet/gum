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
	"syscall"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
)

type model struct {
	spinner    spinner.Model
	title      string
	align      string
	command    []string
	quitting   bool
	status     int
	stdout     string
	stderr     string
	output     string
	showOutput bool
	showError  bool
}

var (
	bothbuf strings.Builder
	outbuf  strings.Builder
	errbuf  strings.Builder

	executing *exec.Cmd
)

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

		executing = exec.Command(command[0], args...) //nolint:gosec
		if term.IsTerminal(os.Stdout.Fd()) {
			executing.Stdout = io.MultiWriter(&bothbuf, &outbuf)
			executing.Stderr = io.MultiWriter(&bothbuf, &errbuf)
		} else {
			executing.Stdout = os.Stdout
			executing.Stderr = os.Stderr
		}
		executing.Stdin = os.Stdin
		_ = executing.Run()
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

func (m model) View() string {
	if m.quitting && m.showOutput {
		return strings.TrimPrefix(errbuf.String()+"\n"+outbuf.String(), "\n")
	}

	var header string
	if m.align == "left" {
		header = m.spinner.View() + " " + m.title
	} else {
		header = m.title + " " + m.spinner.View()
	}
	if !m.showOutput {
		return header
	}
	return header + errbuf.String() + "\n" + outbuf.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
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
			return m, commandAbort
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}
