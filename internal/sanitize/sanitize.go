package sanitize

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var invalidFileChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
var consecutiveUnderscores = regexp.MustCompile(`_+`)

// Filename returns a filesystem-safe name from title, falling back to id if title is empty.
func Filename(title, id string) string {
	name := strings.TrimSpace(title)
	if name == "" {
		name = id
	}
	name = invalidFileChars.ReplaceAllString(name, "_")
	name = consecutiveUnderscores.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	if name == "" {
		name = "untitled"
	}
	if len(name) > 100 {
		name = name[:100]
	}
	return name
}

// SafePath returns an absolute path for filename+ext within outputDir.
// Returns an error if the resulting path would escape outputDir.
func SafePath(outputDir, filename, ext string) (string, error) {
	absDir, err := filepath.Abs(outputDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve output directory: %w", err)
	}
	candidate := filepath.Join(absDir, filename+ext)
	if !strings.HasPrefix(candidate, absDir+string(filepath.Separator)) {
		return "", fmt.Errorf("filename %q escapes output directory", filename)
	}
	return candidate, nil
}

// MakeUnique appends a count suffix to filename if it has already been used.
func MakeUnique(filename string, used map[string]int) string {
	if count, exists := used[filename]; exists {
		return fmt.Sprintf("%s_%d", filename, count+1)
	}
	return filename
}
