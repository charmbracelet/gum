package file

import (
	"time"

	"github.com/charmbracelet/gum/style"
)

// Options are the options for the file command.
type Options struct {
	// Path is the path to the folder / directory to begin traversing.
	Path string `arg:"" optional:"" name:"path" help:"The path to the folder to begin traversing" env:"GUM_FILE_PATH"`
	// Cursor is the character to display in front of the current selected items.
	Cursor      string        `short:"c" help:"The cursor character" default:">" env:"GUM_FILE_CURSOR"`
	All         bool          `short:"a" help:"Show hidden and 'dot' files" default:"false" env:"GUM_FILE_ALL"`
	Permissions bool          `short:"p" help:"Show file permissions" default:"true" negatable:"" env:"GUM_FILE_PERMISSION"`
	Size        bool          `short:"s" help:"Show file size" default:"true" negatable:"" env:"GUM_FILE_SIZE"`
	File        bool          `help:"Allow files selection" default:"true" env:"GUM_FILE_FILE"`
	Directory   bool          `help:"Allow directories selection" default:"false" env:"GUM_FILE_DIRECTORY"`
	ShowHelp    bool          `help:"Show help key binds" negatable:"" default:"true" env:"GUM_FILE_SHOW_HELP"`
	Timeout     time.Duration `help:"Timeout until command aborts without a selection" default:"0s" env:"GUM_FILE_TIMEOUT"`
	Header      string        `help:"Header value" default:"" env:"GUM_FILE_HEADER"`
	Height      int           `help:"Maximum number of files to display" default:"10" env:"GUM_FILE_HEIGHT"`

	CursorStyle      style.Styles `embed:"" prefix:"cursor." help:"The cursor style" set:"defaultForeground=212" envprefix:"GUM_FILE_CURSOR_"`
	SymlinkStyle     style.Styles `embed:"" prefix:"symlink." help:"The style to use for symlinks" set:"defaultForeground=36" envprefix:"GUM_FILE_SYMLINK_"`
	DirectoryStyle   style.Styles `embed:"" prefix:"directory." help:"The style to use for directories" set:"defaultForeground=99" envprefix:"GUM_FILE_DIRECTORY_"`
	FileStyle        style.Styles `embed:"" prefix:"file." help:"The style to use for files" envprefix:"GUM_FILE_FILE_"`
	PermissionsStyle style.Styles `embed:"" prefix:"permissions." help:"The style to use for permissions" set:"defaultForeground=244" envprefix:"GUM_FILE_PERMISSIONS_"`
	SelectedStyle    style.Styles `embed:"" prefix:"selected." help:"The style to use for the selected item" set:"defaultBold=true" set:"defaultForeground=212" envprefix:"GUM_FILE_SELECTED_"`                    //nolint:staticcheck
	FileSizeStyle    style.Styles `embed:"" prefix:"file-size." help:"The style to use for file sizes" set:"defaultWidth=8" set:"defaultAlign=right" set:"defaultForeground=240"  envprefix:"GUM_FILE_FILE_SIZE_"` //nolint:staticcheck
	HeaderStyle      style.Styles `embed:"" prefix:"header." set:"defaultForeground=99" envprefix:"GUM_FILE_HEADER_"`
	Padding          string       `help:"Padding" default:"${defaultPadding}" group:"Style Flags" env:"GUM_FILE_PADDING"`
}
