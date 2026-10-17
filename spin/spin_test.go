package spin

import (
	"testing"
	"time"
)

func TestSpinnerFPS(t *testing.T) {
	tests := []struct {
		name    string
		spinner string
		fps     int
		wantFPS time.Duration
	}{
		{
			name:    "default fps unchanged when fps is 0",
			spinner: "line",
			fps:     0,
			wantFPS: spinnerMap["line"].FPS,
		},
		{
			name:    "custom fps 5",
			spinner: "line",
			fps:     5,
			wantFPS: 200 * time.Millisecond,
		},
		{
			name:    "custom fps 2",
			spinner: "dot",
			fps:     2,
			wantFPS: 500 * time.Millisecond,
		},
		{
			name:    "custom fps 60",
			spinner: "dot",
			fps:     60,
			wantFPS: time.Second / 60,
		},
		{
			name:    "negative fps ignored",
			spinner: "line",
			fps:     -5,
			wantFPS: spinnerMap["line"].FPS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{
				Spinner: tt.spinner,
				FPS:     tt.fps,
			}
			s := opts.newSpinner()
			if s.Spinner.FPS != tt.wantFPS {
				t.Errorf("got FPS %v, want %v", s.Spinner.FPS, tt.wantFPS)
			}
		})
	}
}
