package exit

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
)

// StatusTimeout is the exit code for timed out commands.
const StatusTimeout = 124

// StatusAborted is the exit code for aborted commands.
const StatusAborted = 130

// ErrAborted is the error to return when a gum command is aborted by Ctrl+C.
var ErrAborted = fmt.Errorf("aborted")

// NewTimeout returns a new ErrTimeout.
func NewTimeout(d time.Duration) ErrTimeout {
	return ErrTimeout{d: d}
}

// ErrTimeout is a time out error.
type ErrTimeout struct {
	d time.Duration
}

func (e ErrTimeout) Error() string {
	return "timed out after " + e.d.String()
}

func (e ErrTimeout) Unwrap() error {
	return huh.ErrTimeout
}
