package tail

import (
	"bufio"
	"fmt"
	"os"
)

// Run executes the tail command.
func (o Options) Run() error {
	linesChan := make(chan string)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			linesChan <- scanner.Text()
		}
		close(linesChan)
	}()

	lastLines := make([]string, 0, o.NumLines)
	frameHeight := 0

	for line := range linesChan {
		lastLines = appendLine(lastLines, line, o.NumLines)

		// Clear previous frame
		if frameHeight > 0 {
			fmt.Printf("\033[%dA", frameHeight) // Move cursor up
		}

		// Print current frame
		for _, l := range lastLines {
			fmt.Printf("\r\033[K%s\n", l) // Clear line + print
		}

		frameHeight = len(lastLines)
	}

	// After finishing, clear the live frame and print only final lines
	if frameHeight > 0 {
		fmt.Printf("\033[%dA", frameHeight)
	}
	fmt.Print("\033[J") // Clear everything below cursor
	for _, l := range lastLines {
		fmt.Println(l)
	}

	return nil
}

func appendLine(lines []string, line string, maxLines int) []string {
	lines = append(lines, line)
	if len(lines) > maxLines {
		lines = lines[1:]
	}
	return lines
}
