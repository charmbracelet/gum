package exit

import "testing"

func TestErrExitError(t *testing.T) {
	for name, tt := range map[string]struct {
		code ErrExit
		want string
	}{
		"zero":         {code: 0, want: "exit 0"},
		"one":          {code: 1, want: "exit 1"},
		"timeout":      {code: StatusTimeout, want: "exit 124"},
		"aborted":      {code: StatusAborted, want: "exit 130"},
		"large value":  {code: 255, want: "exit 255"},
		"negative one": {code: -1, want: "exit -1"},
	} {
		t.Run(name, func(t *testing.T) {
			if got := tt.code.Error(); got != tt.want {
				t.Errorf("ErrExit(%d).Error() = %q, want %q", int(tt.code), got, tt.want)
			}
		})
	}
}

func TestStatusConstants(t *testing.T) {
	if StatusTimeout != 124 {
		t.Errorf("StatusTimeout = %d, want 124", StatusTimeout)
	}
	if StatusAborted != 130 {
		t.Errorf("StatusAborted = %d, want 130", StatusAborted)
	}
}
