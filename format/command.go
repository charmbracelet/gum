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
	"fmt"
	"strings"

	"github.com/charmbracelet/gum/internal/stdin"
)

// Func is a function that formats some text.
type Func func(string) (string, error)

var formatType = map[string]Func{
	"code":     code,
	"emoji":    emoji,
	"markdown": markdown,
	"template": template,
}

// Run runs the format command.
func (o Options) Run() error {
	var input string
	if len(o.Template) > 0 {
		input = strings.Join(o.Template, "\n")
	} else {
		input, _ = stdin.Read()
	}

	v, err := formatType[o.Type](input)
	if err != nil {
		return err
	}

	fmt.Println(v)
	return nil
}
