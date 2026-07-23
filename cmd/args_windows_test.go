//go:build windows

package cmd

import (
	"reflect"
	"testing"
)

func TestParseCommandLine(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    `code2nlm -i "C:\My Folder\Sub Folder\" -o "C:\My Folder\Output\"`,
			expected: []string{"code2nlm", "-i", "C:\\My Folder\\Sub Folder\\", "-o", "C:\\My Folder\\Output\\"},
		},
		{
			input:    `code2nlm.exe -i "d:\AI Test\code2nlm\"`,
			expected: []string{"code2nlm.exe", "-i", "d:\\AI Test\\code2nlm\\"},
		},
		{
			input:    `code2nlm -i C:\Folder\ -o C:\Output\`,
			expected: []string{"code2nlm", "-i", "C:\\Folder\\", "-o", "C:\\Output\\"},
		},
		{
			input:    `code2nlm -i "C:\Folder with spaces" -o "C:\Output Folder"`,
			expected: []string{"code2nlm", "-i", "C:\\Folder with spaces", "-o", "C:\\Output Folder"},
		},
		{
			input:    `code2nlm   -i    "C:\Folder"`,
			expected: []string{"code2nlm", "-i", "C:\\Folder"},
		},
		{
			input:    `"C:\Program Files\code2nlm.exe" -i "d:\path\"`,
			expected: []string{"C:\\Program Files\\code2nlm.exe", "-i", "d:\\path\\"},
		},
		{
			input:    `code2nlm -i "" -o "C:\Out"`,
			expected: []string{"code2nlm", "-i", "", "-o", "C:\\Out"},
		},
		{
			input:    `code2nlm -i "" -o ""`,
			expected: []string{"code2nlm", "-i", "", "-o", ""},
		},
	}

	for _, tt := range tests {
		result := parseCommandLine(tt.input)
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf("For input:\n  %s\nExpected:\n  %v\nGot:\n  %v", tt.input, tt.expected, result)
		}
	}
}
