// Package panel provides a multi-panel TUI for switching between
// choose and filter panels side by side.
//
// Panel blocks are separated by '--'. Each block starts with 'choose' or
// 'filter' and accepts all flags from 'gum choose' / 'gum filter'.
//
// Example:
//
//	$ gum panel -- choose apple banana cherry -- filter mango papaya
//	$ gum panel -- choose --limit 3 --header "Fruit" apple banana \
//	            -- filter --no-fuzzy --placeholder "Search" mango papaya
//
// Per-panel flags: run 'gum choose --help' or 'gum filter --help' to see
// all flags available inside each panel block.
package panel

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options for the panel command.
type Options struct {
	Panel []string `arg:"" help:"Panel blocks: -- choose [choose-flags] items... -- filter [filter-flags] items... (see 'gum choose --help' / 'gum filter --help' for per-panel flags)"`

	// Layout
	Vertical bool `help:"Arrange panels vertically instead of horizontally" group:"Layout"`
	Gap      int  `help:"Space between panels" default:"1" env:"GUM_PANEL_GAP" group:"Layout"`
	Height   int  `help:"Height of each panel" default:"10" env:"GUM_PANEL_HEIGHT" group:"Layout"`
	Border   string `help:"Border style (none, single, double, rounded)" default:"single" env:"GUM_PANEL_BORDER" group:"Layout"`
	Reverse  bool   `help:"Display from the bottom of the screen" default:"true" negatable:"" env:"GUM_PANEL_REVERSE" group:"Layout"`

	// Output
	Stacked         bool   `help:"Separate panel outputs with --delimiter" default:"true" negatable:"" group:"Output"`
	Delimiter       string `help:"Separator printed between panel outputs (requires --stacked)" default:"---" group:"Output"`
	OutputDelimiter string `help:"Delimiter between multiple selections within one panel" default:"|" env:"GUM_PANEL_OUTPUT_DELIMITER" group:"Output"`

	// Selection
	Single bool `help:"Enter selects current item and exits immediately" group:"Selection"`
	All    bool `help:"Require at least one selection in every panel before submitting" group:"Selection"`
	Active int  `help:"Index of the initially active panel (0-based)" default:"0" group:"Selection"`

	// Misc
	ShowHelp   bool          `help:"Show key bindings footer" default:"true" negatable:"" env:"GUM_PANEL_SHOW_HELP" group:"Misc"`
	CustomHelp string        `help:"Custom text for the help footer (overrides auto-generated key bindings)" default:"" env:"GUM_PANEL_CUSTOM_HELP" group:"Misc"`
	Timeout    time.Duration `help:"Auto-exit after this duration (0 = disabled)" default:"0s" env:"GUM_PANEL_TIMEOUT" group:"Misc"`
	Debug      bool          `help:"Print debug info to stderr" default:"false" negatable:"" env:"GUM_PANEL_DEBUG" group:"Misc"`

	// STDIN
	InputDelimiter string `help:"Delimiter used when reading items from STDIN" default:" " env:"GUM_PANEL_INPUT_DELIMITER" group:"STDIN"`
	StripANSI      bool   `help:"Strip ANSI color codes from STDIN input" default:"true" negatable:"" env:"GUM_PANEL_STRIP_ANSI" group:"STDIN"`

	// Border styles (global — apply to all panels)
	ActiveBorderStyle   style.Styles `embed:"" prefix:"active-border." set:"defaultForeground=212" envprefix:"GUM_PANEL_ACTIVE_BORDER_" group:"Style Flags"`
	InactiveBorderStyle style.Styles `embed:"" prefix:"inactive-border." set:"defaultForeground=240" envprefix:"GUM_PANEL_INACTIVE_BORDER_" group:"Style Flags"`

	// Common styles (global — apply to all panels)
	CursorStyle         style.Styles `embed:"" prefix:"cursor." set:"defaultForeground=212" envprefix:"GUM_PANEL_CURSOR_" group:"Style Flags"`
	HeaderStyle         style.Styles `embed:"" prefix:"header." set:"defaultForeground=99" envprefix:"GUM_PANEL_HEADER_" group:"Style Flags"`
	ItemStyle           style.Styles `embed:"" prefix:"item." envprefix:"GUM_PANEL_ITEM_" group:"Style Flags"`
	SelectedItemStyle   style.Styles `embed:"" prefix:"selected." set:"defaultForeground=212" envprefix:"GUM_PANEL_SELECTED_" group:"Style Flags"`
	MatchStyle          style.Styles `embed:"" prefix:"match." set:"defaultForeground=212" envprefix:"GUM_PANEL_MATCH_" group:"Style Flags"`
	IndicatorStyle      style.Styles `embed:"" prefix:"indicator." set:"defaultForeground=212" envprefix:"GUM_PANEL_INDICATOR_" group:"Style Flags"`

	// Choose-specific styles (global)
	SelectedPrefixStyle   style.Styles `embed:"" prefix:"selected-indicator." set:"defaultForeground=212" envprefix:"GUM_PANEL_SELECTED_PREFIX_" group:"Style Flags"`
	UnselectedPrefixStyle style.Styles `embed:"" prefix:"unselected-prefix." set:"defaultForeground=240" envprefix:"GUM_PANEL_UNSELECTED_PREFIX_" group:"Style Flags"`

	// Filter-specific styles (global)
	TextStyle        style.Styles `embed:"" prefix:"text." envprefix:"GUM_PANEL_TEXT_" group:"Style Flags"`
	CursorTextStyle  style.Styles `embed:"" prefix:"cursor-text." envprefix:"GUM_PANEL_CURSOR_TEXT_" group:"Style Flags"`
	PromptStyle      style.Styles `embed:"" prefix:"prompt." set:"defaultForeground=240" envprefix:"GUM_PANEL_PROMPT_" group:"Style Flags"`
	PlaceholderStyle style.Styles `embed:"" prefix:"placeholder." set:"defaultForeground=240" envprefix:"GUM_PANEL_PLACEHOLDER_" group:"Style Flags"`
}
