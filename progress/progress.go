package progress

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	/// where to read the content from (stdin)
	reader io.Reader
	/// should the content read from stdin be printed out
	output bool
	/// offset in model.buff where to write to next
	offset int
	/// buffer for storing content read from stdin
	buff [1024]byte
	/// tells if the program's output (stdout) is a tty
	isTTY bool

	/// what string in the content indicates progress
	progressIndicator string
	/// should the progress indicator be printed out when model.output is set to true
	hideProgressIndicator bool

	/// stores metadata of the progress
	binfo *barInfo
	/// renderer for the progress depending on what format was specified
	bfmt *barFormatter

	/// stores the width of the terminal received via tea.Msg
	width int

	/// stores the error that occured so that it can later be communicated
	/// to the user
	err error
}

type progressMsg uint
type finishedMsg struct{}
type tickMsg struct{}

func (m *model) readUntilProgressOrEOF() tea.Cmd {
	return func() tea.Msg {
		amountRead, err := m.reader.Read(m.buff[m.offset:])
		if err != nil {
			if err == io.EOF {
				return finishedMsg{}
			}
			m.err = fmt.Errorf("failed to read the input: %w\n", err)
			return tea.Quit
		}

		read := m.buff[m.offset : m.offset+amountRead]
		progress := bytes.Count(read, []byte(m.progressIndicator))
		if m.hideProgressIndicator {
			read = bytes.ReplaceAll(read, []byte(m.progressIndicator), []byte{})
			copy(m.buff[m.offset:], read)
		}
		m.offset += len(read)
		return progressMsg(progress)
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(m.readUntilProgressOrEOF(), tick())
}

func (m *model) View() string {
	padding := 2
	rendered := m.bfmt.Render(m.binfo, max(0, m.width-(padding*2)))

	return lipgloss.NewStyle().
		PaddingLeft(padding).
		PaddingRight(padding).
		MaxWidth(m.width).
		Width(m.width).
		Render(rendered)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tickMsg:
		return m, tick()
	case finishedMsg:
		if m.output && m.offset > 0 && m.isTTY {
			return m, tea.Batch(tea.Println(string(m.buff[:m.offset])), tea.Quit)
		}
		return m, tea.Quit
	case progressMsg:
		var cmd tea.Cmd
		switch {
		// given that we have something to print, and stdout is a tty the
		// user sees both this model and the output in a terminal. Therefore
		// we should/can use tea.Println so that rendering is correct (here meaining
		// that the model and the output don't become intermingled). The output
		// will go to stderr as configured in command.go but that does not matter.
		case m.output && m.offset > 0 && m.isTTY:
			// tea.Println always adds a new-line at the end so we can only print
			// full lines otherwise the added \n cut's the string that we want to print
			end := bytes.LastIndexByte(m.buff[:m.offset], '\n')
			if end < 0 {
				cmd = m.readUntilProgressOrEOF()
				break
			}

			start := end
			if m.buff[max(0, end-1)] == '\r' {
				start = end - 1
			}

			toPrint := m.buff[:start]
			remaining := m.buff[end+1 : m.offset]
			cmd = tea.Batch(m.readUntilProgressOrEOF(), tea.Println(string(toPrint)))
			copy(m.buff[:], remaining)
			m.offset = len(remaining)
		// we have some output that we want to print but it is most likely
		// being piped to another process thus we can just print
		// the message to stdout for the other process to read.
		case m.output && m.offset > 0:
			os.Stdout.Write(m.buff[:m.offset])
			m.offset = 0
			cmd = m.readUntilProgressOrEOF()
		// nothing to print here
		default:
			cmd = m.readUntilProgressOrEOF()
			m.offset = 0
		}
		m.binfo.Update(uint(msg))
		return m, cmd
	}
	return m, nil
}
