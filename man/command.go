// Package man the man command.
package man

import (
	"fmt"

	"github.com/alecthomas/kong"
	mangokong "github.com/alecthomas/mango-kong"
	"github.com/muesli/roff"
)

// Man is a gum sub-command that generates man pages.
type Man struct{}

// BeforeApply implements Kong BeforeApply hook.
func (m Man) BeforeApply(ctx *kong.Context) error {
	// Set the correct man pages description without color escape sequences.
	ctx.Model.Help = "A tool for glamorous shell scripts."
	man := mangokong.NewManPage(1, ctx.Model)
	man = man.WithSection("Copyright", "(c) 2022-2024 Charmbracelet, Inc.\n"+
		"Released under MIT license.")
	_, _ = fmt.Fprint(ctx.Stdout, man.Build(roff.NewDocument()))
	ctx.Exit(0)
	return nil
}
