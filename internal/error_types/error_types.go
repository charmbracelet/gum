package error_types;

// QuietError is a custom exit error.
type QuietError string

// Error implements error.
func (e QuietError) Error() string { return string(e) }
