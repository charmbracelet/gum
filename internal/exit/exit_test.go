package exit

import "testing"

func TestErrExitError(t *testing.T) {
	e := ErrExit(42)
	if e.Error() != "exit 42" {
		t.Errorf("expected 'exit 42', got %q", e.Error())
	}
}

func TestStatusConstants(t *testing.T) {
	if StatusTimeout != 124 {
		t.Errorf("expected StatusTimeout=124, got %d", StatusTimeout)
	}
	if StatusAborted != 130 {
		t.Errorf("expected StatusAborted=130, got %d", StatusAborted)
	}
}
