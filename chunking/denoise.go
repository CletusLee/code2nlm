package chunking

import "regexp"

var (
	// Regex for source map comments (inline base64 mappings) which consume huge tokens unnecessarily
	sourceMapRegex = regexp.MustCompile(`(?m)^(//|/\*)\s*#\s*sourceMappingURL=.*$`)
	
	// Regex for inline base64 data URIs often found in HTML/JS/CSS
	base64DataURIRegex = regexp.MustCompile(`data:[^;]+;base64,[a-zA-Z0-9/+=]+`)
)

// DenoiseContent strips away known token-heavy non-semantic strings from the content 
// to ensure the LLM's context window isn't wasted on unreadable machine artifacts.
func DenoiseContent(content string) string {
	// Remove source maps entirely from the content
	content = sourceMapRegex.ReplaceAllString(content, "")
	
	// Replace heavy Base64 strings with a tiny placeholder
	content = base64DataURIRegex.ReplaceAllString(content, "[BASE64_DATA_REMOVED_FOR_LLM]")
	
	return content
}
