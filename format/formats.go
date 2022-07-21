package format

import (
	"bytes"
	"fmt"
	tpl "text/template"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
)

var code Func = func(input string) string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return ""
	}
	out, err := renderer.Render(fmt.Sprintf("```\n%s\n```", input))
	if err != nil {
		return ""
	}
	return out
}

var emoji Func = func(input string) string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithEmoji(),
	)
	if err != nil {
		return ""
	}
	out, err := renderer.Render(input)
	if err != nil {
		return ""
	}
	return out
}

var markdown Func = func(input string) string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("pink"),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return ""
	}
	out, err := renderer.Render(input)
	if err != nil {
		return ""
	}
	return out
}

var template Func = func(input string) string {
	f := termenv.TemplateFuncs(termenv.ColorProfile())
	t, err := tpl.New("tpl").Funcs(f).Parse(input)
	if err != nil {
		return ""
	}
	var buf bytes.Buffer
	t.Execute(&buf, nil)
	return buf.String()
}
