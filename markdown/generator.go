package markdown

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// GenerateIndex creates the 00_Project_Index.md file based on the virtual tree
func GenerateIndex(fs afero.Fs, outpuDir string, virtualTree []string) error {
	// Ensure output formatting uses ToSlash
	outDirNorm := filepath.ToSlash(outpuDir)
	err := fs.MkdirAll(outDirNorm, 0755)
	if err != nil {
		return err
	}

	indexPath := filepath.ToSlash(filepath.Join(outDirNorm, "00_Project_Index.md"))
	
	var sb strings.Builder
	sb.WriteString("# Project Index\n\n")
	sb.WriteString("## Directory Structure\n\n")

	for _, p := range virtualTree {
		// Paths are already normalized in scanner, but we ensure output string is clean
		sb.WriteString(fmt.Sprintf("- `%s`\n", p))
	}

	return afero.WriteFile(fs, indexPath, []byte(sb.String()), 0644)
}

// FormatContextualHeader generates the standardized header injected into each chunk.
func FormatContextualHeader(domainName string, part int, totalParts int, projectName string, paths []string) string {
	var sb strings.Builder
	
	partStr := ""
	if totalParts > 1 {
		partStr = fmt.Sprintf(" (Part %d of %d)", part, totalParts)
	}

	sb.WriteString(fmt.Sprintf("# Module: %s%s\n", domainName, partStr))
	sb.WriteString(fmt.Sprintf("**Project**: %s\n", projectName))
	sb.WriteString("**Global Context**: Please refer to `00_Project_Index.md` for the complete directory structure and dependency map.\n\n")
	
	sb.WriteString("## Included Paths in this Chunk\n")
	for _, p := range paths {
		sb.WriteString(fmt.Sprintf("* `%s`\n", filepath.ToSlash(p)))
	}
	sb.WriteString("\n---\n\n")
	return sb.String()
}
