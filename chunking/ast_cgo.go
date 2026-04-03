//go:build cgo

package chunking

import (
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

func getLanguage(path string) *sitter.Language {
	ext := filepath.Ext(path)
	if ext == ".go" {
		return golang.GetLanguage()
	}
	return nil
}

func astSplit(path string, content []byte, maxWords int) []string {
	lang := getLanguage(path)
	if lang == nil {
		return nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	tree := parser.Parse(nil, content)
	if tree == nil {
		return nil
	}

	var parts []string
	var currentPart strings.Builder
	var currentWords int

	root := tree.RootNode()
	count := int(root.ChildCount())

	for i := 0; i < count; i++ {
		child := root.Child(i)
		startByte := child.StartByte()
		endByte := child.EndByte()

		if startByte >= endByte || int(endByte) > len(content) {
			continue
		}

		nodeText := string(content[startByte:endByte])
		nodeWords := CountWords(nodeText)

		if currentWords+nodeWords > maxWords && currentWords > 0 {
			parts = append(parts, currentPart.String())
			currentPart.Reset()
			currentWords = 0
		}

		currentPart.WriteString(nodeText)
		currentPart.WriteString("\n")
		currentWords += nodeWords
	}

	if currentPart.Len() > 0 {
		parts = append(parts, currentPart.String())
	}

	if len(parts) > 1 {
		return parts
	}

	return nil
}
