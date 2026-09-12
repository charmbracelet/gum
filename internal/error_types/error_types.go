package error_types;

type QuietError struct {
    Err error
}

func (e *QuietError) Error() string {
    if e.Err != nil {
        return e.Err.Error()
    }
    return "quiet error"
}

func (e *QuietError) Unwrap() error {
    return e.Err
}
