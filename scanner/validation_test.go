package scanner

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
)

// TestValidation_DirectoryWithExtension ensures that directories are NEVER force-included
// even if their name ends with a forced extension.
func TestValidation_DirectoryWithExtension(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Create a directory named "config.json"
	fs.MkdirAll("config.json", 0755)
	afero.WriteFile(fs, "config.json/real_file.txt", []byte("content"), 0644)

	// Ignore all .json
	ignoreContent := "*.json\n"
	afero.WriteFile(fs, ".gitignore", []byte(ignoreContent), 0644)

	// Force-ext .json
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{".json"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// The directory "config.json" matches the extension, but it should NOT be force-included
	// as a file. Instead, it should be processed normally (checked against ignore rules).
	// In gitignore, "config.json" matches "*.json", so the directory itself is ignored.
	// Therefore, "config.json/real_file.txt" should NOT be in the tree.
	for _, path := range virtualTree {
		if path == "config.json" || path == "config.json/real_file.txt" {
			t.Errorf("Path %s should have been ignored (it's a directory or inside an ignored directory)", path)
		}
	}
}

// TestValidation_ComplexNestedForce ensures force-ext works for deeply nested files
// and respects the hidden file guard even in subdirectories.
func TestValidation_ComplexNestedForce(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("pkg/lib/config", 0755)
	afero.WriteFile(fs, "pkg/lib/config/settings.json", []byte("{}"), 0644)
	afero.WriteFile(fs, "pkg/lib/config/.secret.json", []byte("{}"), 0644)
	afero.WriteFile(fs, "pkg/lib/main.go", []byte("package main"), 0644)

	// Ignore everything in pkg/lib/config/
	ignoreContent := "pkg/lib/config/\n"
	afero.WriteFile(fs, ".gitignore", []byte(ignoreContent), 0644)

	// Force-ext .json
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{".json"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should contain:
	// 1. pkg/lib/main.go (not ignored)
	// 2. pkg/lib/config/settings.json (force-included despite parent dir being ignored)
	// Should NOT contain:
	// 3. pkg/lib/config/.secret.json (hidden file guard wins)

	expected := map[string]bool{
		"pkg/lib/main.go":             true,
		"pkg/lib/config/settings.json": true,
	}

	if len(virtualTree) != 2 {
		t.Errorf("Expected 2 files, got %d: %v", len(virtualTree), virtualTree)
	}

	for _, p := range virtualTree {
		pth := filepath.ToSlash(p)
		if !expected[pth] {
			t.Errorf("Unexpected path in tree: %s", pth)
		}
		delete(expected, pth)
	}

	if len(expected) > 0 {
		t.Errorf("Missing expected paths: %v", expected)
	}
}

// TestValidation_MultipleDots ensures that extension matching works with files
// having multiple dots (e.g. .d.ts).
func TestValidation_MultipleDots(t *testing.T) {
	fs := afero.NewMemMapFs()

	afero.WriteFile(fs, "types.d.ts", []byte(""), 0644)
	
	// Ignore .d.ts files
	ignoreContent := "*.d.ts\n"
	afero.WriteFile(fs, ".gitignore", []byte(ignoreContent), 0644)

	// Force-ext .ts -> this should NOT rescue types.d.ts if filepath.Ext returns .ts
	// Actually filepath.Ext("types.d.ts") returns ".ts".
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{".ts"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	found := false
	for _, p := range virtualTree {
		if filepath.ToSlash(p) == "types.d.ts" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected types.d.ts to be rescued by .ts extension")
	}

	// Now try specifically matching .d.ts
	// Note: filepath.Ext only returns the LAST extension.
	// So to rescue .d.ts, user must pass .ts.
}
