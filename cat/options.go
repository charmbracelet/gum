package cat

// Options is the customization options for the cat command.
type Options struct {
	// This field is used to store the text read from stdin or a file.
	// It is not shown to the user as a flag.
	Text string `hidden:""`

	File  string `arg:"" description:"Name of file to read" optional:""`
	Theme string `description:"Styling theme" enum:"dark,light,dracula,notty" default:"dark"`
	Width int    `description:"Wrap width" default:"80"`
}
