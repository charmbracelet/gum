package files

import (
	"os"
	"path/filepath"
	"strings"
)

// List returns a list of all files in the current directory.
// It ignores the .git directory.
func List() []string {
	var files []string
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if shouldIgnore(path) || info.IsDir() || err != nil {
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

func shouldIgnore(path string) bool {
	if strings.HasPrefix(path, ".git") {
		return true
	}
	if strings.HasPrefix(path, ".") {
		return true
	}
	return false
}
