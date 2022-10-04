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
	"bufio"
	"bytes"
	"os/exec"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	spinner spinner.Model
	title   string
	align   string
	command []string
	aborted bool

	status int
	stdout string
	stderr string

	stdin         chan string
	readFromStdin bool
}

type finishCommandMsg struct {
	stdout string
	stderr string
	status int
}

func (m model) commandStart(command []string) tea.Cmd {
	return func() tea.Msg {
		var args []string
		if len(command) > 1 {
			args = command[1:]
		}
		cmd := exec.Command(command[0], args...) //nolint:gosec

		// Create stdout, stderr buffers to store final output
		var outbuf, errbuf strings.Builder

		// Create stdout, stderr pipes to read data asynchronously
		cmdStdoutReader, _ := cmd.StdoutPipe()
		cmdStderrReader, _ := cmd.StderrPipe()
		scannerStdout := bufio.NewScanner(cmdStdoutReader)
		// Setting split function to capture '\n'
		scannerStdout.Split(splitWithNewLine)
		scannerStderr := bufio.NewScanner(cmdStderrReader)
		scannerStderr.Split(splitWithNewLine)
		var wg sync.WaitGroup
		// Settings wg = 2 to store data from stdin and stderr of executable file
		wg.Add(2)
		go readStdin(&wg, scannerStdout, m.stdin, &outbuf, m.readFromStdin)
		go readStderr(&wg, scannerStderr, &errbuf)
		_ = cmd.Start()
		wg.Wait()
		_ = cmd.Wait()

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
		m.commandStart(m.command),
	)
}
func (m model) View() string {
	if m.align == "left" {
		return m.spinner.View() + " " + m.title
	}

	return m.title + " " + m.spinner.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
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

	select {
	case title := <-m.stdin:
		m.title = title
	default:
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func readStdin(wg *sync.WaitGroup, scanner *bufio.Scanner, stdin chan<- string, outBuf *strings.Builder, updateTitle bool) {
	defer wg.Done()
	for scanner.Scan() {
		text := scanner.Text()
		if updateTitle {
			stdin <- text
		}
		outBuf.WriteString(text)
	}
}

func readStderr(wg *sync.WaitGroup, scanner *bufio.Scanner, errBuf *strings.Builder) {
	defer wg.Done()
	for scanner.Scan() {
		text := scanner.Text()
		errBuf.WriteString(text)
	}
}

func splitWithNewLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
