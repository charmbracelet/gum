package file

import "github.com/charmbracelet/gum/style"

// Options are the options for the file command.
type Options struct {
	// Path is the path to the folder / directory to begin traversing.
	Path string `arg:"" optional:"" name:"path" help:"The path to the folder to begin traversing"`
	// Cursor is the character to display in front of the current selected items.
	Cursor string `short:"c" help:"The cursor character" default:">"`
	All    bool   `short:"a" help:"Show hidden and 'dot' files" default:"true"`

	Height           int          `help:"Maximum number of files to display" default:"0"`
	CursorStyle      style.Styles `embed:"" prefix:"cursor." help:"The cursor style" set:"defaultForeground=212" envprefix:"GUM_FILE_CURSOR_"`
	SymlinkStyle     style.Styles `embed:"" prefix:"symlink." help:"The style to use for symlinks" set:"defaultForeground=36" envprefix:"GUM_FILE_SYMLINK_"`
	DirectoryStyle   style.Styles `embed:"" prefix:"directory." help:"The style to use for directories" set:"defaultForeground=99" envprefix:"GUM_FILE_DIRECTORY_"`
	FileStyle        style.Styles `embed:"" prefix:"file." help:"The style to use for files" envprefix:"GUM_FILE_FILE_"`
	PermissionsStyle style.Styles `embed:"" prefix:"permissions." help:"The style to use for permissions" set:"defaultForeground=244" envprefix:"GUM_FILE_PERMISSIONS_"`
	SelectedStyle    style.Styles `embed:"" prefix:"selected." help:"The style to use for the selected item" set:"defaultBold=true" set:"defaultForeground=212" envprefix:"GUM_FILE_SELECTED_"`
	FileSizeStyle    style.Styles `embed:"" prefix:"file-size." help:"The style to use for file sizes" set:"defaultWidth=8" set:"defaultAlign=right" set:"defaultForeground=240"  envprefix:"GUM_FILE_FILE_SIZE_"`
}
