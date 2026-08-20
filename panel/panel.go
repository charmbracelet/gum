package panel

import (
	"github.com/charmbracelet/gum/choose"
	"github.com/charmbracelet/gum/filter"
)

// PanelType represents the type of panel (choose or filter).
//
//revive:disable:var-naming
type PanelType string

//revive:enable:var-naming

const (
	// PanelChoose is a choose-style panel where users select from a list.
	PanelChoose PanelType = "choose"
	// PanelFilter is a filter-style panel with fuzzy search.
	PanelFilter PanelType = "filter"
)

// Panel represents a single panel configuration.
type Panel struct {
	Type       PanelType
	ModelIdx   int
	ChooseOpts *choose.Options // set when Type == PanelChoose
	FilterOpts *filter.Options // set when Type == PanelFilter
}
