package join

// Options is the set of options that can configure a join.
type Options struct {
	Text []string `arg:"" help:"Text to join."`

	Align      string `help:"Text alignment" enum:"left,center,right,bottom,middle,top" default:"left"`
	Horizontal bool   `help:"Join (potentially multi-line) strings horizontally"`
	Vertical   bool   `help:"Join (potentially multi-line) strings vertically"`
}
