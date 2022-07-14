// Package format allows you to render formatted text from the command line.
//
// It supports the following types:
//
//   1. Markdown
//   2. Code
//   3. Emoji
//   4. Template
//
// For more information, see the format/README.md file.
//
package format

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/muesli/termenv"
)

// Run runs the format command
func (o Options) Run() error {
	var input string
	if len(o.Template) > 0 {
		input = strings.Join(o.Template, "\n")
	} else {
		input, _ = stdin.Read()
	}

	switch o.Type {
	case "code":
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(0),
		)
		if err != nil {
			return err
		}
		out, err := renderer.Render(fmt.Sprintf("```\n%s\n```", input))
		if err != nil {
			return err
		}
		fmt.Println(out)
	case "markdown", "md":
		out, err := glamour.Render(input, "pink")
		if err != nil {
			return err
		}
		fmt.Println(out)
	case "emoji":
		renderer, err := glamour.NewTermRenderer(
			glamour.WithEmoji(),
		)
		if err != nil {
			return err
		}
		out, err := renderer.Render(input)
		if err != nil {
			return err
		}
		fmt.Println(out)
	case "template", "tpl":
		f := termenv.TemplateFuncs(termenv.ColorProfile())
		tpl, err := template.New("tpl").Funcs(f).Parse(input)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		tpl.Execute(&buf, nil)
		fmt.Println(&buf)
	}
	return nil
}
