package files

import (
	"os"
	"path/filepath"
	"strings"
)

// List returns a list of all files in the current directory.
// It ignores the .git directory.
func List(opts ListOptions) []string {
	var files []string

	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			switch {
			case err != nil:
				return nil
			case shouldIgnore(path):
				return nil
			case info.IsDir() && !opts.OnlyDirectories:
				return nil
			case !info.IsDir() && opts.OnlyDirectories:
				return nil
			}

			files = append(files, path)
			return nil
		})

	if err != nil {
		return []string{}
	}

	return files
}

var defaultIgnorePatterns = []string{"node_modules", ".git", "."}

func shouldIgnore(path string) bool {
	for _, prefix := range defaultIgnorePatterns {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	return false
}
