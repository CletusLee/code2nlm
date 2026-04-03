package scanner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sabhiram/go-gitignore"
	"github.com/spf13/afero"
)

// ScanDirectory sequentially walks a directory, evaluating ignore files,
// omitting hidden directories, and summing total bytes from file sizes
// without reading their full contents (Double-Pass I/O avoidance).
func ScanDirectory(fs afero.Fs, root string, ignoreFilePath string) (int64, []string, error) {
	var totalBytes int64
	var virtualTree []string

	// Initialize gitignore matcher if file exists
	var ignoreMatcher *ignore.GitIgnore
	exists, err := afero.Exists(fs, ignoreFilePath)
	if err == nil && exists {
		content, err := afero.ReadFile(fs, ignoreFilePath)
		if err == nil {
			ignoreMatcher = ignore.CompileIgnoreLines(strings.Split(string(content), "\n")...)
		}
	}

	err = afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Normalize path for check
		relPath, err := filepath.Rel(root, path)
		if err != nil {
			relPath = path
		}
		
		// Skip root itself
		if relPath == "." {
			return nil
		}

		// Skip hidden files/directories by default
		if strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check against ignore matcher
		if ignoreMatcher != nil && ignoreMatcher.MatchesPath(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.IsDir() {
			// Fast byte accumulation strategy avoiding file read
			totalBytes += info.Size()
			// Use robust normalization to POSIX style for tree
			virtualTree = append(virtualTree, strings.ReplaceAll(relPath, "\\", "/"))
		}

		return nil
	})

	return totalBytes, virtualTree, err
}
