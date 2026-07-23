//go:build windows

package cmd

import (
	"os"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

func init() {
	rawCmdLine := windowsUTF16PtrToString(syscall.GetCommandLine())
	os.Args = parseCommandLine(rawCmdLine)
}

func windowsUTF16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	length := 0
	for {
		val := *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(length)*2))
		if val == 0 {
			break
		}
		length++
	}
	return syscall.UTF16ToString(unsafe.Slice(p, length))
}


func parseCommandLine(cmdLine string) []string {
	var args []string
	var current strings.Builder
	inQuotes := false
	seenQuotes := false
	i := 0
	n := len(cmdLine)

	for i < n {
		r := rune(cmdLine[i])

		if inQuotes {
			if r == '"' {
				inQuotes = false
				i++
			} else if r == '\\' && i+1 < n && cmdLine[i+1] == '"' {
				// If we see \", on Windows this is usually the trailing backslash escape issue.
				// We treat it as a literal backslash and the closing quote.
				current.WriteRune('\\')
				inQuotes = false
				i += 2 // skip both \ and "
			} else {
				current.WriteRune(r)
				i++
			}
		} else {
			if unicode.IsSpace(r) {
				if current.Len() > 0 || seenQuotes {
					args = append(args, current.String())
					current.Reset()
					seenQuotes = false
				}
				i++
			} else if r == '"' {
				inQuotes = true
				seenQuotes = true
				i++
			} else {
				current.WriteRune(r)
				i++
			}
		}
	}

	if current.Len() > 0 || seenQuotes {
		args = append(args, current.String())
	}

	return args
}
