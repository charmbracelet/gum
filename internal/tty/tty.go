// Package tty provides tty-aware printing.
package tty

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/charmbracelet/colorprofile"
)

var stdout = sync.OnceValue(func() io.Writer {
	return colorprofile.NewWriter(os.Stdout, os.Environ())
})

// Writer returns a writer to stdout that downgrades or strips ansi
// sequences according to the environment (NO_COLOR, CLICOLOR, CLICOLOR_FORCE)
// and the capabilities of the terminal.
func Writer() io.Writer {
	return stdout()
}

// Println handles println, striping ansi sequences if stdout is not a tty.
func Println(s string) {
	_, _ = fmt.Fprintln(stdout(), s)
}
