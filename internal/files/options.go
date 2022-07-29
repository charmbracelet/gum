package files

// ListOptions is a flag set containing options
// which control the output of List.
//
// This flag set is primarily used by the filter command.
type ListOptions struct {
	OnlyDirectories bool `help:"Only return directory names in results (in other words, omit file names)" default:"false" group:"File Listing Flags"`
}
