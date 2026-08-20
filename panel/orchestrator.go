package panel

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/choose"
	"github.com/charmbracelet/gum/filter"
	"github.com/charmbracelet/lipgloss"
	"github.com/sahilm/fuzzy"
)

type helpOverlayMsg string

type orchestrator struct {
	panels      []Panel
	activeIdx   int
	height      int
	borderStyle lipgloss.Style
	showHelp    bool
	customHelp  string
	gap         int
	vertical    bool
	reverse     bool
	stacked     bool
	delimiter   string
	debug       bool
	all         bool
	helpOverlay string

	// Border styles
	activeBorderStyle   lipgloss.Style
	inactiveBorderStyle lipgloss.Style

	// Common styles
	matchStyle        lipgloss.Style
	cursorStyle       lipgloss.Style
	headerStyle       lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	indicatorStyle    lipgloss.Style

	// Choose-specific styles
	selectedPrefixStyle   lipgloss.Style
	unselectedPrefixStyle lipgloss.Style

	// Filter-specific styles
	textStyle        lipgloss.Style
	cursorTextStyle  lipgloss.Style
	promptStyle      lipgloss.Style
	placeholderStyle lipgloss.Style

	// Global options
	single bool

	chooseModels []chooseModel
	filterModels []filterModel
	quitting,
	submitted bool
	errorMessage string
	help         help.Model
	keymap       panelKeymap
}

type chooseModel struct {
	index             int
	currentOrder      int
	height            int
	padding           []int
	cursor            string
	header            string
	execHelp          string
	selectedPrefix    string
	unselectedPrefix  string
	cursorPrefix      string
	items             []chooseItem
	limit             int
	noLimit           bool
	numSelected       int
	paginator         paginator.Model
	showHelp          bool
	keymap            chooseKeymap
	cursorStyle       lipgloss.Style
	headerStyle       lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	quitting          bool
	submitted         bool
}

type chooseItem struct {
	Text     string
	Selected bool
	Order    int
}

type chooseKeymap struct {
	Down,
	Up,
	Right,
	Left,
	Home,
	End,
	ToggleAll,
	Toggle,
	HelpToggle,
	Abort,
	Quit,
	Submit key.Binding
}

type filterModel struct {
	textinput             textinput.Model
	viewport              *viewport.Model
	choices               map[string]string
	filteringChoices      []string
	matches               []fuzzy.Match
	cursor                int
	header                string
	execHelp              string
	selected              map[string]struct{}
	limit                 int
	noLimit               bool
	numSelected           int
	indicator             string
	selectedPrefix        string
	unselectedPrefix      string
	height                int
	padding               []int
	quitting              bool
	headerStyle           lipgloss.Style
	matchStyle            lipgloss.Style
	textStyle             lipgloss.Style
	cursorTextStyle       lipgloss.Style
	indicatorStyle        lipgloss.Style
	selectedPrefixStyle   lipgloss.Style
	unselectedPrefixStyle lipgloss.Style
	reverse               bool
	fuzzy                 bool
	sort                  bool
	strict                bool
	value                 string
	showHelp              bool
	keymap                filterKeymap
	submitted             bool
}

type filterKeymap struct {
	FocusInSearch,
	FocusOutSearch,
	Down,
	Up,
	NDown,
	NUp,
	Home,
	End,
	ToggleAndNext,
	ToggleAndPrevious,
	ToggleAll,
	Toggle,
	HelpToggle,
	Abort,
	Quit,
	Submit key.Binding
}

type panelKeymap struct {
	NextPanel  key.Binding
	PrevPanel  key.Binding
	LeftPanel  key.Binding
	RightPanel key.Binding
	Submit     key.Binding
	SubmitAll  key.Binding // Submit all current selections without auto-selecting
	PanelHelp  key.Binding
	Quit       key.Binding
}

// FullHelp implements help.KeyMap.
func (k panelKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextPanel, k.PrevPanel, k.LeftPanel, k.RightPanel},
		{k.Submit, k.SubmitAll, k.PanelHelp, k.Quit},
	}
}

// ShortHelp implements help.KeyMap.
func (k panelKeymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.NextPanel,
		k.PrevPanel,
		k.Submit,
		k.SubmitAll,
		k.PanelHelp,
		k.Quit,
	}
}

func defaultPanelKeymap() panelKeymap {
	return panelKeymap{
		NextPanel: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next panel"),
		),
		PrevPanel: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev panel"),
		),
		LeftPanel: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "prev panel"),
		),
		RightPanel: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "next panel"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "next/submit"),
		),
		SubmitAll: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "submit all"),
		),
		PanelHelp: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit"),
		),
	}
}

