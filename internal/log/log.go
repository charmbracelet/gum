package log

import "fmt"

// Error prints an error message to the user.
func Error(message string) {
	fmt.Println("Error:", message)
}
