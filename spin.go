package main

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
type Spin struct {
	Color   string `help:"Spinner color" default:"#FF06B7"`
	Display string `help:"Text to display to user while spinning" default:"Loading..."`
	Spinner string `help:"Spinner type" enum:"line,dot,minidot,jump,pulse,points,globe,moon,monkey,meter,hamburger" default:"dot"`
}
