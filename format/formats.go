package format

import (
	"bytes"
	"fmt"
	tpl "text/template"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
)

var code Func = func(input string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return "", err
	}
	return renderer.Render(fmt.Sprintf("```\n%s\n```", input))
}

var emoji Func = func(input string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithEmoji(),
	)
	if err != nil {
		return "", err
	}
	return renderer.Render(input)
}

var markdown Func = func(input string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("pink"),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return "", err
	}
	return renderer.Render(input)
}

var template Func = func(input string) (string, error) {
	f := termenv.TemplateFuncs(termenv.ColorProfile())
	t, err := tpl.New("tpl").Funcs(f).Parse(input)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, nil)
	return buf.String(), err
}
