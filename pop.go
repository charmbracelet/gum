package main

// Pop is the command-line interface for Soda Pop.
type Pop struct {
	Input  `cmd:"" help:"Prompt for input."`
	Search `cmd:"" help:"Fuzzy search options."`
	Spin   `cmd:"" help:"Show spinner while executing a command."`
	Style  `cmd:"" help:"Style some text."`
}
