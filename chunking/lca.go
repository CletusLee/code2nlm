package chunking

import (
	"path/filepath"
	"strings"
)

// GetLCA finds the Lowest Common Ancestor directory for a set of file paths.
// The paths are expected to be in POSIX format (slashes).
func GetLCA(paths []string) string {
	if len(paths) == 0 {
		return ""
	}
	if len(paths) == 1 {
		dir := filepath.Dir(paths[0])
		if dir == "." {
			return ""
		}
		return dir
	}

	// Split each path into components
	var splitPaths [][]string
	for _, p := range paths {
		dir := filepath.ToSlash(filepath.Dir(p))
		if dir == "." || dir == "" {
			splitPaths = append(splitPaths, []string{})
			continue
		}
		splitPaths = append(splitPaths, strings.Split(dir, "/"))
	}

	// Find minimum length among all split paths
	minLen := len(splitPaths[0])
	for i := 1; i < len(splitPaths); i++ {
		if len(splitPaths[i]) < minLen {
			minLen = len(splitPaths[i])
		}
	}

	if minLen == 0 {
		return ""
	}

	// Compare components one by one
	var common []string
	for i := 0; i < minLen; i++ {
		target := splitPaths[0][i]
		allMatch := true
		for j := 1; j < len(splitPaths); j++ {
			if splitPaths[j][i] != target {
				allMatch = false
				break
			}
		}
		if allMatch {
			common = append(common, target)
		} else {
			break
		}
	}

	return strings.Join(common, "/")
}

// NormalizeLCA converts a path-style LCA to an underscore-separated filename prefix.
func NormalizeLCA(lca string) string {
	if lca == "" || lca == "." {
		return "global"
	}
	
	// Replace slashes with underscores and remove any leading/trailing underscores
	normalized := strings.ReplaceAll(lca, "/", "_")
	normalized = strings.Trim(normalized, "_")
	
	if normalized == "" {
		return "global"
	}
	return normalized
}
