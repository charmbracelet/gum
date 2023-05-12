package format

// Options is customization options for the format command.
type Options struct {
	Template []string `arg:"" optional:"" help:"Template string to format (can also be provided via stdin)"`
	Theme    string   `help:"Glamour theme to use for markdown formatting" default:"pink" env:"GUM_FORMAT_THEME"`
	Language string   `help:"Programming language to parse code" short:"l" default:"" env:"GUM_FORMAT_LANGUAGE"`

	Type string `help:"Format to use (markdown,template,code,emoji)" enum:"markdown,template,code,emoji" short:"t" default:"markdown" env:"GUM_FORMAT_TYPE"`
}