// BorderStyle returns the lipgloss styling for panel borders.
func (o Options) BorderStyle() lipgloss.Style {
	borderStyle := lipgloss.NewStyle()

	switch o.Border {
	case "none":
		borderStyle = borderStyle.BorderStyle(lipgloss.Border{})
	case "single":
		borderStyle = borderStyle.BorderStyle(lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "┌",
			TopRight:    "┐",
			BottomLeft:  "└",
			BottomRight: "┘",
		})
	case "double":
		borderStyle = borderStyle.BorderStyle(lipgloss.Border{
			Top:         "═",
			Bottom:      "═",
			Left:        "║",
			Right:       "║",
			TopLeft:     "╔",
			TopRight:    "╗",
			BottomLeft:  "╚",
			BottomRight: "╝",
		})
	case "rounded":
		borderStyle = borderStyle.BorderStyle(lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "╰",
			BottomRight: "╯",
		})
	}

	return borderStyle
}

// kongVars are the default variables required by choose/filter Options structs.
var kongVars = kong.Vars{
	"defaultHeight":           "0",
	"defaultWidth":            "0",
	"defaultAlign":            "left",
	"defaultBorder":           "none",
	"defaultBorderForeground": "",
	"defaultBorderBackground": "",
	"defaultBackground":       "",
	"defaultForeground":       "",
	"defaultMargin":           "0 0",
	"defaultPadding":          "0 0",
	"defaultUnderline":        "false",
	"defaultBold":             "false",
	"defaultFaint":            "false",
	"defaultItalic":           "false",
	"defaultStrikethrough":    "false",
}

func parseChooseBlock(args []string) (*choose.Options, error) {
	var opts choose.Options
	parser, err := kong.New(&opts, kong.Exit(func(int) {}), kongVars)
	if err != nil {
		return nil, fmt.Errorf("create choose parser: %w", err)
	}
	if _, err = parser.Parse(args); err != nil {
		return nil, fmt.Errorf("invalid choose options: %w", err)
	}
	return &opts, nil
}

func parseFilterBlock(args []string) (*filter.Options, error) {
	var opts filter.Options
	parser, err := kong.New(&opts, kong.Exit(func(int) {}), kongVars)
	if err != nil {
		return nil, fmt.Errorf("create filter parser: %w", err)
	}
	if _, err = parser.Parse(args); err != nil {
		return nil, fmt.Errorf("invalid filter options: %w", err)
	}
	return &opts, nil
}

// parsePanelsFromArgs parses panel blocks separated by "--" tokens.
// Each block starts with "choose" or "filter" followed by flags and items.
// A "--" token is only treated as a panel separator when followed by
// "choose" or "filter"; otherwise it is treated as end-of-flags for kong
// within the current block, allowing items starting with "-" to be used.
//
// Example: ["choose", "--limit", "3", "a", "b", "--", "filter", "x", "y"]
// Example: ["choose", "--no-limit", "--", "-Q", "-C", "pkg"]  (-- inside block)
func parsePanelsFromArgs(args []string, _ string, _ bool) ([]Panel, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no panels specified\n\nExample:\n  gum panel -- choose apple banana -- filter mango papaya")
	}

	blocks := splitPanelBlocks(args)

	var panels []Panel
	for _, block := range blocks {
		if len(block) == 0 {
			continue
		}
		p, err := parseBlock(block)
		if err != nil {
			return nil, err
		}
		panels = append(panels, p)
	}

	if len(panels) == 0 {
		return nil, fmt.Errorf("no panels specified\n\nExample:\n  gum panel -- choose apple banana -- filter mango papaya")
	}

	return panels, nil
}

// splitPanelBlocks splits args at "--" tokens that are followed by a panel
// type keyword ("choose" or "filter"). Other "--" tokens (e.g. end-of-flags
// within a panel block) are kept as-is.
func splitPanelBlocks(args []string) [][]string {
	var blocks [][]string
	var current []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" && i+1 < len(args) && isPanelType(args[i+1]) {
			if len(current) > 0 || len(blocks) == 0 {
				blocks = append(blocks, current)
			}
			current = []string{args[i+1]}
			i++
			continue
		}
		current = append(current, arg)
	}
	if len(current) > 0 || len(blocks) == 0 {
		blocks = append(blocks, current)
	}

	return blocks
}

// isPanelType returns true if s (case-insensitive) is "choose" or "filter".
func isPanelType(s string) bool {
	s = strings.ToLower(s)
	return s == "choose" || s == "filter"
}

// parseBlock parses a single panel block: first token is type, rest are options+items.
func parseBlock(block []string) (Panel, error) {
	typeName := strings.ToLower(block[0])
	rest := block[1:]

	switch PanelType(typeName) {
	case PanelChoose:
		if len(rest) == 0 {
			return Panel{}, fmt.Errorf("choose panel has no items")
		}
		opts, err := parseChooseBlock(rest)
		if err != nil {
			return Panel{}, err
		}
		if len(opts.Options) == 0 {
			return Panel{}, fmt.Errorf("choose panel has no items")
		}
		return Panel{Type: PanelChoose, ChooseOpts: opts}, nil
	case PanelFilter:
		if len(rest) == 0 {
			return Panel{}, fmt.Errorf("filter panel has no items")
		}
		opts, err := parseFilterBlock(rest)
		if err != nil {
			return Panel{}, err
		}
		if len(opts.Options) == 0 {
			return Panel{}, fmt.Errorf("filter panel has no items")
		}
		return Panel{Type: PanelFilter, FilterOpts: opts}, nil
	default:
		return Panel{}, fmt.Errorf("unknown panel type %q: expected 'choose' or 'filter'", block[0])
	}
}

