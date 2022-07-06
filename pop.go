package main

import (
	"github.com/charmbracelet/sodapop/input"
	"github.com/charmbracelet/sodapop/search"
	"github.com/charmbracelet/sodapop/spin"
	"github.com/charmbracelet/sodapop/style"
	"github.com/charmbracelet/sodapop/write"
)

// Pop is the command-line interface for Soda Pop.
type Pop struct {
	// Input provides a shell script interface for the text input bubble.
	// https://github.com/charmbracelet/bubbles/textinput
	//
	// It can be used to prompt the user for some input. The text the user
	// entered will be sent to stdout.
	//
	//   $ pop input --placeholder "What's your favorite pop?" > answer.text
	//
	Input input.Options `cmd:"" help:"Prompt for input."`

	// Write provides a shell script interface for the text area bubble.
	// https://github.com/charmbracelet/bubbles/textarea
	//
	// It can be used to ask the user to write some long form of text
	// (multi-line) input. The text the user entered will be sent to stdout.
	//
	//   $ pop write > output.text
	//
	Write write.Options `cmd:"" help:"Prompt for text"`

	// Search provides a fuzzy searching text input to allow filtering a list of
	// options to select one option.
	//
	// By default it will list all the files (recursively) in the current directory
	// for the user to choose one, but the script (or user) can provide different
	// new-line separated options to choose from.
	//
	// I.e. let's pick from a list of soda pop flavors:
	//
	//   $ cat flavors.text | pop search
	//
	Search search.Options `cmd:"" help:"Fuzzy search options."`

	// Spin provides a shell script interface for the spinner bubble.
	// https://github.com/charmbracelet/bubbles/spinner
	//
	// It is useful for displaying that some task is running in the background
	// while consuming it's output so that it is not shown to the user.
	//
	// For example, let's do a long running task:
	//   $ sleep 5
	//
	// We can simply prepend a spinner to this task to show it to the user,
	// while performing the task / command in the background.
	//
	//   $ pop spin -t "Taking a nap..." -- sleep 5
	//
	// The spinner will automatically exit when the task is complete.
	Spin spin.Options `cmd:"" help:"Show spinner while executing a command."`

	// Style provides a shell script interface for Lip Gloss.
	// https://github.com/charmbracelet/lipgloss
	//
	// It allows you to use Lip Gloss to style text without needing to use Go.
	// All of the styling options are available as flags.
	//
	// Let's make some text glamorous using bash:
	//
	//   $ pop style \
	//		--foreground "#FF06B7" --border "double" \
	// 		--margin 2 --padding "2 4" --width 50 \
	//			"And oh gosh, how delicious the fabulous frizzy frobscottle" \
	//			"was! It was sweet and refreshing. It tasted of vanilla and" \
	//			"cream, with just the faintest trace of raspberries on the" \
	//			"edge of the flavour. And the bubbles were wonderful."
	//
	//
	//     ╔══════════════════════════════════════════════════╗
	//     ║                                                  ║
	//     ║                                                  ║
	//     ║    And oh gosh, how delicious the fabulous       ║
	//     ║    frizzy frobscottle was It was sweet and       ║
	//     ║    refreshing. It tasted of vanilla and          ║
	//     ║    cream, with just the faintest trace of        ║
	//     ║    raspberries on the edge of the flavour.       ║
	//     ║    And the bubbles were wonderful.               ║
	//     ║                                                  ║
	//     ║                                                  ║
	//     ╚══════════════════════════════════════════════════╝
	//
	Style style.Options `cmd:"" help:"Style some text."`
}
