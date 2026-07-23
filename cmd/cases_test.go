package cmd

import (
	"strings"
	"testing"
	"path/filepath"

	"github.com/spf13/afero"
)

func TestMultiPartChunking(t *testing.T) {
	// Setup
	FS = afero.NewMemMapFs()
	// Large file: 300 words (approx)
	largeContent := strings.Repeat("word ", 300)
	afero.WriteFile(FS, "large.txt", []byte(largeContent), 0644)

	InputPath = "./"
	OutputPath = "./output"
	MaxWords = 100 // Should split into ~3 chunks
	MaxSources = 50
	IgnoreFile = ".gitignore"

	err := runChunking()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify chunks
	// Chunker.Process should emit global_001.md, global_002.md...
	files, _ := afero.ReadDir(FS, OutputPath)
	count := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") && strings.Contains(f.Name(), "_00") {
			count++
		}
	}

	if count < 2 {
		t.Errorf("Expected at least 2 chunks for 300 words with MaxWords=100, got %d", count)
	}
}

func TestSmallFileGrouping(t *testing.T) {
	// Setup
	FS = afero.NewMemMapFs()
	afero.WriteFile(FS, "a.txt", []byte("small file a"), 0644)
	afero.WriteFile(FS, "b.txt", []byte("small file b"), 0644)

	InputPath = "./"
	OutputPath = "./output"
	MaxWords = 1000
	MaxSources = 50

	err := runChunking()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should be grouped into <projectName>_global_001.md
	projectName := filepath.Base(filepath.Clean(InputPath))
	chunkName := projectName + "_global_001.md"
	chunkPath := filepath.ToSlash(filepath.Join(OutputPath, chunkName))
	exists, _ := afero.Exists(FS, chunkPath)
	if !exists {
		t.Fatalf("Expected %s to exist", chunkName)
	}

	content, _ := afero.ReadFile(FS, chunkPath)
	contentStr := string(content)
	if !strings.Contains(contentStr, "### File: `a.txt`") || !strings.Contains(contentStr, "### File: `b.txt`") {
		t.Errorf("Chunk should contain both files")
	}
}

func TestPathNormalizationDeep(t *testing.T) {
	// Setup
	FS = afero.NewMemMapFs()
	// Force a deep Windows-like path manually in virtual tree mock if possible, 
	// but scanner already converts to Slash. 
	// We'll verify that the output reflects this.
	afero.WriteFile(FS, "dir\\sub\\file.txt", []byte("content"), 0644)

	InputPath = "./"
	OutputPath = "./output"
	MaxWords = 1000
	MaxSources = 50

	runChunking()

	projectName := filepath.Base(filepath.Clean(InputPath))
	indexName := projectName + "_000_Project_Index.md"
	indexPath := filepath.ToSlash(filepath.Join(OutputPath, indexName))
	idx, _ := afero.ReadFile(FS, indexPath)
	
	if strings.Contains(string(idx), "\\") {
		t.Errorf("Index should not contain backslashes, found in: %s", string(idx))
	}
}
