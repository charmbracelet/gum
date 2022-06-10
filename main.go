package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/sea-shell/input"
	"github.com/charmbracelet/sea-shell/style"
	"github.com/muesli/coral"
)

var rootCmd = &coral.Command{
	Use:   "sea",
	Short: "Sea Shell is a shell integration for bubbles and lipgloss for input and layout management",
	Long:  "Sea Shell is a wrapper for using charmbracelet/bubbles components\nand charmbracelet/lipgloss layouts in the terminal directly in shell scripts",
}

// Execute executes the root command
func main() {
	rootCmd.AddCommand(input.Cmd())
	rootCmd.AddCommand(style.Cmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
