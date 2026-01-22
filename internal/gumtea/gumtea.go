// Package gumtea wraps tea.NewProgram to allow custom options.
package gumtea

import (
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// NewProgram wraps tea.NewProgram, injecting options from the environment.
func NewProgram(model tea.Model, baseOpts ...tea.ProgramOption) *tea.Program {
	opts := append(baseOpts, loadOptionsFromEnv()...)
	return tea.NewProgram(model, opts...)
}

// loadOptionsFromEnv loads options from environment variables.
// This feature is provisional. It may be altered or removed in a future version of this package.
func loadOptionsFromEnv() []tea.ProgramOption {
	var opts []tea.ProgramOption

	if fps := os.Getenv("GUM_FPS"); fps != "" {
		if v, err := strconv.Atoi(fps); err == nil && v > 0 {
			opts = append(opts, tea.WithFPS(v))
		}
	}

	if alt := os.Getenv("GUM_ALTSCREEN"); alt != "" {
		if altEnabled, err := strconv.ParseBool(alt); err == nil && altEnabled {
			opts = append(opts, tea.WithAltScreen())
		}
	}

	return opts
}
