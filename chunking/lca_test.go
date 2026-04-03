package chunking

import (
	"testing"
)

func TestLCAAlgorithm(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		expected string
	}{
		{"Empty", []string{}, ""},
		{"Single", []string{"src/main.go"}, "src"},
		{"Same Dir", []string{"src/a.go", "src/b.go"}, "src"},
		{"Deep Intersection", []string{"src/pkg/a.go", "src/pkg/b.go"}, "src/pkg"},
		{"Root Intersection", []string{"src/a.go", "cmd/main.go"}, ""},
		{"Triple Split", []string{"a/b/c/1.go", "a/b/d/2.go", "a/e/f/3.go"}, "a"},
		{"Deep Root Intersection", []string{"a/b/c.go", "d/e/f.go"}, ""},
		{"Same Path", []string{"src/a.go", "src/a.go"}, "src"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := GetLCA(tt.paths)
			if res != tt.expected {
				t.Errorf("GetLCA(%v) = %q, want %q", tt.paths, res, tt.expected)
			}
		})
	}
}

func TestNormalizeLCA(t *testing.T) {
	tests := []struct {
		lca      string
		expected string
	}{
		{"", "global"},
		{".", "global"},
		{"src", "src"},
		{"src/pkg", "src_pkg"},
		{"a/b/c", "a_b_c"},
	}

	for _, tt := range tests {
		t.Run(tt.lca, func(t *testing.T) {
			res := NormalizeLCA(tt.lca)
			if res != tt.expected {
				t.Errorf("NormalizeLCA(%q) = %q, want %q", tt.lca, res, tt.expected)
			}
		})
	}
}
