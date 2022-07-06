package main

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
type Search struct {
	AccentColor string `help:"Accent color for prompt, indicator, and matches" default:"#FF06B7"`
	Indicator   string `help:"Character for selection" default:"•"`
	Placeholder string `help:"Placeholder value" default:"..."`
	Prompt      string `help:"Prompt to display" default:"> "`
	Width       int    `help:"Input width" default:"20"`
}
