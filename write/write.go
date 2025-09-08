// Package write provides a shell script interface for the text area bubble.
// https://github.com/charmbracelet/bubbles/tree/master/textarea
//
// It can be used to ask the user to write some long form of text (multi-line)
// input. The text the user entered will be sent to stdout.
// Text entry is completed with CTRL+D and aborted with CTRL+C or Escape.
//
// $ gum write > output.text
package write

import (
	"io"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/editor"
)

type keymap struct {
	textarea.KeyMap
	Submit       key.Binding
	Quit         key.Binding
	Abort        key.Binding
	OpenInEditor key.Binding
}

// FullHelp implements help.KeyMap.
func (k keymap) FullHelp() [][]key.Binding { return nil }

// ShortHelp implements help.KeyMap.
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.InsertNewline,
		k.OpenInEditor,
		k.Submit,
	}
}

func defaultKeymap() keymap {
	km := textarea.DefaultKeyMap
	km.InsertNewline = key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("ctrl+j", "insert newline"),
	)
	return keymap{
		KeyMap: km,
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "cancel"),
		),
		OpenInEditor: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "open editor"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),
	}
}

type model struct {
	autoWidth   bool
	header      string
	headerStyle lipgloss.Style
	quitting    bool
	submitted   bool
	textarea    textarea.Model
	showHelp    bool
	help        help.Model
	keymap      keymap
	padding     []int
}

func (m model) Init() tea.Cmd { return textarea.Blink }

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var parts []string

	// Display the header above the text area if it is not empty.
	if m.header != "" {
		parts = append(parts, m.headerStyle.Render(m.header))
	}
	parts = append(parts, m.textarea.View())
	if m.showHelp {
		parts = append(parts, "", m.help.View(m.keymap))
	}
	return lipgloss.NewStyle().
		Padding(m.padding...).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			parts...,
		))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.autoWidth {
			m.textarea.SetWidth(msg.Width - m.padding[1] - m.padding[3])
		}
	case tea.FocusMsg, tea.BlurMsg:
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	case startEditorMsg:
		return m, openEditor(msg.path, msg.lineno)
	case editorFinishedMsg:
		if msg.err != nil {
			m.quitting = true
			return m, tea.Interrupt
		}
		m.textarea.SetValue(msg.content)
	case tea.KeyMsg:
		km := m.keymap
		switch {
		case key.Matches(msg, km.Abort):
			m.quitting = true
			return m, tea.Interrupt
		case key.Matches(msg, km.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, km.Submit):
			m.quitting = true
			m.submitted = true
			return m, tea.Quit
		case key.Matches(msg, km.OpenInEditor):
			//nolint: gosec
			return m, createTempFile(m.textarea.Value(), uint(m.textarea.Line())+1)
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

type startEditorMsg struct {
	path   string
	lineno uint
}

type editorFinishedMsg struct {
	content string
	err     error
}

func createTempFile(content string, lineno uint) tea.Cmd {
	return func() tea.Msg {
		f, err := os.CreateTemp("", "gum.*.md")
		if err != nil {
			return editorFinishedMsg{err: err}
		}
		_, err = io.WriteString(f, content)
		if err != nil {
			return editorFinishedMsg{err: err}
		}
		_ = f.Close()
		return startEditorMsg{
			path:   f.Name(),
			lineno: lineno,
		}
	}
}

func openEditor(path string, lineno uint) tea.Cmd {
	cb := func(err error) tea.Msg {
		if err != nil {
			return editorFinishedMsg{
				err: err,
			}
		}
		bts, err := os.ReadFile(path)
		if err != nil {
			return editorFinishedMsg{err: err}
		}
		return editorFinishedMsg{
			content: string(bts),
		}
	}
	cmd, err := editor.Cmd(
		"Gum",
		path,
		editor.LineNumber(lineno),
		editor.EndOfLine(),
	)
	if err != nil {
		return func() tea.Msg { return cb(err) }
	}
	return tea.ExecProcess(cmd, cb)
}
