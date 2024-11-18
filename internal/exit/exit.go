package exit

import (
	"errors"
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
)

// StatusTimeout is the exit code for timed out commands.
const StatusTimeout = 124

// StatusAborted is the exit code for aborted commands.
const StatusAborted = 130

// ErrAborted is the error to return when a gum command is aborted by Ctrl+C.
var ErrAborted = huh.ErrUserAborted

// Handle handles the error.
func Handle(err error, d time.Duration) error {
	if errors.Is(err, huh.ErrTimeout) {
		return fmt.Errorf("%w after %s", huh.ErrTimeout, d)
	}
	return err
}
