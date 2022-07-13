package choose

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for choosing between different through
// options.
func (o Options) Run() error {
	if len(o.Options) == 0 {
		input, _ := stdin.Read()
		o.Options = strings.Split(strings.TrimSpace(input), "\n")
	}

	var items []item
	for _, option := range o.Options {
		items = append(items, item{text: option, selected: false})
	}

	// We don't need to display prefixes if we are only picking one option.
	// Simply displaying the cursor is enough.
	if o.Limit == 1 {
		o.SelectedPrefix = ""
		o.UnselectedPrefix = ""
		o.CursorPrefix = ""
	}

	m, err := tea.NewProgram(model{
		height:            o.Height,
		cursor:            o.Cursor,
		selectedPrefix:    o.SelectedPrefix,
		unselectedPrefix:  o.UnselectedPrefix,
		cursorPrefix:      o.CursorPrefix,
		items:             items,
		limit:             o.Limit,
		cursorStyle:       o.CursorStyle.ToLipgloss(),
		itemStyle:         o.ItemStyle.ToLipgloss(),
		selectedItemStyle: o.SelectedItemStyle.ToLipgloss(),
	}, tea.WithOutput(os.Stderr)).StartReturningModel()

	var s strings.Builder

	for _, item := range m.(model).items {
		if item.selected {
			s.WriteString(item.text)
			s.WriteRune('\n')
		}
	}

	fmt.Println(strings.TrimSuffix(s.String(), "\n"))

	return err
}
