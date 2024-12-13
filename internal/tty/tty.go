// Package tty provides tty-aware printing.
package tty

import (
	"fmt"
	"os"
	"sync"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/term"
)

var isTTY = sync.OnceValue(func() bool {
	return term.IsTerminal(os.Stdout.Fd())
})

// Println handles println, striping ansi sequences if stdout is not a tty.
func Println(s string) {
	if isTTY() {
		fmt.Println(s)
		return
	}
	fmt.Println(ansi.Strip(s))
}