func (m *orchestrator) initModels(o Options) error {
	m.chooseModels = make([]chooseModel, 0)
	m.filterModels = make([]filterModel, 0)

	chooseIdx := 0
	filterIdx := 0

	for i, panel := range m.panels {
		switch panel.Type {
		case PanelChoose:
			co := panel.ChooseOpts
			items := make([]chooseItem, len(co.Options))
			for j, item := range co.Options {
				items[j] = chooseItem{Text: item, Selected: false, Order: j}
			}

			// Handle pre-selected items
			for _, sel := range co.Selected {
				for j := range items {
					if items[j].Text == sel || sel == "*" {
						items[j].Selected = true
					}
				}
			}

			limit := co.Limit
			if co.NoLimit {
				limit = len(items)
			}

			pager := paginator.New()
			pager.SetTotalPages((len(items) + o.Height - 1) / o.Height)
			pager.PerPage = o.Height
			pager.Type = paginator.Dots

			km := defaultChooseKeymap()
			if co.NoLimit || co.Limit > 1 {
				km.Toggle.SetEnabled(true)
			}
			if co.NoLimit {
				km.ToggleAll.SetEnabled(true)
			}

			cm := chooseModel{
				index:             0,
				currentOrder:      0,
				height:            o.Height,
				padding:           []int{0, 0, 0, 0},
				cursor:            co.Cursor,
				header:            co.Header,
				execHelp:          co.ExecHelp,
				selectedPrefix:    co.SelectedPrefix,
				unselectedPrefix:  co.UnselectedPrefix,
				cursorPrefix:      co.CursorPrefix,
				items:             items,
				limit:             limit,
				noLimit:           co.NoLimit,
				numSelected:       0,
				paginator:         pager,
				showHelp:          false,
				keymap:            km,
				cursorStyle:       o.CursorStyle.ToLipgloss(),
				headerStyle:       o.HeaderStyle.ToLipgloss(),
				itemStyle:         o.ItemStyle.ToLipgloss(),
				selectedItemStyle: o.SelectedItemStyle.ToLipgloss(),
				quitting:          false,
				submitted:         false,
			}
			m.chooseModels = append(m.chooseModels, cm)
			m.panels[i].ModelIdx = chooseIdx
			chooseIdx++

		case PanelFilter:
			fo := panel.FilterOpts
			ti := textinput.New()
			ti.Focus()
			ti.Prompt = fo.Prompt
			ti.PromptStyle = o.PromptStyle.ToLipgloss()
			ti.Placeholder = fo.Placeholder
			ti.PlaceholderStyle = o.PlaceholderStyle.ToLipgloss()
			if fo.Value != "" {
				ti.SetValue(fo.Value)
			}

			v := viewport.New(0, o.Height)

			choices := make(map[string]string, len(fo.Options))
			filteringChoices := make([]string, 0, len(fo.Options))
			for _, opt := range fo.Options {
				choices[opt] = opt
				filteringChoices = append(filteringChoices, opt)
			}

			matches := matchAll(filteringChoices)

			limit := fo.Limit
			if fo.NoLimit {
				limit = len(fo.Options)
			}

			fkm := defaultFilterKeymap()
			if fo.NoLimit || fo.Limit > 1 {
				fkm.Toggle.SetEnabled(true)
				fkm.ToggleAndPrevious.SetEnabled(true)
				fkm.ToggleAndNext.SetEnabled(true)
				fkm.ToggleAll.SetEnabled(true)
			}

			// Handle pre-selected items
			preSelected := make(map[string]struct{})
			for _, sel := range fo.Selected {
				for _, c := range fo.Options {
					if c == sel || sel == "*" {
						preSelected[c] = struct{}{}
					}
				}
			}

			fm := filterModel{
				textinput:             ti,
				viewport:              &v,
				choices:               choices,
				filteringChoices:      filteringChoices,
				matches:               matches,
				cursor:                0,
				header:                fo.Header,
				execHelp:              fo.ExecHelp,
				selected:              preSelected,
				limit:                 limit,
				noLimit:               fo.NoLimit,
				numSelected:           len(preSelected),
				indicator:             fo.Indicator,
				selectedPrefix:        fo.SelectedPrefix,
				unselectedPrefix:      fo.UnselectedPrefix,
				height:                o.Height,
				padding:               []int{0, 0, 0, 0},
				quitting:              false,
				headerStyle:           o.HeaderStyle.ToLipgloss(),
				matchStyle:            o.MatchStyle.ToLipgloss(),
				textStyle:             o.TextStyle.ToLipgloss(),
				cursorTextStyle:       o.CursorTextStyle.ToLipgloss(),
				indicatorStyle:        o.IndicatorStyle.ToLipgloss(),
				selectedPrefixStyle:   o.SelectedPrefixStyle.ToLipgloss(),
				unselectedPrefixStyle: o.UnselectedPrefixStyle.ToLipgloss(),
				reverse:               fo.Reverse,
				fuzzy:                 fo.Fuzzy,
				sort:                  fo.FuzzySort,
				strict:                fo.Strict,
				value:                 fo.Value,
				showHelp:              false,
				keymap:                fkm,
				submitted:             false,
			}
			m.filterModels = append(m.filterModels, fm)
			m.panels[i].ModelIdx = filterIdx
			filterIdx++
		}
	}

	return nil
}

