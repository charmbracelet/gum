package completion

type Completion struct {
	Complete Complete `cmd:"" hidden:"" help:"Request shell completion"`
	Bash     Bash     `cmd:"" help:"Generate the autocompletion script for bash"`
	Zsh      Zsh      `cmd:"" help:"Generate the autocompletion script for zsh"`
	// Fish Fish `cmd:"" help:"Generate the autocompletion script for fish"`
}
