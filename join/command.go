package join

import (
	"fmt"

	"github.com/charmbracelet/gum/internal/decode"
	"github.com/charmbracelet/lipgloss"
)

// Run is the command-line interface for the joining strings through lipgloss.
func (o Options) Run() error {
	join := lipgloss.JoinHorizontal
	if o.Vertical {
		join = lipgloss.JoinVertical
	}
	fmt.Println(join(decode.Align[o.Align], o.Text...))
	return nil
}
