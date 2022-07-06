package main

// Pop is the command-line interface for Soda Pop.
type Pop struct {
	// Input provides a shell script interface for the text input bubble.
	// https://github.com/charmbracelet/bubbles/textinput
	//
	// It can be used to prompt the user for some input. The text the user
	// entered will be sent to stdout.
	//
	//   $ pop input --prompt "? " --placeholder "What's your favorite pop?" > answer.text
	//
	Input struct {
		Placeholder string `help:"Placeholder value" default:"..."`
		Prompt      string `help:"Prompt to display" default:"> "`
		Width       int    `help:"Input width" default:"20"`
	} `cmd:"" help:"Prompt for input."`

	// Search provides a fuzzy searching text input to allow filtering a list of
	// options to select one option.
	//
	// By default it will list all the files (recursively) in the current directory
	// for the user to choose one, but the script (or user) can provide different
	// new-line separated options to choose from.
	//
	// I.e. let's pick from a list of soda pop flavors:
	//
	// $ cat flavors.text | pop search
	//
	Search struct {
		AccentColor string `help:"Accent color for prompt, indicator, and matches" default:"#FF06B7"`
		Indicator   string `help:"Character for selection" default:"•"`
		Placeholder string `help:"Placeholder value" default:"..."`
		Prompt      string `help:"Prompt to display" default:"> "`
		Width       int    `help:"Input width" default:"20"`
	} `cmd:"" help:"Fuzzy search options."`

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
	Spin struct {
		Color   string `help:"Spinner color" default:"#FF06B7"`
		Display string `help:"Text to display to user while spinning" default:"Loading..."`
		Spinner string `help:"Spinner type" enum:"line,dot,minidot,jump,pulse,points,globe,moon,monkey,meter,hamburger" default:"dot"`
	} `cmd:"" help:"Show spinner while executing a command."`

	// Style styles some text.
	Style struct {
		Background       string `help:"Background color"`
		Foreground       string `help:"Foreground color"`
		BorderBackground string `help:"Border background color"`
		BorderForeground string `help:"Border foreground color"`
		Align            string `help:"Text alignment" enum:"left,center,right,bottom,middle,top" default:"left"`
		Border           string `help:"Border style to apply" enum:"none,hidden,normal,rounded,thick,double" default:"none"`
		Height           int    `help:"Height of output"`
		Width            int    `help:"Width of output"`
		Margin           string `help:"Margin to apply around the text."`
		Padding          string `help:"Padding to apply around the text."`
		Bold             bool   `help:"Apply bold formatting"`
		Faint            bool   `help:"Apply faint formatting"`
		Italic           bool   `help:"Apply italic formatting"`
		Strikethrough    bool   `help:"Apply strikethrough formatting"`
	} `cmd:"" help:"Style some text."`
}
