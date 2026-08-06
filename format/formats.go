//nolint
package format

import (
	"bytes"
	"fmt"
	tpl "text/template"

	"charm.land/glamour/v2"
	"charm.land/lipgloss/v2"
)

func code(input, language string) (string, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithEnvironmentConfig(),
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
	t, err := tpl.New("tpl").Funcs(templateFuncs()).Parse(input)
	if err != nil {
		return "", fmt.Errorf("unable to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, nil)
	return buf.String(), err
}

// templateFuncs returns a few useful template helpers for styling text.
func templateFuncs() tpl.FuncMap {
	return tpl.FuncMap{
		"Color": func(values ...any) string {
			s := lipgloss.NewStyle()
			switch len(values) {
			case 2:
				s = s.Foreground(lipgloss.Color(values[0].(string)))
			case 3:
				s = s.
					Foreground(lipgloss.Color(values[0].(string))).
					Background(lipgloss.Color(values[1].(string)))
			}
			return s.Render(values[len(values)-1].(string))
		},
		"Foreground": func(values ...any) string {
			s := lipgloss.NewStyle()
			if len(values) == 2 {
				s = s.Foreground(lipgloss.Color(values[0].(string)))
			}
			return s.Render(values[len(values)-1].(string))
		},
		"Background": func(values ...any) string {
			s := lipgloss.NewStyle()
			if len(values) == 2 {
				s = s.Background(lipgloss.Color(values[0].(string)))
			}
			return s.Render(values[len(values)-1].(string))
		},
		"Bold":      styleFunc(func(s lipgloss.Style) lipgloss.Style { return s.Bold(true) }),
		"Faint":     styleFunc(func(s lipgloss.Style) lipgloss.Style { return s.Faint(true) }),
		"Italic":    styleFunc(func(s lipgloss.Style) lipgloss.Style { return s.Italic(true) }),
		"Underline": styleFunc(func(s lipgloss.Style) lipgloss.Style { return s.Underline(true) }),
		"Overline":  styleFunc(func(s lipgloss.Style) lipgloss.Style { return s }),
		"Blink":     styleFunc(func(s lipgloss.Style) lipgloss.Style { return s.Blink(true) }),
		"Reverse":   styleFunc(func(s lipgloss.Style) lipgloss.Style { return s.Reverse(true) }),
		"CrossOut":  styleFunc(func(s lipgloss.Style) lipgloss.Style { return s.Strikethrough(true) }),
	}
}

func styleFunc(f func(lipgloss.Style) lipgloss.Style) func(...any) string {
	return func(values ...any) string {
		return f(lipgloss.NewStyle()).Render(values[0].(string))
	}
}
