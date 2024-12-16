package stdin

// StdinOptions provide switches for ReadWithOptions, indicating how to read
// input from stdin
type StdinOptions interface {
	// DoStripANSICodes returns true if the ANSI codes should be stripped from the
	// input.
	DoStripANSICodes() bool
}

// ReadWithOptions delegates to the module's Read functions based on the
// provided options instance.
func ReadWithOptions(opts StdinOptions) (string, error) {
	if opts.DoStripANSICodes() {
		return ReadStrip()
	} else {
		return Read()
	}
}
