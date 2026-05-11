package version

// Options is the set of options that can be used with version.
type Options struct {
	Constraint string `arg:"" help:"Semantic version constraint"`
}
