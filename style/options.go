package style

// Options is the customization options for the style command.
type Options struct {
	Text  []string `arg:"" optional:"" help:"Text to style"`
	Style Styles   `help:"Style to apply" embed:""`
}
