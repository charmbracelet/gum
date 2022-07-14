// Package format allows you to render formatted text from the command line.
// There are a few different types of formats that can be used:
//
// 1. Markdown
// -------------
// Render any input as markdown text. This uses glamour behind the scenes.
// Input passed as arguments are placed on new lines.
// Input passed as stdin are rendered as is.
// https://github.com/charmbracelet/glamour
//
//   $ gum format '**Tasty** *Bubble* `Gum`'
//   $ gum format "1. List" "2. Of" "3. Items"
//   $ echo "# Bubble Gum \n - List \n - of \n - Items" | gum format
//   $ gum format -- "- Bulleted" "- List"
//
// Or, pass input via stdin:
//
//   $ echo '**Tasty** *Bubble* `Gum`' | gum format
//
// 2. Template
// -------------
// Render styled input from a template. Templates are handled by termenv.
// https://github.com/muesli/termenv
//
//   $ gum format '{{ Bold "Tasty" }} {{ Italic "Bubble" }} {{ Color "99" "0" " Gum " }}' --type template
//
// 3. Code
// -------------
// Perform syntax highlighting on some code snippets. Styling is handled by
// glamour, which in turn uses chroma. https://github.com/alecthomas/chroma
//
//   $ cat code.go | gum format --type code
//
// 4. Emoji
// -------------
// Parse emojis within text and render emojis. Emoji rendering is handled by
// glamour, which in turn uses goldmark-emoji.
// https://github.com/yuin/goldmark-emoji
//
//   $ gum format 'I :heart: Bubble Gum :candy:' --type emoji
//   $ echo "I :heart: Bubble Gum :candy:" | gum format --type emoji
//
// Output:
//
//   I ❤️ Bubble Gum 🍬
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
		out, err := glamour.Render(fmt.Sprintf("```\n%s\n```", input), "auto")
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
		renderer, err := glamour.NewTermRenderer(glamour.WithEmoji())
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
