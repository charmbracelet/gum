package info

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

// Options is the customization options for the info command.
type Options struct{}

// TermInfo holds information about the current terminal.
type TermInfo struct {
	ColorProfile string `json:"color_profile"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	IsTTY        bool   `json:"is_tty"`
	HasDarkBG    bool   `json:"has_dark_background"`
}

// Run provides a shell script interface for querying terminal information.
func (o Options) Run() error {
	w, h, _ := term.GetSize(os.Stdout.Fd())

	info := TermInfo{
		ColorProfile: fmt.Sprintf("%v", lipgloss.ColorProfile()),
		Width:        w,
		Height:       h,
		IsTTY:        term.IsTerminal(os.Stdout.Fd()),
		HasDarkBG:    lipgloss.HasDarkBackground(),
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal info: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
