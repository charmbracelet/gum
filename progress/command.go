package progress

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/internal/timeout"
	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for the progress bar bubble.
// https://github.com/charmbracelet/bubbles/progress
func (o Options) Run() error {
	if stdin.IsEmpty() {
		return fmt.Errorf("no input provided, pipe percentage values (0-100) to stdin")
	}

	var opts []progress.Option
	if o.Solid != "" {
		opts = append(opts, progress.WithSolidFill(o.Solid))
	} else {
		opts = append(opts, progress.WithGradient(o.GradientFrom, o.GradientTo))
	}
	if !o.ShowPercentage {
		opts = append(opts, progress.WithoutPercentage())
	}
	if o.Width > 0 {
		opts = append(opts, progress.WithWidth(o.Width))
	}

	p := progress.New(opts...)
	top, right, bottom, left := style.ParsePadding(o.Padding)

	m := model{
		progress:   p,
		title:      o.Title,
		titleStyle: o.TitleStyle.ToLipgloss(),
		padding:    []int{top, right, bottom, left},
	}

	ctx, cancel := timeout.Context(o.Timeout)
	defer cancel()

	program := tea.NewProgram(
		m,
		tea.WithOutput(os.Stderr),
		tea.WithContext(ctx),
		tea.WithInput(nil),
	)

	// Read percentage values from stdin in a goroutine.
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			val, err := strconv.ParseFloat(line, 64)
			if err != nil {
				continue
			}
			// Accept values 0-100, normalize to 0.0-1.0
			if val > 1.0 {
				val /= 100.0
			}
			if val < 0 {
				val = 0
			}
			if val > 1.0 {
				val = 1.0
			}
			program.Send(tickMsg(val))
		}
		// When stdin closes, signal completion.
		program.Send(doneMsg{})
	}()

	_, err := program.Run()
	if err != nil {
		return fmt.Errorf("unable to run progress: %w", err)
	}

	return nil
}
