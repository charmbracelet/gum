package format

// Options is customization options for the format command
type Options struct {
	Template []string `arg:"" optional:"" help:"Template string to format (can also be provided via stdin)"`

	Type string `help:"Format to use (markdown,template,code,emoji)" enum:"markdown,template,code,emoji" short:"t" default:"markdown"`
}
