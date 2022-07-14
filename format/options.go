package format

// Options is customization options for the format command
type Options struct {
	Template string `arg:"" optional:"" help:"Template string to format"`

	Type string `help:"Format to use" enum:"markdown,md,template,tpl" short:"t" default:"template"`
}
