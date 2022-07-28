package main

import (
	"github.com/alecthomas/kong"
	"github.com/charmbracelet/gum/choose"
	"github.com/charmbracelet/gum/completion"
	"github.com/charmbracelet/gum/confirm"
	"github.com/charmbracelet/gum/filter"
	"github.com/charmbracelet/gum/format"
	"github.com/charmbracelet/gum/input"
	"github.com/charmbracelet/gum/join"
	"github.com/charmbracelet/gum/man"
	"github.com/charmbracelet/gum/spin"
	"github.com/charmbracelet/gum/style"
	"github.com/charmbracelet/gum/write"
)

// Gum is the command-line interface for Gum.
type Gum struct {
	// Version is a flag that can be used to display the version number.
	Version kong.VersionFlag `short:"v" help:"Print the version number"`

	// Completion generates Gum shell completion scripts.
	Completion completion.Completion `cmd:"" hidden:"" help:"Request shell completion"`

	// Man is a hidden command that generates Gum man pages.
	Man man.Man `cmd:"" hidden:"" help:"Generate man pages"`

	// Choose provides an interface to choose one option from a given list of
	// options. The options can be provided as (new-line separated) stdin or a
	// list of arguments.
	//
	// It is different from the filter command as it does not provide a fuzzy
	// finding input, so it is best used for smaller lists of options.
	//
	// Let's pick from a list of gum flavors:
	//
	//   $ gum choose "Strawberry" "Banana" "Cherry"
	//
	Choose choose.Options `cmd:"" help:"Choose an option from a list of choices"`

	// Confirm provides an interface to ask a user to confirm an action.
	// The user is provided with an interface to choose an affirmative or
	// negative answer, which is then reflected in the exit code for use in
	// scripting.
	//
	// If the user selects the affirmative answer, the program exits with 0.
	// If the user selects the negative answer, the program exits with 1.
	//
	// I.e. confirm if the user wants to delete a file
	//
	//   $ gum confirm "Are you sure?" && rm file.txt
	//
	Confirm confirm.Options `cmd:"" help:"Ask a user to confirm an action"`

	// Filter provides a fuzzy searching text input to allow filtering a list of
	// options to select one option.
	//
	// By default it will list all the files (recursively) in the current directory
	// for the user to choose one, but the script (or user) can provide different
	// new-line separated options to choose from.
	//
	// I.e. let's pick from a list of gum flavors:
	//
	//   $ cat flavors.text | gum filter
	//
	Filter filter.Options `cmd:"" help:"Filter items from a list"`

	// Format allows you to render styled text from `markdown`, `code`,
	// `template` strings, or embedded `emoji` strings.
	// For more information see the format/README.md file.
	Format format.Options `cmd:"" help:"Format a string using a template"`

	// Input provides a shell script interface for the text input bubble.
	// https://github.com/charmbracelet/bubbles/tree/master/textinput
	//
	// It can be used to prompt the user for some input. The text the user
	// entered will be sent to stdout.
	//
	//   $ gum input --placeholder "What's your favorite gum?" > answer.text
	//
	Input input.Options `cmd:"" help:"Prompt for some input"`

	// Join provides a shell script interface for the lipgloss JoinHorizontal
	// and JoinVertical commands. It allows you to join multi-line text to
	// build different layouts.
	//
	// For example, you can place two bordered boxes next to each other:
	// Note: We wrap the variable in quotes to ensure the new lines are part of a
	// single argument. Otherwise, the command won't work as expected.
	//
	//   $ gum join --horizontal "$BUBBLE_BOX" "$GUM_BOX"
	//
	//   ╔══════════════════════╗╔═════════════╗
	//   ║                      ║║             ║
	//   ║        Bubble        ║║     Gum     ║
	//   ║                      ║║             ║
	//   ╚══════════════════════╝╚═════════════╝
	//
	Join join.Options `cmd:"" help:"Join text vertically or horizontally"`

	// Spin provides a shell script interface for the spinner bubble.
	// https://github.com/charmbracelet/bubbles/tree/master/spinner
	//
	// It is useful for displaying that some task is running in the background
	// while consuming it's output so that it is not shown to the user.
	//
	// For example, let's do a long running task: $ sleep 5
	//
	// We can simply prepend a spinner to this task to show it to the user,
	// while performing the task / command in the background.
	//
	//   $ gum spin -t "Taking a nap..." -- sleep 5
	//
	// The spinner will automatically exit when the task is complete.
	//
	Spin spin.Options `cmd:"" help:"Display spinner while running a command"`

	// Style provides a shell script interface for Lip Gloss.
	// https://github.com/charmbracelet/lipgloss
	//
	// It allows you to use Lip Gloss to style text without needing to use Go.
	// All of the styling options are available as flags.
	//
	// Let's make some text glamorous using bash:
	//
	//   $ gum style \
	//  	--foreground 212 --border double --align center \
	//  	--width 50 --margin 2 --padding "2 4" \
	//  	"Bubble Gum (1¢)" "So sweet and so fresh\!"
	//
	//
	//    ╔══════════════════════════════════════════════════╗
	//    ║                                                  ║
	//    ║                                                  ║
	//    ║                 Bubble Gum (1¢)                  ║
	//    ║              So sweet and so fresh!              ║
	//    ║                                                  ║
	//    ║                                                  ║
	//    ╚══════════════════════════════════════════════════╝
	//
	Style style.Options `cmd:"" help:"Apply coloring, borders, spacing to text"`

	// Write provides a shell script interface for the text area bubble.
	// https://github.com/charmbracelet/bubbles/tree/master/textarea
	//
	// It can be used to ask the user to write some long form of text
	// (multi-line) input. The text the user entered will be sent to stdout.
	//
	//   $ gum write > output.text
	//
	Write write.Options `cmd:"" help:"Prompt for long-form text"`
}