func defaultChooseKeymap() chooseKeymap {
	return chooseKeymap{
		Down: key.NewBinding(
			key.WithKeys("down", "j", "ctrl+j", "ctrl+n"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k", "ctrl+k", "ctrl+p"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l", "ctrl+f"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h", "ctrl+b"),
		),
		Home: key.NewBinding(
			key.WithKeys("g", "home"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
		),
		ToggleAll: key.NewBinding(
			key.WithKeys("a", "A", "ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
			key.WithDisabled(),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" ", "tab", "x", "ctrl+@"),
			key.WithHelp("x", "toggle"),
			key.WithDisabled(),
		),
		HelpToggle: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "exec help"),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "abort"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter", "ctrl+q"),
			key.WithHelp("enter", "submit"),
		),
	}
}

func defaultFilterKeymap() filterKeymap {
	return filterKeymap{
		Down: key.NewBinding(
			key.WithKeys("down", "ctrl+j", "ctrl+n"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "ctrl+k", "ctrl+p"),
		),
		NDown: key.NewBinding(
			key.WithKeys("j"),
		),
		NUp: key.NewBinding(
			key.WithKeys("k"),
		),
		Home: key.NewBinding(
			key.WithKeys("g", "home"),
		),
		End: key.NewBinding(
			key.WithKeys("G", "end"),
		),
		ToggleAndNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle"),
			key.WithDisabled(),
		),
		ToggleAndPrevious: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "toggle"),
			key.WithDisabled(),
		),
		Toggle: key.NewBinding(
			key.WithKeys("ctrl+@", "shift+space"),
			key.WithHelp("ctrl+@", "toggle"),
			key.WithDisabled(),
		),
		ToggleAll: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "select all"),
			key.WithDisabled(),
		),
		HelpToggle: key.NewBinding(
			key.WithKeys("ctrl+h"),
			key.WithHelp("ctrl+h", "exec help"),
		),
		FocusInSearch: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "search"),
		),
		FocusOutSearch: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "blur search"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit"),
		),
		Abort: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "abort"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter", "ctrl+q"),
			key.WithHelp("enter", "submit"),
		),
	}
}

func matchAll(options []string) []fuzzy.Match {
	matches := make([]fuzzy.Match, len(options))
	for i, option := range options {
		matches[i] = fuzzy.Match{Str: option}
	}
	return matches
}

func (m orchestrator) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(m.chooseModels)+len(m.filterModels))

	for range m.chooseModels {
		cmds = append(cmds, nil)
	}
	for range m.filterModels {
		cmds = append(cmds, textinput.Blink)
	}

	return tea.Batch(cmds...)
}

