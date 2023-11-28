package date

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fxtlabs/date"
)

type interval int

// Default styles.
var (
	weekdayStyle           = lipgloss.NewStyle().Faint(true)
	defaultCursorTextStyle = lipgloss.NewStyle().Underline(true).Foreground(lipgloss.Color("201"))
)

const (
	day   interval = 0
	month interval = 1
	year  interval = 2
)

// incr i in direction d; bodge mod-3 indexing.
func (i interval) incr(d direction) interval {
	mod := (int(i) + int(d)) % 3
	if mod < 0 {
		return year
	}
	return interval(mod)
}

type direction int

const (
	forward  direction = 1
	backward direction = -1
)

// picker implements tea.Model for a date.Date.
type picker struct {
	date.Date
	focus interval

	promptStyle lipgloss.Style
	prompt      string

	cursorTextStyle lipgloss.Style
}

func basePicker() *picker {
	return &picker{
		Date:            date.Today(),
		focus:           day,
		prompt:          "> ",
		cursorTextStyle: defaultCursorTextStyle,
	}
}

func (p *picker) formatDate() string {
	raw := p.Date.Format("02 Jan 2006")
	parts := strings.Split(raw, " ")
	parts[int(p.focus)] = p.cursorTextStyle.Render(parts[int(p.focus)])
	return strings.Join(parts, " ") + " " + p.formatWeekday()
}

func (p *picker) formatWeekday() string {
	name := ""
	switch p.Date.Weekday() {
	case time.Monday:
		name = "Mon"
	case time.Tuesday:
		name = "Tue"
	case time.Wednesday:
		name = "Wed"
	case time.Thursday:
		name = "Thu"
	case time.Friday:
		name = "Fri"
	case time.Saturday:
		name = "Sat"
	case time.Sunday:
		name = "Sun"
	}
	return weekdayStyle.Render(name)
}

func (p *picker) incr(d direction) {
	switch p.focus {
	case day:
		p.Date = p.Date.AddDate(0, 0, int(d))
	case month:
		p.Date = p.Date.AddDate(0, int(d), 0)
	case year:
		p.Date = p.Date.AddDate(int(d), 0, 0)
	}
}

// Init implements tea.Model.
func (p *picker) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (p *picker) Update(msg tea.Msg) (*picker, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// "up"/"down" increment/decrement the focused component, respectively
		case "up", "k":
			p.incr(forward)
		case "down", "j":
			p.incr(backward)

		// "left"/"right" cycle the focused component
		case "left", "h":
			p.focus = p.focus.incr(backward)
		case "right", "l":
			p.focus = p.focus.incr(forward)
		}
	}
	return p, nil
}

// View implements tea.Model.
func (p *picker) View() string {
	return p.promptStyle.Render(p.prompt) + p.formatDate()
}

// Value of p.
func (p *picker) Value() date.Date {
	return p.Date
}
