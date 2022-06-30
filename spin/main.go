// Package spin provides a shell script interface for the spinner bubble.
// https://github.com/charmbracelet/bubbles/spinner
//
// It is useful for displaying that some task is running in the background
// while consuming it's output so that it is not shown to the user.
//
// For example, let's do a long running task:
//   $ sleep 5
//
// We can simply prepend a spinner to this task to show it to the user,
// while performing the task / command in the background.
//
//   $ sea spin -t "Taking a nap..." -- sleep 5
//
// The spinner will automatically exit when the task is complete.
package spin

import (
	"os/exec"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/coral"
)

// model is the state of the spinner program.
// it tracks which spinner to use, the title to show the user, and the command to run in the background.
type model struct {
	spinner spinner.Model
	title   string
	command []string
}

type finishCommandMsg struct{ output string }

func commandStart(command []string) tea.Cmd {
	return func() tea.Msg {
		var args []string
		if len(command) > 1 {
			args = command[1:]
		}
		out, _ := exec.Command(command[0], args...).Output()
		return finishCommandMsg{output: string(out)}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		commandStart(m.command),
	)
}
func (m model) View() string { return m.spinner.View() + " " + m.title }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case finishCommandMsg:
		return m, tea.Quit
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

type options struct {
	spinner *string
	title   *string
	color   *string
}

// Cmd returns the spin command
func Cmd() *coral.Command {
	var opts options

	var cmd = &coral.Command{
		Use:   "spin",
		Short: "Spin displays a spinner and performs some action",
		Args:  coral.ArbitraryArgs,
		RunE: func(cmd *coral.Command, args []string) error {
			s := spinner.New()
			s.Spinner = spinners[*opts.spinner]
			s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(*opts.color))
			m := model{
				spinner: s,
				title:   *opts.title,
				command: args,
			}
			p := tea.NewProgram(m)
			_, err := p.StartReturningModel()
			return err
		},
	}

	f := cmd.Flags()

	opts = options{
		spinner: f.StringP("spinner", "s", "dot", "The type of spinner to use (line, dot, minidot, jump, pulse, points, globe, moon, monkey, meter, hamburger)"),
		title:   f.StringP("title", "t", "Loading", "The title of the action being performed, shown to the right of the spinner"),
		color:   f.StringP("color", "c", "white", "The color of the spinner"),
	}

	return cmd
}

var spinners = map[string]spinner.Spinner{
	"line":      spinner.Line,
	"dot":       spinner.Dot,
	"minidot":   spinner.MiniDot,
	"jump":      spinner.Jump,
	"pulse":     spinner.Pulse,
	"points":    spinner.Points,
	"globe":     spinner.Globe,
	"moon":      spinner.Moon,
	"monkey":    spinner.Monkey,
	"meter":     spinner.Meter,
	"hamburger": spinner.Hamburger,
}
