// Package files handles files.
package files

import (
	"os"
	"path/filepath"
	"strings"
)

// List returns a list of all files in the current directory.
// It ignores the .git directory.
func List() []string {
	return ListFrom(".", false, false)
}

// ListFrom returns a list of files starting from root. If showHidden is true,
// dot-prefixed files/directories are included. If includeDirs is true,
// directories are included in the results.
func ListFrom(root string, showHidden, includeDirs bool) []string {
	var files []string
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil //nolint:nilerr
			}
			rel, relErr := filepath.Rel(root, path)
			if relErr != nil {
				return nil
			}
			if rel == "." {
				return nil
			}
			if shouldIgnoreEntry(rel, showHidden) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if info.IsDir() {
				if includeDirs {
					files = append(files, rel)
				}
				return nil
			}
			files = append(files, rel)
			return nil
		})
	if err != nil {
		return []string{}
	}
	return files
}

var defaultIgnorePatterns = []string{"node_modules", ".git", "."}

func shouldIgnore(path string) bool {
	return shouldIgnoreEntry(path, false)
}

func shouldIgnoreEntry(path string, showHidden bool) bool {
	base := filepath.Base(path)
	if !showHidden && strings.HasPrefix(base, ".") {
		return true
	}
	for _, pattern := range []string{"node_modules", ".git"} {
		if base == pattern {
			return true
		}
	}
	return false
}
