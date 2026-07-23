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

	totalBytes, virtualTree, err := ScanDirectory(fs, ".", ".myignore", nil)

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

// Test 3.1: force-ext rescues a file that would otherwise be ignored
func TestForceExtRescuesIgnoredFile(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("src", 0755)
	afero.WriteFile(fs, "src/config.json", make([]byte, 100), 0644)
	afero.WriteFile(fs, "src/main.go", make([]byte, 200), 0644)

	// Create ignore file that excludes .json files
	ignoreContent := "*.json\n"
	afero.WriteFile(fs, ".gitignore", []byte(ignoreContent), 0644)

	// Force-ext .json to rescue it
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{".json"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should include both files: config.json is force-included, main.go is included normally
	if len(virtualTree) != 2 {
		t.Errorf("Expected 2 files, got %d", len(virtualTree))
	}
}

// Test 3.2: multiple --force-ext values — each matching extension is rescued
func TestMultipleForceExtValues(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("src", 0755)
	afero.WriteFile(fs, "src/config.json", make([]byte, 100), 0644)
	afero.WriteFile(fs, "src/lock.json", make([]byte, 50), 0644)
	afero.WriteFile(fs, "src/main.go", make([]byte, 200), 0644)

	ignoreContent := "*.json\n*.lock\n"
	afero.WriteFile(fs, ".gitignore", []byte(ignoreContent), 0644)

	// Force-ext multiple extensions
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{".json", ".lock"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should include all three files
	if len(virtualTree) != 3 {
		t.Errorf("Expected 3 files, got %d", len(virtualTree))
	}
}

// Test 3.3: non-matching extension with an ignore rule is still excluded
func TestNonMatchingExtensionStillExcluded(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("src", 0755)
	afero.WriteFile(fs, "src/config.json", make([]byte, 100), 0644)
	afero.WriteFile(fs, "src/main.go", make([]byte, 200), 0644)

	ignoreContent := "*.json\n"
	afero.WriteFile(fs, ".gitignore", []byte(ignoreContent), 0644)

	// Force-ext only .lock, not .json
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{".lock"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only include main.go (config.json is still ignored)
	if len(virtualTree) != 1 {
		t.Errorf("Expected 1 file, got %d", len(virtualTree))
	}
}

// Test 3.4: empty forceExts — no change to existing ignore behaviour
func TestEmptyForceExtsNoChange(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("src", 0755)
	afero.WriteFile(fs, "src/config.json", make([]byte, 100), 0644)
	afero.WriteFile(fs, "src/main.go", make([]byte, 200), 0644)

	ignoreContent := "*.json\n"
	afero.WriteFile(fs, ".gitignore", []byte(ignoreContent), 0644)

	// Empty forceExts - same behavior as before
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only include main.go (config.json is ignored)
	if len(virtualTree) != 1 {
		t.Errorf("Expected 1 file, got %d", len(virtualTree))
	}
}

// Test 3.5: hidden file with matching extension is still excluded (hidden guard wins)
func TestHiddenFileWithMatchingExtensionExcluded(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("src", 0755)
	afero.WriteFile(fs, "src/.env", make([]byte, 100), 0644)
	afero.WriteFile(fs, "src/main.go", make([]byte, 200), 0644)

	// Force-ext .env but it should still be excluded because it's hidden
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{".env"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only include main.go (.env is hidden, force-ext doesn't override)
	if len(virtualTree) != 1 {
		t.Errorf("Expected 1 file, got %d", len(virtualTree))
	}
}

// Test 3.6: extension matching is case-insensitive (.JSON matches config.JSON)
func TestForceExtCaseInsensitive(t *testing.T) {
	fs := afero.NewMemMapFs()

	fs.MkdirAll("src", 0755)
	afero.WriteFile(fs, "src/config.JSON", make([]byte, 100), 0644)
	afero.WriteFile(fs, "src/main.go", make([]byte, 200), 0644)

	ignoreContent := "*.JSON\n"
	afero.WriteFile(fs, ".gitignore", []byte(ignoreContent), 0644)

	// Force-ext with lowercase .json should match .JSON file
	_, virtualTree, err := ScanDirectory(fs, ".", ".gitignore", []string{".json"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// config.JSON should be force-included despite being in ignore file
	if len(virtualTree) != 2 {
		t.Errorf("Expected 2 files, got %d", len(virtualTree))
	}
}
