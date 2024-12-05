package choose

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/stdin"
)

// Run provides a shell script interface for choosing between different through
// options.
func (o Options) Run() error {
	var (
		subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
		verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
	)

	if len(o.Options) <= 0 {
		input, _ := stdin.Read()
		if input == "" {
			return errors.New("no options provided, see `gum choose --help`")
		}
		o.Options = strings.Split(input, "\n")
	}

	if o.SelectIfOne && len(o.Options) == 1 {
		fmt.Println(o.Options[0])
		return nil
	}

	// We don't need to display prefixes if we are only picking one option.
	// Simply displaying the cursor is enough.
	if o.Limit == 1 && !o.NoLimit {
		o.SelectedPrefix = ""
		o.UnselectedPrefix = ""
		o.CursorPrefix = ""
	}

	if o.NoLimit {
		o.Limit = len(o.Options) + 1
	}

	if o.Ordered {
		slices.SortFunc(o.Options, strings.Compare)
	}

	// Keep track of the selected items.
	currentSelected := 0
	// Check if selected items should be used.
	hasSelectedItems := len(o.Selected) > 0
	startingIndex := 0
	currentOrder := 0
	items := make([]item, len(o.Options))
	for i, option := range o.Options {
		var order int
		// Check if the option should be selected.
		isSelected := hasSelectedItems && currentSelected < o.Limit && slices.Contains(o.Selected, option)
		// If the option is selected then increment the current selected count.
		if isSelected {
			if o.Limit == 1 {
				// When the user can choose only one option don't select the option but
				// start with the cursor hovering over it.
				startingIndex = i
				isSelected = false
			} else {
				currentSelected++
				order = currentOrder
				currentOrder++
			}
		}
		items[i] = item{text: option, selected: isSelected, order: order}
	}

	// Use the pagination model to display the current and total number of
	// pages.
	pager := paginator.New()
	pager.SetTotalPages((len(items) + o.Height - 1) / o.Height)
	pager.PerPage = o.Height
	pager.Type = paginator.Dots
	pager.ActiveDot = subduedStyle.Render("•")
	pager.InactiveDot = verySubduedStyle.Render("•")
	pager.KeyMap = paginator.KeyMap{}
	pager.Page = startingIndex / o.Height
	// Disable Keybindings since we will control it ourselves.
	tm, err := tea.NewProgram(model{
		index:             startingIndex,
		currentOrder:      currentOrder,
		height:            o.Height,
		cursor:            o.Cursor,
		header:            o.Header,
		selectedPrefix:    o.SelectedPrefix,
		unselectedPrefix:  o.UnselectedPrefix,
		cursorPrefix:      o.CursorPrefix,
		items:             items,
		limit:             o.Limit,
		paginator:         pager,
		cursorStyle:       o.CursorStyle.ToLipgloss(),
		headerStyle:       o.HeaderStyle.ToLipgloss(),
		itemStyle:         o.ItemStyle.ToLipgloss(),
		selectedItemStyle: o.SelectedItemStyle.ToLipgloss(),
		numSelected:       currentSelected,
		hasTimeout:        o.Timeout > 0,
		timeout:           o.Timeout,
	}, tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		return fmt.Errorf("failed to start tea program: %w", err)
	}
	m := tm.(model)
	if m.aborted {
		return exit.ErrAborted
	}
	if m.timedOut {
		return exit.ErrTimeout
	}
	if o.Ordered && o.Limit > 1 {
		sort.Slice(m.items, func(i, j int) bool {
			return m.items[i].order < m.items[j].order
		})
	}
	var s strings.Builder
	for _, item := range m.items {
		if item.selected {
			s.WriteString(item.text)
			s.WriteRune('\n')
		}
	}

	if term.IsTerminal(os.Stdout.Fd()) {
		fmt.Print(s.String())
	} else {
		fmt.Print(ansi.Strip(s.String()))
	}

	return nil
}
