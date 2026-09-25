package style

import "testing"

func TestBorderKeys(t *testing.T) {
	// All border keys advertised in the CLI enum must be present in the map.
	wantKeys := []string{"double", "hidden", "none", "normal", "rounded", "thick"}
	for _, key := range wantKeys {
		if _, ok := Border[key]; !ok {
			t.Errorf("Border map missing key %q", key)
		}
	}
	if got, want := len(Border), len(wantKeys); got != want {
		t.Errorf("Border map size = %d, want %d", got, want)
	}
}
