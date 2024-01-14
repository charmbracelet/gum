package format

import (
	"bytes"
	"fmt"
	tpl "text/template"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
)

func code(input, language string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return "", fmt.Errorf("unable to create renderer: %w", err)
	}
	output, err := renderer.Render(fmt.Sprintf("```%s\n%s\n```", language, input))
	if err != nil {
		return "", fmt.Errorf("unable to render: %w", err)
	}
	return output, nil
}

func emoji(input string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithEmoji(),
	)
	if err != nil {
		return "", fmt.Errorf("unable to create renderer: %w", err)
	}
	output, err := renderer.Render(input)
	if err != nil {
		return "", fmt.Errorf("unable to render: %w", err)
	}
	return output, nil
}

func markdown(input string, theme string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath(theme),
		glamour.WithWordWrap(0),
	)
	if err != nil {
		return "", fmt.Errorf("unable to render: %w", err)
	}
	output, err := renderer.Render(input)
	if err != nil {
		return "", fmt.Errorf("unable to render: %w", err)
	}
	return output, nil
}

func template(input string) (string, error) {
	f := termenv.TemplateFuncs(termenv.ANSI256)
	t, err := tpl.New("tpl").Funcs(f).Parse(input)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, nil)
	return buf.String(), err
}
