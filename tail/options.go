package tail

type Options struct {
	NumLines	int		`short:"n" default:"4" help:"Number of lines you want to tail"`
}