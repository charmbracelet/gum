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
func (m Man) BeforeApply(app *kong.Kong) error {
	man := mangokong.NewManPage(1, app.Model)
	man = man.WithSection("Copyright", "(C) 2021-2022 Charmbracelet, Inc.\n"+
		"Released under MIT license.")
	fmt.Fprint(app.Stdout, man.Build(roff.NewDocument()))
	app.Exit(0)
	return nil
}
