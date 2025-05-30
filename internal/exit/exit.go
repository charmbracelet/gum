// Package exit code implementation.
package exit

import "strconv"

// StatusTimeout is the exit code for timed out commands.
const StatusTimeout = 124

// StatusAborted is the exit code for aborted commands.
const StatusAborted = 130

// ErrExit is a custom exit error.
type ErrExit int

// Error implements error.
func (e ErrExit) Error() string { return "exit " + strconv.Itoa(int(e)) }
