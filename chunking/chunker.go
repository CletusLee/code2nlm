package chunking

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"code2nlm/markdown"
)

type Chunker struct {
	FS          afero.Fs
	MaxWords    int
	InputPath   string
	OutputPath  string
	ProjectName string
}

func CountWords(s string) int {
	return len(strings.Fields(s))
}

// Process drives the sequential file chunking stream. (Task 5.1)
func (c *Chunker) Process(virtualTree []string) error {
	prefixCounts := make(map[string]int)
	currentWords := 0
	var currentPaths []string
	var currentContent strings.Builder

	flushChunk := func() error {
		if currentContent.Len() == 0 {
			return nil
		}

		lca := GetLCA(currentPaths)
		prefix := NormalizeLCA(lca)
		
		prefixCounts[prefix]++
		count := prefixCounts[prefix]

		fileName := fmt.Sprintf("%s_%03d.md", prefix, count)
		filePath := strings.ReplaceAll(filepath.Join(c.OutputPath, fileName), "\\", "/")

		header := markdown.FormatContextualHeader(prefix, count, -1, c.ProjectName, currentPaths)
		finalContent := header + currentContent.String()

		err := afero.WriteFile(c.FS, filePath, []byte(finalContent), 0644)
		if err != nil {
			return err
		}

		currentWords = 0
		currentPaths = nil
		currentContent.Reset()
		return nil
	}

	for _, path := range virtualTree {
		actualFilePath := filepath.Join(c.InputPath, path)
		contentBytes, err := afero.ReadFile(c.FS, actualFilePath)
		if err != nil {
			continue // Skip unreadable files
		}
		
		contentStr := string(contentBytes)
		contentStr = DenoiseContent(contentStr)
		words := CountWords(contentStr)
		contentBytes = []byte(contentStr) // Update bytes for structural parsing if needed

		if currentWords+words > c.MaxWords {
			if currentWords > 0 {
				if err := flushChunk(); err != nil {
					return err
				}
			}

			if words > c.MaxWords {
				// File is massive, must split it. AST (Task 5.2) or Fallback (Task 5.3)
				parts := c.splitLargeFile(path, contentBytes)
				for _, part := range parts {
					partWords := CountWords(part)
					if currentWords+partWords > c.MaxWords && currentWords > 0 {
						if err := flushChunk(); err != nil {
							return err
						}
					}
					
					c.addPath(&currentPaths, path)
					currentContent.WriteString(fmt.Sprintf("\n### File: `%s`\n\n```\n%s\n```\n", strings.ReplaceAll(path, "\\", "/"), part))
					currentWords += partWords
				}
				continue
			}
		}

		c.addPath(&currentPaths, path)
		currentContent.WriteString(fmt.Sprintf("\n### File: `%s`\n\n```\n%s\n```\n", strings.ReplaceAll(path, "\\", "/"), contentStr))
		currentWords += words
	}

	return flushChunk()
}

func (c *Chunker) addPath(paths *[]string, path string) {
	for _, p := range *paths {
		if p == path {
			return
		}
	}
	*paths = append(*paths, path)
}

// Fallback mechanism to astSplit which we'll define in separate files
func (c *Chunker) splitLargeFile(path string, content []byte) []string {
	parts := astSplit(path, content, c.MaxWords)
	if len(parts) > 0 {
		return parts
	}

	return fallbackSplit(string(content), c.MaxWords)
}

// fallbackSplit performs raw text boundaries when parsers are unavailable (Task 5.3)
func fallbackSplit(content string, maxWords int) []string {
	var parts []string
	lines := strings.Split(content, "\n")

	var currentPart strings.Builder
	currentWordsCount := 0

	for _, line := range lines {
		lineWords := CountWords(line)

		if currentWordsCount+lineWords <= maxWords {
			currentPart.WriteString(line)
			currentPart.WriteString("\n")
			currentWordsCount += lineWords
			continue
		}

		if currentWordsCount > 0 {
			parts = append(parts, currentPart.String())
			currentPart.Reset()
			currentWordsCount = 0
		}

		if lineWords > maxWords {
			fields := strings.Fields(line)
			for i := 0; i < len(fields); i += maxWords {
				end := i + maxWords
				if end > len(fields) {
					end = len(fields)
				}
				parts = append(parts, strings.Join(fields[i:end], " "))
			}
		} else {
			currentPart.WriteString(line)
			currentPart.WriteString("\n")
			currentWordsCount = lineWords
		}
	}

	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	return parts
}
