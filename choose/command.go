package choose

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/mattn/go-runewidth"
)

// Run provides a shell script interface for choosing between different through
// options.
func (o Options) Run() error {
	if len(o.Options) == 0 {
		input, _ := stdin.Read()
		o.Options = strings.Split(input, "\n")
	}

	items := []list.Item{}
	for _, option := range o.Options {
		if option == "" {
			continue
		}
		items = append(items, item(option))
	}

	const defaultWidth = 20

	id := itemDelegate{
		indicator:         o.Indicator,
		indicatorStyle:    o.IndicatorStyle.ToLipgloss(),
		selectedItemStyle: o.SelectedStyle.ToLipgloss(),
		itemStyle:         o.ItemStyle.ToLipgloss().MarginLeft(runewidth.StringWidth(o.Indicator)),
	}

	l := list.New(items, id, defaultWidth, o.Height)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(!o.HidePagination)

	m, err := tea.NewProgram(model{list: l}, tea.WithOutput(os.Stderr)).StartReturningModel()
	fmt.Println(m.(model).choice)
	return err
}
