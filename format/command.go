// Package format allows you to render formatted text from the command line.
//
// It supports the following types:
//
// 1. Markdown
// 2. Code
// 3. Emoji
// 4. Template
//
// For more information, see the format/README.md file.
package format

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/gum/internal/stdin"
)

// Run runs the format command.
func (o Options) Run() error {
	var input, output string
	var err error
	if len(o.Template) > 0 {
		input = strings.Join(o.Template, "\n")
	} else {
		input, _ = stdin.Read(stdin.StripANSI(o.StripANSI))
	}

	switch o.Type {
	case "code":
		output, err = code(input, o.Language)
	case "emoji":
		output, err = emoji(input)
	case "template":
		output, err = template(input)
	default:
		output, err = markdown(input, o.Theme)
	}
	if err != nil {
		return err
	}

	fmt.Print(output)
	return nil
}
