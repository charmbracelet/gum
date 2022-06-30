package input

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/coral"
	"github.com/muesli/termenv"
)

type model struct{ textinput textinput.Model }

func (m model) Init() tea.Cmd { return nil }
func (m model) View() string  { return m.textinput.View() }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "enter":
			return m, tea.Quit
		}
	}

	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

type options struct {
	prompt      *string
	placeholder *string
	width       *int
}

// Cmd returns the input command
func Cmd() *coral.Command {
	var opts options

	var cmd = &coral.Command{
		Use:   "input",
		Short: "Input prompts the user for input.",
		RunE: func(cmd *coral.Command, args []string) error {
			ti := textinput.New()

			// Flags + Options
			ti.Prompt = *opts.prompt
			ti.Placeholder = *opts.placeholder
			ti.Width = *opts.width

			ti.Focus()

			p := tea.NewProgram(model{ti}, tea.WithOutput(os.Stderr))
			m, err := p.StartReturningModel()

			fmt.Println(m.(model).textinput.Value())

			return err
		},
	}

	opts = options{
		prompt:      cmd.Flags().String("prompt", "> ", "Prompt to display"),
		placeholder: cmd.Flags().String("placeholder", "Enter a value...", "Placeholder value"),
		width:       cmd.Flags().Int("width", 20, "Input width"),
	}

	return cmd
}
