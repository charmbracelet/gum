package spin

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/coral"
)

type model struct {
	spinner spinner.Model
	title   string
	command string
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}
func (m model) View() string { return m.spinner.View() + " " + m.title }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "enter":
			return m, tea.Quit
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

type options struct {
	frames *string
	title  *string
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
			s.Spinner = spinners[*opts.frames]
			m := model{
				spinner: s,
				title:   *opts.title,
				command: strings.Join(args, " "),
			}
			p := tea.NewProgram(m)
			_, err := p.StartReturningModel()
			return err
		},
	}

	f := cmd.Flags()

	opts = options{
		frames: f.StringP("frames", "f", "dot", "The type of spinner to use"),
		title:  f.StringP("title", "t", "Loading", "The title of the action being performed, shown to the right of the spinner"),
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