func (m orchestrator) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case helpOverlayMsg:
		m.helpOverlay = string(msg)
		return m, nil

	case tea.WindowSizeMsg:
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.helpOverlay != "" {
			m.helpOverlay = ""
			return m, nil
		}

		km := m.keymap

		switch {
		case key.Matches(msg, km.NextPanel), key.Matches(msg, km.RightPanel):
			m.activeIdx = (m.activeIdx + 1) % len(m.panels)
			m.errorMessage = ""
			return m, nil

		case key.Matches(msg, km.PrevPanel), key.Matches(msg, km.LeftPanel):
			m.activeIdx = (m.activeIdx - 1 + len(m.panels)) % len(m.panels)
			m.errorMessage = ""
			return m, nil

		case key.Matches(msg, km.Submit):
			return m.handleSubmit()

		case key.Matches(msg, km.SubmitAll):
			return m.handleSubmitAll()

		case key.Matches(msg, km.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Safety check: activeIdx in range
	if m.activeIdx < 0 || m.activeIdx >= len(m.panels) {
		return m, nil
	}

	panel := m.panels[m.activeIdx]

	switch panel.Type {
	case PanelChoose:
		if panel.ModelIdx < 0 || panel.ModelIdx >= len(m.chooseModels) {
			err := fmt.Errorf("internal error: choose model index %d out of range (have %d models)", panel.ModelIdx, len(m.chooseModels))
			m.quitting = true
			return m, func() tea.Msg { return err }
		}
		cm := &m.chooseModels[panel.ModelIdx]
		return m, m.updateChooseModel(cm, msg)
	case PanelFilter:
		if panel.ModelIdx < 0 || panel.ModelIdx >= len(m.filterModels) {
			err := fmt.Errorf("internal error: filter model index %d out of range (have %d models)", panel.ModelIdx, len(m.filterModels))
			m.quitting = true
			return m, func() tea.Msg { return err }
		}
		fm := &m.filterModels[panel.ModelIdx]
		return m, m.updateFilterModel(fm, msg)
	default:
		err := fmt.Errorf("internal error: unknown panel type: %s", panel.Type)
		m.quitting = true
		return m, func() tea.Msg { return err }
	}
}

func (m *orchestrator) updateChooseModel(cm *chooseModel, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		km := cm.keymap
		switch {
		case key.Matches(msg, km.Down):
			cm.index++
			if cm.index >= len(cm.items) {
				cm.index = 0
				cm.paginator.Page = 0
			}
			_, end := cm.paginator.GetSliceBounds(len(cm.items))
			if cm.index >= end {
				cm.paginator.NextPage()
			}
		case key.Matches(msg, km.Up):
			cm.index--
			if cm.index < 0 {
				cm.index = len(cm.items) - 1
				cm.paginator.Page = cm.paginator.TotalPages - 1
			}
			_, end := cm.paginator.GetSliceBounds(len(cm.items))
			if cm.index >= end {
				cm.paginator.NextPage()
			}
			start, _ := cm.paginator.GetSliceBounds(len(cm.items))
			if cm.index < start {
				cm.paginator.PrevPage()
			}
		case key.Matches(msg, km.Home):
			cm.index = 0
			cm.paginator.Page = 0
		case key.Matches(msg, km.End):
			cm.index = len(cm.items) - 1
			cm.paginator.Page = cm.paginator.TotalPages - 1
		case key.Matches(msg, km.Toggle):
			if cm.items[cm.index].Selected {
				cm.items[cm.index].Selected = false
				cm.numSelected--
			} else if cm.numSelected < cm.limit {
				cm.items[cm.index].Selected = true
				cm.items[cm.index].Order = cm.currentOrder
				cm.numSelected++
				cm.currentOrder++
			}
		case key.Matches(msg, km.ToggleAll):
			if cm.limit <= 1 && !cm.noLimit {
				break
			}
			if cm.numSelected < len(cm.items) && cm.numSelected < cm.limit {
				for i := range cm.items {
					if cm.numSelected >= cm.limit {
						break
					}
					if !cm.items[i].Selected {
						cm.items[i].Selected = true
						cm.items[i].Order = cm.currentOrder
						cm.numSelected++
						cm.currentOrder++
					}
				}
			} else {
				for i := range cm.items {
					cm.items[i].Selected = false
					cm.items[i].Order = 0
				}
				cm.numSelected = 0
				cm.currentOrder = 0
			}
		case key.Matches(msg, km.HelpToggle):
			return m.execHelp(cm.header, cm.execHelp)
		case key.Matches(msg, km.Submit):
			cm.submitted = true
			cm.quitting = true
			if cm.limit <= 1 && cm.numSelected < 1 {
				cm.items[cm.index].Selected = true
			}
		case key.Matches(msg, km.Quit):
			cm.quitting = true
		}
	}

	var cmd tea.Cmd
	cm.paginator, cmd = cm.paginator.Update(msg)
	return cmd
}

func (m *orchestrator) updateFilterModel(fm *filterModel, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		km := fm.keymap
		switch {
		case key.Matches(msg, km.FocusInSearch):
			fm.textinput.Focus()
		case key.Matches(msg, km.FocusOutSearch):
			fm.textinput.Blur()
		case key.Matches(msg, km.Quit):
			fm.quitting = true
		case key.Matches(msg, km.Submit):
			fm.submitted = true
			fm.quitting = true
		case key.Matches(msg, km.HelpToggle):
			return m.execHelp(fm.header, fm.execHelp)
		case key.Matches(msg, km.Down, km.NDown):
			m.filterCursorDown(fm)
		case key.Matches(msg, km.Up, km.NUp):
			m.filterCursorUp(fm)
		case key.Matches(msg, km.Home):
			fm.cursor = 0
		case key.Matches(msg, km.End):
			fm.cursor = len(fm.choices) - 1
		case key.Matches(msg, km.ToggleAndNext):
			if fm.limit == 1 && !fm.noLimit {
				break
			}
			m.toggleFilterSelection(fm)
			m.filterCursorDown(fm)
		case key.Matches(msg, km.ToggleAndPrevious):
			if fm.limit == 1 && !fm.noLimit {
				break
			}
			m.toggleFilterSelection(fm)
			m.filterCursorUp(fm)
		case key.Matches(msg, km.Toggle):
			if fm.limit == 1 && !fm.noLimit {
				break
			}
			m.toggleFilterSelection(fm)
		case key.Matches(msg, km.ToggleAll):
			if fm.limit <= 1 && !fm.noLimit {
				break
			}
			maxSelect := fm.limit
			if fm.numSelected < len(fm.matches) && fm.numSelected < maxSelect {
				for i := range fm.matches {
					if fm.numSelected >= maxSelect {
						break
					}
					if _, ok := fm.selected[fm.matches[i].Str]; !ok {
						fm.selected[fm.matches[i].Str] = struct{}{}
						fm.numSelected++
					}
				}
			} else {
				fm.selected = make(map[string]struct{})
				fm.numSelected = 0
			}
		default:
			var icmd tea.Cmd
			fm.textinput, icmd = fm.textinput.Update(msg)

			newValue := fm.textinput.Value()
			var choices []string
			if !fm.strict {
				choices = append(choices, newValue)
			}
			choices = append(choices, fm.filteringChoices...)
			if fm.fuzzy {
				if fm.sort {
					fm.matches = fuzzy.Find(newValue, choices)
				} else {
					fm.matches = fuzzy.FindNoSort(newValue, choices)
				}
			} else {
				fm.matches = exactMatches(newValue, choices)
			}
			if newValue == "" {
				fm.matches = matchAll(fm.filteringChoices)
			}
			fm.cursor = 0

			return icmd
		}
	}

	fm.keymap.FocusInSearch.SetEnabled(!fm.textinput.Focused())
	fm.keymap.FocusOutSearch.SetEnabled(fm.textinput.Focused())
	fm.keymap.NUp.SetEnabled(!fm.textinput.Focused())
	fm.keymap.NDown.SetEnabled(!fm.textinput.Focused())
	fm.keymap.Home.SetEnabled(!fm.textinput.Focused())
	fm.keymap.End.SetEnabled(!fm.textinput.Focused())

	if len(fm.matches) > 0 {
		fm.cursor = (fm.cursor + len(fm.matches)) % len(fm.matches)
	}

	return nil
}

