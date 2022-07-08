package log

import "fmt"

// Error prints an error message to the user.
func Error(message string) {
	fmt.Println("Error:", message)
}

// Info prints an informational message to the user.
func Info(message string) {
	fmt.Println(message)
}
