package main

// Input provides a shell script interface for the text input bubble.
// https://github.com/charmbracelet/bubbles/textinput
//
// It can be used to prompt the user for some input. The text the user entered
// will be sent to stdout.
//
//   $ pop input --prompt "? " --placeholder \
//		 What's your favorite soda pop flavor?" > answer.text
//
type Input struct {
	Placeholder string `help:"Placeholder value" default:"..."`
	Prompt      string `help:"Prompt to display" default:"> "`
	Width       int    `help:"Input width" default:"20"`
}