func (m *orchestrator) filterCursorDown(fm *filterModel) {
	if len(fm.matches) == 0 {
		return
	}
	fm.cursor = (fm.cursor + 1) % len(fm.matches)
}

func (m *orchestrator) filterCursorUp(fm *filterModel) {
	if len(fm.matches) == 0 {
		return
	}
	fm.cursor = (fm.cursor - 1 + len(fm.matches)) % len(fm.matches)
}

func (m *orchestrator) toggleFilterSelection(fm *filterModel) {
	if len(fm.matches) == 0 {
		return
	}
	currentMatch := fm.matches[fm.cursor].Str
	if _, ok := fm.selected[currentMatch]; ok {
		delete(fm.selected, currentMatch)
		fm.numSelected--
	} else if fm.numSelected < fm.limit {
		fm.selected[currentMatch] = struct{}{}
		fm.numSelected++
	}
}

func exactMatches(search string, choices []string) []fuzzy.Match {
	matches := fuzzy.Matches{}
	search = strings.ToLower(search)
	for i, choice := range choices {
		matchedString := strings.ToLower(choice)
		index := strings.Index(matchedString, search)
		if index >= 0 {
			matchedIndexes := make([]int, 0, len(search))
			for s := range search {
				matchedIndexes = append(matchedIndexes, index+s)
			}
			matches = append(matches, fuzzy.Match{
				Str:            choice,
				Index:          i,
				MatchedIndexes: matchedIndexes,
			})
		}
	}
	return matches
}

func (m orchestrator) View() string {
	if m.quitting {
		return ""
	}

	if len(m.panels) == 0 {
		return ""
	}

	// Safety check: activeIdx in range
	if m.activeIdx < 0 || m.activeIdx >= len(m.panels) {
		m.activeIdx = 0
	}

	var panelViews []string

	for i, panel := range m.panels {
		panelBorder := m.borderStyle
		if i != m.activeIdx {
			panelBorder = panelBorder.BorderForeground(lipgloss.Color("240"))
		}

		var view string

		switch panel.Type {
		case PanelChoose:
			if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.chooseModels) {
				cm := m.chooseModels[panel.ModelIdx]
				view = m.renderChooseView(&cm)
			}
		case PanelFilter:
			if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.filterModels) {
				fm := m.filterModels[panel.ModelIdx]
				view = m.renderFilterView(&fm)
			}
		}

		// Use per-panel header if set, otherwise use panel type as title
		var panelHeader string
		switch panel.Type {
		case PanelChoose:
			if panel.ChooseOpts != nil {
				panelHeader = panel.ChooseOpts.Header
			}
		case PanelFilter:
			if panel.FilterOpts != nil {
				panelHeader = panel.FilterOpts.Header
			}
		}
		if panelHeader == "" {
			panelHeader = strings.ToUpper(string(panel.Type))
		}
		headerView := m.headerStyle.Render(" " + panelHeader + " ")
		if i == m.activeIdx {
			headerView = " ● " + headerView
		}

		panelView := panelBorder.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Render(headerView),
				view,
			),
		)
		panelViews = append(panelViews, panelView)
	}

	var joinedView string
	if m.vertical {
		joinedView = lipgloss.JoinVertical(lipgloss.Top, panelViews...)
	} else {
		joinedView = m.joinHorizontalWithGap(panelViews)
	}

	if m.showHelp {
		var helpText string
		if m.errorMessage != "" {
			helpText = " " + m.errorMessage + " "
		} else if m.customHelp != "" {
			helpText = " " + m.customHelp
		} else {
			helpText = " " + m.help.View(m.keymap)
		}
		joinedView = lipgloss.JoinVertical(lipgloss.Left, joinedView, helpText)
	}

	if m.helpOverlay != "" {
		overlayBody := m.renderOverlay(m.helpOverlay)
		if m.reverse {
			joinedView = lipgloss.JoinVertical(lipgloss.Left, overlayBody, joinedView)
		} else {
			joinedView = lipgloss.JoinVertical(lipgloss.Left, joinedView, overlayBody)
		}
	}

	return joinedView
}

