package exit

import (
	"errors"
	"fmt"
	"time"
)

// StatusTimeout is the exit code for timed out commands.
const StatusTimeout = 124

// StatusAborted is the exit code for aborted commands.
const StatusAborted = 130

// ErrAborted is the error to return when a gum command is aborted by Ctrl+C.
var ErrAborted = errors.New("user aborted")

// ErrTimeout is the error returned when the timeout is reached.
var ErrTimeout = errors.New("timeout")

// Handle handles the error.
func Handle(err error, d time.Duration) error {
	if errors.Is(err, ErrTimeout) {
		return fmt.Errorf("%w after %s", ErrTimeout, d)
	}
	return err
}
