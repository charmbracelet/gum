package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListDoesNotPanicOnPermissionError(t *testing.T) {
	// Create a temp directory with a file we can't access
	dir := t.TempDir()

	// Create a normal file
	if err := os.WriteFile(filepath.Join(dir, "normal.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a directory with no read permissions
	unreadableDir := filepath.Join(dir, "unreadable")
	if err := os.Mkdir(unreadableDir, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		// Restore permissions so cleanup works
		os.Chmod(unreadableDir, 0o755) //nolint:errcheck
	})

	// Change to the temp directory
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Chdir(origDir) //nolint:errcheck
	})

	// List should not panic even with permission errors
	files := List()
	if len(files) == 0 {
		t.Log("List returned empty (may be due to permission error), but did not panic")
	}
}

func TestListIgnoresGitDir(t *testing.T) {
	dir := t.TempDir()

	// Create a normal file
	if err := os.WriteFile(filepath.Join(dir, "normal.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Create a .git directory with a file
	gitDir := filepath.Join(dir, ".git")
	if err := os.Mkdir(gitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte("git"), 0o644); err != nil {
		t.Fatal(err)
	}

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Chdir(origDir) //nolint:errcheck
	})

	files := List()
	for _, f := range files {
		if filepath.Base(f) == "config" && filepath.Dir(f) == ".git" {
			t.Errorf("List should not include .git/config, but got %s", f)
		}
	}
}
