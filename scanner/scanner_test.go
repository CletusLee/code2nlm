package scanner

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestScanner(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Setup mock FS
	fs.MkdirAll("src/domain", 0755)
	fs.MkdirAll("node_modules", 0755)
	fs.MkdirAll(".git", 0755)

	afero.WriteFile(fs, "src/domain/file1.ts", make([]byte, 100), 0644)
	afero.WriteFile(fs, "src/domain/file2.exe", make([]byte, 200), 0644)
	afero.WriteFile(fs, "src/domain/.hidden", make([]byte, 50), 0644)
	afero.WriteFile(fs, "node_modules/lib.js", make([]byte, 500), 0644)
	afero.WriteFile(fs, ".git/config", make([]byte, 100), 0644)

	// Create ignore file
	ignoreContent := "node_modules/\n*.exe\n"
	afero.WriteFile(fs, ".myignore", []byte(ignoreContent), 0644)

	totalBytes, virtualTree, err := ScanDirectory(fs, ".", ".myignore")
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only include src/domain/file1.ts (100 bytes)
	// Hidden files omitted (.hidden, .git)
	// Ignored omitted (node_modules, *.exe)
	
	if totalBytes != 100 {
		t.Errorf("Expected 100 total bytes, got %d", totalBytes)
	}

	expectedPaths := []string{
		filepath.ToSlash("src/domain/file1.ts"),
	}

	if !reflect.DeepEqual(virtualTree, expectedPaths) {
		t.Errorf("Expected paths %v, got %v", expectedPaths, virtualTree)
	}
}
