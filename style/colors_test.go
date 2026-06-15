package style

import "testing"

func TestResolveColorNamed(t *testing.T) {
	t.Parallel()

	if got := resolveColor("red"); got != "#FF0000" {
		t.Fatalf("resolveColor(red) = %q, want #FF0000", got)
	}
	if got := resolveColor("#AABBCC"); got != "#AABBCC" {
		t.Fatalf("resolveColor(hex) = %q, want passthrough", got)
	}
}
