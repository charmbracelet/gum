// Package format provides a way to format a string using a template.
// Behind the scenes it uses the termenv templating helpers to format
// the string. https://github.com/muesli/termenv
//
// Use it to quickly print styled strings without needing to manually
// mess with lots of style commands. If you need more powerful styling
// use the `gum style` and `gum join` to build up output.
//
//   $ gum format '{{ Bold "Tasty" }} {{ Underline "Bubble" }} {{ Foreground "212" "Gum" }}'
//
// Or, pass the format string over stdin:
//
//   $ printf '{{ Bold "Tasty" }} {{ Underline "Bubble" }} {{ Foreground "212" "Gum" }}' | gum format
//   $ printf 'Inline {{ Bold (Color "#eb5757" "#292927" " code ") }} block' | gum format
//
// Markdown also works!
//
//  $ gum format --type markdown '**Tasty** ~~Bubble~~ `Gum`'
//
package format

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/muesli/termenv"
)

// Run runs the format command
func (o Options) Run() error {
	in, err := stdin.Read()
	if in != "" || err != nil {
		o.Template = in
	}

	switch o.Type {
	case "markdown", "md":
		renderer, err := glamour.NewTermRenderer(glamour.WithStandardStyle("pink"))
		if err != nil {
			fmt.Println(err)
			return err
		}
		out, err := renderer.Render(o.Template)
		if err != nil {
			return err
		}
		fmt.Println(out)
	case "template", "tpl":
		f := termenv.TemplateFuncs(termenv.ColorProfile())
		tpl := template.New("tpl").Funcs(f)
		tpl, err = tpl.Parse(o.Template)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		tpl.Execute(&buf, nil)
		fmt.Println(&buf)
	}
	return nil
}
