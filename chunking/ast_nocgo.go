//go:build !cgo

package chunking

// astSplit is a stub for environments without CGO.
// It returns nil to force the fallback text splitter.
func astSplit(path string, content []byte, maxWords int) []string {
	return nil
}
