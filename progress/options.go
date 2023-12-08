package progress

import (
	"github.com/charmbracelet/gum/style"
)

type Options struct {
	Title         string       `help:"Text to display to user while spinning" env:"GUM_PROGRESS_TITLE"`
	TitleStyle    style.Styles `embed:"" prefix:"title." envprefix:"GUM_PROGRESS_TITLE_"`
	Format        string       `short:"f" help:"What format to use for rendering the bar. Choose from: {Iter}, {Elapsed}, {Title} and {Avg} or see --limit for more options. Unknown options remain untouched." envprefix:"GUM_PROGRESS_FORMAT"`
	ProgressColor string       `help:"Set the color for the progress" envprefix:"GUM_PROGRESS_PROGRESS_COLOR"`

	ProgressIndicator     string `help:"What indicator to use for counting progress" default:"\n" env:"GUM_PROGRESS_PROGRESS_INDICATOR"`
	HideProgressIndicator bool   `help:"Don't show the --progress-indicator in the output. Only makes sense in combination with --show-output" default:"false"`
	ShowOutput            bool   `short:"o" help:"Print what gum reads" default:"false"`

	Limit uint `short:"l" help:"Species how many items there are (enables {Bar}, {Limit}, {Remaining}, {Eta} and {Pct} to be used in --format)"`
}
