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
	if !o.shouldRenderAsMarkdown() {
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

func (o Options) shouldRenderAsMarkdown() bool {
	return strings.HasSuffix(o.File, ".md") || o.File == ""
}
