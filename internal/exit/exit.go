package exit

import "fmt"

// StatusAborted is the exit code for aborted commands.
const StatusAborted = 130

// ErrAborted is the error to return when a gum command is aborted by Ctrl + C.
var ErrAborted = fmt.Errorf("aborted")
