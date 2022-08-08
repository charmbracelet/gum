package exit

import "fmt"

// StatusAborted is the exit code for aborted commands.
const StatusAborted = 130

// ErrAborted is the error to return when a gum command is aborted by Ctrl + C.
var ErrAborted = fmt.Errorf("aborted")

// StatusEmptyList is the exit code for empty list
const StatusEmptyList = 1

// ErrEmptyList is the error to return when a gum command that should operate
// on a list receives an empty list.
var ErrEmptyList = fmt.Errorf("empty")
