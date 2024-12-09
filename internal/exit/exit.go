package exit

import (
	"errors"
	"strconv"
)

// StatusTimeout is the exit code for timed out commands.
const StatusTimeout = 124

// StatusAborted is the exit code for aborted commands.
const StatusAborted = 130

// ErrAborted is the error to return when a gum command is aborted by Ctrl+C.
var ErrAborted = errors.New("user aborted")

// ErrTimeout is the error returned when the timeout is reached.
var ErrTimeout = errors.New("timeout")

// ErrExit is a custom exit error.
type ErrExit int

// Error implements error.
func (e ErrExit) Error() string { return "exit " + strconv.Itoa(int(e)) }
