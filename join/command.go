// Package join provides a shell script interface for the lipgloss
// JoinHorizontal and JoinVertical commands. It allows you to join multi-line
// text to build different layouts.
//
// For example, you can place two bordered boxes next to each other: Note: We
// wrap the variable in quotes to ensure the new lines are part of a single
// argument. Otherwise, the command won't work as expected.
//
//   $ gum join --horizontal "$BUBBLE_BOX" "$GUM_BOX"
//
//   ╔══════════════════════╗╔═════════════╗
//   ║                      ║║             ║
//   ║        Bubble        ║║     Gum     ║
//   ║                      ║║             ║
//   ╚══════════════════════╝╚═════════════╝
//
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
