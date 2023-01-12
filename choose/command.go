package choose

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/style"
)

var (
	subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
	verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
)

// Run provides a shell script interface for choosing between different through
// options.
func (o Options) Run() error {
	if len(o.Options) == 0 {
		input, _ := stdin.Read()
		if input == "" {
			return errors.New("no options provided, see `gum choose --help`")
		}
		o.Options = strings.Split(strings.TrimSuffix(input, "\n"), "\n")
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

	if len(o.Selected) > o.Limit {
		return errors.New("number of selected options cannot be greater than the limit")
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
		isSelected := hasSelectedItems && currentSelected < o.Limit && arrayContains(o.Selected, option)
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

	// Disable Keybindings since we will control it ourselves.
	pager.UseHLKeys = false
	pager.UseLeftRightKeys = false
	pager.UseJKKeys = false
	pager.UsePgUpPgDownKeys = false

	tm, err := tea.NewProgram(model{
		index:             startingIndex,
		currentOrder:      currentOrder,
		height:            o.Height,
		cursor:            o.Cursor,
		selectedPrefix:    o.SelectedPrefix,
		unselectedPrefix:  o.UnselectedPrefix,
		cursorPrefix:      o.CursorPrefix,
		items:             items,
		limit:             o.Limit,
		paginator:         pager,
		cursorStyle:       o.CursorStyle.ToLipgloss(),
		itemStyle:         o.ItemStyle.ToLipgloss(),
		selectedItemStyle: o.SelectedItemStyle.ToLipgloss(),
		numSelected:       currentSelected,
	}, tea.WithOutput(os.Stderr)).Run()

	if err != nil {
		return fmt.Errorf("failed to start tea program: %w", err)
	}

	m := tm.(model)
	if m.aborted {
		return exit.ErrAborted
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

	fmt.Print(s.String())

	return nil
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	style.HideFlags(ctx)
	return nil
}

// Check if an array contains a value.
func arrayContains(strArray []string, value string) bool {
	for _, str := range strArray {
		if str == value {
			return true
		}
	}
	return false
}
