package tty

import (
	"os"
	"testing"

	"github.com/charmbracelet/colorprofile"
)

func TestClicolorForce(t *testing.T) {
	// Test that when CLICOLOR_FORCE=1 is set, the writer uses TrueColor profile
	// This is a regression test for issue #1127
	originalValue := os.Getenv("CLICOLOR_FORCE")
	os.Setenv("CLICOLOR_FORCE", "1")
	defer os.Setenv("CLICOLOR_FORCE", originalValue)

	w := Writer()

	// When CLICOLOR_FORCE=1, the writer should have TrueColor profile
	// which passes ANSI sequences through without stripping or downsampling
	wrapper, ok := w.(*colorprofile.Writer)
	if !ok {
		t.Fatalf("Writer() returned %T, expected *colorprofile.Writer", w)
	}

	if wrapper.Profile != colorprofile.TrueColor {
		t.Errorf("expected Profile to be TrueColor when CLICOLOR_FORCE=1, got %v", wrapper.Profile)
	}
}

func TestNoClicolorForce(t *testing.T) {
	// Test that when CLICOLOR_FORCE is not set, the writer uses default detection
	originalValue := os.Getenv("CLICOLOR_FORCE")
	os.Unsetenv("CLICOLOR_FORCE")
	defer func() {
		if originalValue != "" {
			os.Setenv("CLICOLOR_FORCE", originalValue)
		}
	}()

	w := Writer()

	// When CLICOLOR_FORCE is not set, the writer should use default detection
	// which may strip or downgrade ANSI sequences based on TTY status
	if w == nil {
		t.Error("Writer() returned nil")
	}
}