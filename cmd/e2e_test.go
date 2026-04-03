package cmd

import (
	"strings"
	"testing"
	"path/filepath"

	"github.com/spf13/afero"
)

func TestE2E_CLI(t *testing.T) {
	// 1. Setup mock FS fixture
	FS = afero.NewMemMapFs()
	afero.WriteFile(FS, "src/main.go", []byte("package main\n\nfunc main(){\n}\n"), 0644)
	afero.WriteFile(FS, "src/utils.go", []byte("package main\n\nfunc util(){\n}\n"), 0644)

	// Set required arguments
	InputPath = "./"
	OutputPath = "./output"
	MaxWords = 50000
	MaxSources = 50
	IgnoreFile = ".gitignore"

	// 2. Execute process
	err := runChunking()
	if err != nil {
		t.Fatalf("Unexpected E2E error: %v", err)
	}

	// 3. Verify Index exists and has paths
	indexPath := filepath.ToSlash(filepath.Join(OutputPath, "00_Project_Index.md"))
	exists, err := afero.Exists(FS, indexPath)
	if err != nil || !exists {
		t.Fatalf("Index file %s missing", indexPath)
	}
	
	idxContent, _ := afero.ReadFile(FS, indexPath)
	if !strings.Contains(string(idxContent), "src/main.go") {
		t.Errorf("Index should map src/main.go")
	}

	// 4. Verify Chunk generation and context headers
	chunk1Path := filepath.ToSlash(filepath.Join(OutputPath, "01_chunk.md"))
	exists, err = afero.Exists(FS, chunk1Path)
	if err != nil || !exists {
		t.Fatalf("First chunk file missing")
	}

	chunkContent, _ := afero.ReadFile(FS, chunk1Path)
	contentStr := string(chunkContent)
	if !strings.Contains(contentStr, "# Module: Global") {
		t.Errorf("Missing contextual module header")
	}
	if !strings.Contains(contentStr, "src/main.go") || !strings.Contains(contentStr, "src/utils.go") {
		t.Errorf("Missing included paths context in chunk header")
	}
	if !strings.Contains(contentStr, "func main()") {
		t.Errorf("Missing code content in chunk")
	}
}
