package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func TestParseArgs(t *testing.T) {
	// Setup flags
	rootCmd.SetArgs([]string{
		"--input", "my_input",
		"--output", "my_output",
		"--max-sources", "10",
		"--max-words", "5000",
		"--ignore-file", ".myignore",
		"--strategy", "dir",
	})

	originalRunE := rootCmd.RunE
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error { return nil }
	defer func() { rootCmd.RunE = originalRunE }()

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if InputPath != "my_input" {
		t.Errorf("Expected InputPath 'my_input', got '%s'", InputPath)
	}
	if OutputPath != "my_output" {
		t.Errorf("Expected OutputPath 'my_output', got '%s'", OutputPath)
	}
	if MaxSources != 10 {
		t.Errorf("Expected MaxSources 10, got %d", MaxSources)
	}
	if MaxWords != 5000 {
		t.Errorf("Expected MaxWords 5000, got %d", MaxWords)
	}
	if IgnoreFile != ".myignore" {
		t.Errorf("Expected IgnoreFile '.myignore', got '%s'", IgnoreFile)
	}
	if Strategy != "dir" {
		t.Errorf("Expected Strategy 'dir', got '%s'", Strategy)
	}
}

func TestGranularity(t *testing.T) {
	FS = afero.NewMemMapFs()
	afero.WriteFile(FS, "file1.txt", make([]byte, 500000), 0644) // 500k bytes -> ~100k words

	InputPath = "./"
	MaxWords = 10000
	
	// Test 1: Should pass if max-sources is high enough (100k words / 10k = 10 files needed)
	MaxSources = 15
	err := runChunking()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test 2: Should fail strictly if max-sources is too low
	MaxSources = 5
	err = runChunking()
	if err == nil {
		t.Fatalf("Expected strictly failing error due to max-sources limit, got nil")
	}

	if !strings.Contains(err.Error(), "FATAL: Project is too large") {
		t.Errorf("Expected FATAL capacity error, got: %v", err)
	}
}