func (m orchestrator) renderChooseView(cm *chooseModel) string {
	if cm.quitting {
		return ""
	}

	var s strings.Builder

	start, end := cm.paginator.GetSliceBounds(len(cm.items))
	for i, item := range cm.items[start:end] {
		actualIdx := start + i

		if actualIdx == cm.index%cm.height {
			s.WriteString(cm.cursorStyle.Render(cm.cursor))
		} else {
			s.WriteString(strings.Repeat(" ", lipgloss.Width(cm.cursor)))
		}

		if item.Selected {
			s.WriteString(cm.selectedItemStyle.Render(cm.selectedPrefix + item.Text))
		} else if actualIdx == cm.index%cm.height {
			s.WriteString(cm.cursorStyle.Render(cm.cursorPrefix + item.Text))
		} else {
			s.WriteString(cm.itemStyle.Render(cm.unselectedPrefix + item.Text))
		}
		if i < len(cm.items[start:end])-1 {
			s.WriteRune('\n')
		}
	}

	if cm.paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", cm.height-cm.paginator.ItemsOnPage(len(cm.items))+1))
		s.WriteString("  " + cm.paginator.View())
	}

	return s.String()
}

func (m orchestrator) renderFilterView(fm *filterModel) string {
	if fm.quitting {
		return ""
	}

	var s strings.Builder
	var lineTextStyle lipgloss.Style

	inputView := fm.textinput.View()
	s.WriteString(inputView)
	s.WriteString("\n")

	for i := range fm.matches {
		if i == fm.cursor {
			s.WriteString(fm.indicatorStyle.Render(fm.indicator))
			lineTextStyle = fm.cursorTextStyle
		} else {
			s.WriteString(strings.Repeat(" ", lipgloss.Width(fm.indicator)))
			lineTextStyle = fm.textStyle
		}

		if _, ok := fm.selected[fm.matches[i].Str]; ok {
			s.WriteString(fm.selectedPrefixStyle.Render(fm.selectedPrefix))
		} else if fm.limit > 1 {
			s.WriteString(fm.unselectedPrefixStyle.Render(fm.unselectedPrefix))
		} else {
			s.WriteString(" ")
		}

		styledOption := fm.choices[fm.matches[i].Str]
		if len(fm.matches[i].MatchedIndexes) == 0 {
			s.WriteString(lineTextStyle.Render(styledOption))
			s.WriteRune('\n')
			continue
		}

		s.WriteString(lineTextStyle.Render(styledOption))
		s.WriteRune('\n')
	}

	return s.String()
}

func (m orchestrator) getResults(outputDelimiter string) []string {
	results := make([]string, len(m.panels))

	for i, panel := range m.panels {
		var panelItems []string

		switch panel.Type {
		case PanelChoose:
			if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.chooseModels) {
				cm := m.chooseModels[panel.ModelIdx]
				for _, item := range cm.items {
					if item.Selected {
						panelItems = append(panelItems, item.Text)
					}
				}
			}
		case PanelFilter:
			if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.filterModels) {
				fm := m.filterModels[panel.ModelIdx]
				for k := range fm.selected {
					panelItems = append(panelItems, k)
				}
			}
		}

		results[i] = strings.Join(panelItems, outputDelimiter)
	}

	return results
}

// getSingleResult returns the single selected item (used with --single flag).
func (m orchestrator) getSingleResult() string {
	if m.activeIdx < 0 || m.activeIdx >= len(m.panels) {
		return ""
	}
	panel := m.panels[m.activeIdx]
	switch panel.Type {
	case PanelChoose:
		if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.chooseModels) {
			cm := m.chooseModels[panel.ModelIdx]
			if cm.noLimit {
				var items []string
				for _, item := range cm.items {
					if item.Selected {
						items = append(items, item.Text)
					}
				}
				return strings.Join(items, "\n")
			}
			for _, item := range cm.items {
				if item.Selected {
					return item.Text
				}
			}
		}
	case PanelFilter:
		if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.filterModels) {
			fm := m.filterModels[panel.ModelIdx]
			for k := range fm.selected {
				return k
			}
		}
	}
	return ""
}

