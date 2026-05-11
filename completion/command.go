package completion

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

// Completion command.
type Completion struct {
	Bash Bash `cmd:"" help:"Generate the autocompletion script for bash"`
	Zsh  Zsh  `cmd:"" help:"Generate the autocompletion script for zsh"`
	Fish Fish `cmd:"" help:"Generate the autocompletion script for fish"`
}

func commandName(cmd *kong.Node) string {
	commandName := cmd.FullPath()
	commandName = strings.ReplaceAll(commandName, " ", "_")
	commandName = strings.ReplaceAll(commandName, ":", "__")
	return commandName
}

func hasCommands(cmd *kong.Node) bool {
	for _, c := range cmd.Children {
		if !c.Hidden {
			return true
		}
	}
	return false
}

//nolint:deadcode,unused
func isArgument(cmd *kong.Node) bool {
	return cmd.Type == kong.ArgumentNode
}

// writeString writes a string into a buffer, and checks if the error is not nil.
func writeString(b io.StringWriter, s string) {
	if _, err := b.WriteString(s); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func nonCompletableFlag(flag *kong.Flag) bool {
	return flag.Hidden
}

func flagPossibleValues(flag *kong.Flag) []string {
	values := make([]string, 0)
	for _, enum := range flag.EnumSlice() {
		if strings.TrimSpace(enum) != "" {
			values = append(values, enum)
		}
	}
	return values
}
