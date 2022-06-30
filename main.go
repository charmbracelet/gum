package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/seashell/input"
	"github.com/charmbracelet/seashell/search"
	"github.com/charmbracelet/seashell/spin"
	"github.com/charmbracelet/seashell/style"
	"github.com/muesli/coral"
)

var rootCmd = &coral.Command{
	Use:   "sea",
	Short: "Sea Shell is a shell integration for bubbles and lipgloss for input and layout management",
	Long:  "Sea Shell is a wrapper for using Bubbles components\nand Lip Gloss layouts in the terminal directly in shell scripts",
}

// Execute executes the root command
func main() {
	rootCmd.AddCommand(input.Cmd())
	rootCmd.AddCommand(search.Cmd())
	rootCmd.AddCommand(spin.Cmd())
	rootCmd.AddCommand(style.Cmd())

	// We re-add a help command to remove the shorthand, so that we can use the
	// -h shorthand for height
	rootCmd.PersistentFlags().BoolP("help", "", false, "Show help for a command")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