func (m *orchestrator) handleSubmit() (orchestrator, tea.Cmd) {
	panel := m.panels[m.activeIdx]

	// Handle --single mode: select current item and quit immediately
	if m.single {
		switch panel.Type {
		case PanelChoose:
			if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.chooseModels) {
				cm := &m.chooseModels[panel.ModelIdx]
				if !cm.noLimit {
					for i := range cm.items {
						cm.items[i].Selected = false
					}
					cm.items[cm.index].Selected = true
					cm.numSelected = 1
				}
				if cm.numSelected == 0 {
					cm.items[cm.index].Selected = true
					cm.numSelected = 1
				}
			}
		case PanelFilter:
			if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.filterModels) {
				fm := &m.filterModels[panel.ModelIdx]
				if len(fm.matches) > 0 {
					fm.selected = make(map[string]struct{})
					fm.selected[fm.matches[fm.cursor].Str] = struct{}{}
					fm.numSelected = 1
				}
			}
		}
		m.submitted = true
		m.quitting = true
		return *m, tea.Quit
	}

	// Default mode: only auto-select when limit == 1 (explicit selection required when noLimit)
	if panel.Type == PanelChoose {
		if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.chooseModels) {
			cm := &m.chooseModels[panel.ModelIdx]
			if cm.limit <= 1 && cm.numSelected < 1 {
				cm.items[cm.index].Selected = true
				cm.numSelected = 1
			}
		}
	} else if panel.Type == PanelFilter {
		if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.filterModels) {
			fm := &m.filterModels[panel.ModelIdx]
			if len(fm.selected) == 0 && fm.limit <= 1 && len(fm.matches) > 0 {
				if fm.strict && len(fm.matches) == 0 {
					m.errorMessage = "No matches found"
					return *m, nil
				}
				fm.selected[fm.matches[fm.cursor].Str] = struct{}{}
			}
		}
	}

	if m.activeIdx < len(m.panels)-1 {
		m.activeIdx++
		m.errorMessage = ""
		return *m, nil
	}

	if m.all {
		firstIncomplete := m.findFirstIncompletePanel()
		if firstIncomplete >= 0 {
			m.activeIdx = firstIncomplete
			m.errorMessage = "You must make a choice in all panels!"
			return *m, nil
		}
	}

	m.submitted = true
	m.quitting = true
	return *m, tea.Quit
}

// handleSubmitAll submits all current selections without auto-selecting.
// This allows users to submit without selecting anything in the last panel.
func (m *orchestrator) handleSubmitAll() (orchestrator, tea.Cmd) {
	// Check --all requirement
	if m.all {
		firstIncomplete := m.findFirstIncompletePanel()
		if firstIncomplete >= 0 {
			m.activeIdx = firstIncomplete
			m.errorMessage = "You must make a choice in all panels!"
			return *m, nil
		}
	}

	m.submitted = true
	m.quitting = true
	return *m, tea.Quit
}

func (m *orchestrator) findFirstIncompletePanel() int {
	for i, panel := range m.panels {
		switch panel.Type {
		case PanelChoose:
			if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.chooseModels) {
				cm := m.chooseModels[panel.ModelIdx]
				if cm.numSelected == 0 {
					return i
				}
			}
		case PanelFilter:
			if panel.ModelIdx >= 0 && panel.ModelIdx < len(m.filterModels) {
				fm := m.filterModels[panel.ModelIdx]
				if len(fm.selected) == 0 {
					return i
				}
			}
		}
	}
	return -1
}

func (m *orchestrator) renderOverlay(content string) string {
	var s strings.Builder
	for _, line := range strings.Split(content, "\n") {
		s.WriteString("│ " + line)
		s.WriteRune('\n')
	}
	s.WriteString("│ (press any key to close)")
	return s.String()
}

func (m *orchestrator) joinHorizontalWithGap(panelViews []string) string {
	if m.gap > 0 && len(panelViews) > 1 {
		gapViews := make([]string, 0, len(panelViews)*2-1)
		for i, pv := range panelViews {
			gapViews = append(gapViews, pv)
			if i < len(panelViews)-1 {
				gapViews = append(gapViews, strings.Repeat(" ", m.gap))
			}
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, gapViews...)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, panelViews...)
}

// execHelp returns a tea.Cmd that runs the given execHelp command.
// If cmdStr is empty it falls back to GUM_PANEL_HELP_<header> env var.
// The command's stdout and stderr are captured and returned as a
// helpOverlayMsg for display inside the TUI overlay.
func (m *orchestrator) execHelp(header, cmdStr string) tea.Cmd {
	if cmdStr == "" {
		cmdStr = os.Getenv("GUM_PANEL_HELP_" + headerToEnvKey(header))
	}
	if cmdStr == "" {
		return nil
	}

	return func() tea.Msg {
		parts := strings.Fields(cmdStr)
		if len(parts) == 0 {
			return nil
		}
		name := parts[0]
		args := parts[1:]
		cmd := exec.Command(name, args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		_ = cmd.Run()
		return helpOverlayMsg(out.String())
	}
}

// headerToEnvKey converts a panel header to an env var key suffix.
// Alphanumeric characters are kept (uppercased); all other characters
// are encoded as _XX where XX is their uppercase hex value, ensuring
// no two distinct headers produce the same key.
func headerToEnvKey(header string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(header) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			fmt.Fprintf(&b, "_%02X", r)
		}
	}
	return b.String()
}
