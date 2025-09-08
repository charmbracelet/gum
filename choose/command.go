package choose

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/internal/tty"
	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/lipgloss"
)

// Run provides a shell script interface for choosing between different through
// options.
func (o Options) Run() error {
	var (
		subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
		verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
	)

	input, _ := stdin.Read(stdin.StripANSI(o.StripANSI))
	if len(o.Options) > 0 && len(o.Selected) == 0 {
		o.Selected = strings.Split(input, o.InputDelimiter)
	} else if len(o.Options) == 0 {
		if input == "" {
			return errors.New("no options provided, see `gum choose --help`")
		}
		o.Options = strings.Split(input, o.InputDelimiter)
	}

	// normalize options into a map
	options := map[string]string{}
	// keep the labels in the user-provided order
	var labels []string //nolint:prealloc
	for _, opt := range o.Options {
		if o.LabelDelimiter == "" {
			options[opt] = opt
			continue
		}
		label, value, ok := strings.Cut(opt, o.LabelDelimiter)
		if !ok {
			return fmt.Errorf("invalid option format: %q", opt)
		}
		labels = append(labels, label)
		options[label] = value
	}
	if o.LabelDelimiter != "" {
		o.Options = labels
	}

	if o.SelectIfOne && len(o.Options) == 1 {
		fmt.Println(options[o.Options[0]])
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

	isSelectAll := len(o.Selected) == 1 && o.Selected[0] == "*"

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
		isSelected := hasSelectedItems && currentSelected < o.Limit && (isSelectAll || slices.Contains(o.Selected, option))
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
	top, right, bottom, left := style.ParsePadding(o.Padding)
	pager := paginator.New()
	pager.SetTotalPages((len(items) + o.Height - 1) / o.Height)
	pager.PerPage = o.Height
	pager.Type = paginator.Dots
	pager.ActiveDot = subduedStyle.Render("•")
	pager.InactiveDot = verySubduedStyle.Render("•")
	pager.KeyMap = paginator.KeyMap{}
	pager.Page = startingIndex / o.Height

	km := defaultKeymap()
	if o.NoLimit || o.Limit > 1 {
		km.Toggle.SetEnabled(true)
	}
	if o.NoLimit {
		km.ToggleAll.SetEnabled(true)
	}

	m := model{
		index:             startingIndex,
		currentOrder:      currentOrder,
		height:            o.Height,
		padding:           []int{top, right, bottom, left},
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
		showHelp:          o.ShowHelp,
		help:              help.New(),
		keymap:            km,
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	// Disable Keybindings since we will control it ourselves.
	tm, err := tea.NewProgram(
		m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return fmt.Errorf("unable to pick selection: %w", err)
	}
	m = tm.(model)
	if !m.submitted {
		return errors.New("nothing selected")
	}
	if o.Ordered && o.Limit > 1 {
		sort.Slice(m.items, func(i, j int) bool {
			return m.items[i].order < m.items[j].order
		})
	}

	var out []string
	for _, item := range m.items {
		if item.selected {
			out = append(out, options[item.text])
		}
	}
	tty.Println(strings.Join(out, o.OutputDelimiter))
	return nil
}
