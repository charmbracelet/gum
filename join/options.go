package join

// Options is the set of options that can configure a join.
type Options struct {
	Text []string `arg:"" help:"Text to join."`

	Align      string `help:"Text alignment" enum:"left,center,right,bottom,middle,top" default:"left" env:"GUM_JOIN_ALIGN"`
	Horizontal bool   `help:"Join (potentially multi-line) strings horizontally" env:"GUM_JOIN_HORIZONTAL"`
	Vertical   bool   `help:"Join (potentially multi-line) strings vertically" env:"GUM_JOIN_VERTICAL"`
}
