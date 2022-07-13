package choose

import (
	"errors"
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
		if input == "" {
			return errors.New("no options provided, see `gum choose --help`")
		}
		o.Options = strings.Split(strings.TrimSpace(input), "\n")
	}

	var items []item
	for _, option := range o.Options {
		items = append(items, item{text: option, selected: false})
	}

	// We don't need to display prefixes if we are only picking one option.
	// Simply displaying the cursor is enough.
	if o.Limit == 1 && !o.NoLimit {
		o.SelectedPrefix = ""
		o.UnselectedPrefix = ""
		o.CursorPrefix = ""
	}

	// If we've set no limit then we can simply select as many options as there
	// are so let's set the limit to the number of options.
	if o.NoLimit {
		o.Limit = len(o.Options)
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
