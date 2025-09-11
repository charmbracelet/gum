package tail

import (
	"bufio"
	"fmt"
	"os"
)

// Run executes the tail command.
func (o Options) Run() error {
	// Channel to communicate lines read from stdin
	lines := make(chan string)

	// Read from stdin in a separate goroutine
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	// Slice to store the last n lines
	lastLines := make([]string, 0, o.NumLines)

	// Continuously read lines from stdin and update the last n lines
	for line := range lines {
		lastLines = appendLine(lastLines, line, o.NumLines)
		clearScreen()
		printLines(lastLines, o.NumLines)
	}

	return nil
}

// appendLine appends a line to the slice, maintaining the maximum number of lines.
func appendLine(lines []string, line string, maxLines int) []string {
	lines = append(lines, line)
	if len(lines) > maxLines {
		lines = lines[1:]
	}
	return lines
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// printLines prints the last n lines to stdout
func printLines(lines []string, n int) {
	// Print an extra newline for separation of previous shell output and the tailing
	fmt.Println()

	// Print the last n lines
	start := len(lines) - n
	if start < 0 {
		start = 0
	}
	for i := start; i < len(lines); i++ {
		fmt.Println(lines[i])
	}
}
