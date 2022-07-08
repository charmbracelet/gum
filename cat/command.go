package cat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/gum/internal/log"
)

// Run provides a shell script interface for Glamour rendering.
// https://github.com/charmbracelet/glamour
func (o Options) Run() {
	// Are we rendering markdown?
	// If not, let's render it as code.
	if !strings.HasSuffix(o.File, ".md") {
		o.Text = fmt.Sprintf("```\n%s\n```", strings.TrimSpace(o.Text))
		// Since this is code, let's not do word wrapping.
		o.Width = 0
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(o.Theme),
		glamour.WithWordWrap(o.Width),
	)
	if err != nil {
		log.Error(err.Error())
		return
	}

	out, err := r.Render(o.Text)
	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Print(out)
}
