package parser

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/weaby/adr-buddy/internal/model"
)

// detectCommentStyle returns the comment marker based on file extension
func detectCommentStyle(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	// Languages using # for comments
	hashCommentExts := map[string]bool{
		".py":   true,
		".rb":   true,
		".sh":   true,
		".bash": true,
		".yml":  true,
		".yaml": true,
	}

	if hashCommentExts[ext] {
		return "#"
	}

	// Default to // for most languages
	return "//"
}

// isAnnotationLine checks if a line contains an annotation
func isAnnotationLine(line, commentStyle string) bool {
	trimmed := strings.TrimSpace(line)

	// Must start with comment marker
	if !strings.HasPrefix(trimmed, commentStyle) {
		return false
	}

	// Remove comment marker and check for @decision
	content := strings.TrimSpace(strings.TrimPrefix(trimmed, commentStyle))
	return strings.HasPrefix(content, "@decision.")
}

// parseAnnotationField extracts the field name and value from an annotation line
// Returns (field, value, true) if valid, ("", "", false) otherwise
func parseAnnotationField(line, commentStyle string) (string, string, bool) {
	trimmed := strings.TrimSpace(line)

	// Remove comment marker
	if !strings.HasPrefix(trimmed, commentStyle) {
		return "", "", false
	}
	content := strings.TrimSpace(strings.TrimPrefix(trimmed, commentStyle))

	// Check for @decision prefix
	if !strings.HasPrefix(content, "@decision.") {
		return "", "", false
	}

	// Remove @decision. prefix
	content = strings.TrimPrefix(content, "@decision.")

	// Split on first colon
	parts := strings.SplitN(content, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	field := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	return field, value, true
}

// isContinuationLine checks if a line is a continuation of a multi-line field
func isContinuationLine(line, commentStyle string) bool {
	trimmed := strings.TrimSpace(line)

	// Must start with comment marker
	if !strings.HasPrefix(trimmed, commentStyle) {
		return false
	}

	// Remove comment marker
	content := strings.TrimPrefix(trimmed, commentStyle)

	// Must have leading spaces (indentation)
	if len(content) == 0 || content[0] != ' ' {
		return false
	}

	// Must not be a new annotation field
	trimmedContent := strings.TrimSpace(content)
	return !strings.HasPrefix(trimmedContent, "@decision.")
}

// extractContinuationValue extracts the text from a continuation line
func extractContinuationValue(line, commentStyle string) string {
	trimmed := strings.TrimSpace(line)
	content := strings.TrimPrefix(trimmed, commentStyle)
	return strings.TrimSpace(content)
}

// ParseFile parses a single file and extracts all annotations
func ParseFile(filePath string) ([]*model.Annotation, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	commentStyle := detectCommentStyle(filePath)
	var annotations []*model.Annotation
	var current *model.Annotation
	var currentField string

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check if this is a new annotation field
		if field, value, ok := parseAnnotationField(line, commentStyle); ok {
			// If we don't have a current annotation, create one
			if current == nil {
				current = &model.Annotation{
					Location: model.SourceLocation{
						File: filePath,
						Line: lineNum,
					},
					CustomFields: make(map[string]string),
				}
			}

			// Set the field value
			currentField = field
			setAnnotationField(current, field, value)
			continue
		}

		// Check if this is a continuation line
		if current != nil && isContinuationLine(line, commentStyle) {
			value := extractContinuationValue(line, commentStyle)
			appendToField(current, currentField, value)
			continue
		}

		// If we have a current annotation and this line is not part of it, save it
		if current != nil {
			annotations = append(annotations, current)
			current = nil
			currentField = ""
		}
	}

	// Save last annotation if exists
	if current != nil {
		annotations = append(annotations, current)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return annotations, nil
}

// setAnnotationField sets a field value on an annotation
func setAnnotationField(ann *model.Annotation, field, value string) {
	switch field {
	case "id":
		ann.ID = value
	case "name":
		ann.Name = value
	case "status":
		ann.Status = value
	case "category":
		ann.Category = value
	case "context":
		ann.Context = value
	case "decision":
		ann.Decision = value
	case "alternatives":
		ann.Alternatives = value
	case "consequences":
		ann.Consequences = value
	default:
		ann.CustomFields[field] = value
	}
}

// appendToField appends a value to a multi-line field
func appendToField(ann *model.Annotation, field, value string) {
	switch field {
	case "context":
		ann.Context += "\n" + value
	case "decision":
		ann.Decision += "\n" + value
	case "alternatives":
		ann.Alternatives += "\n" + value
	case "consequences":
		ann.Consequences += "\n" + value
	default:
		if existing, ok := ann.CustomFields[field]; ok {
			ann.CustomFields[field] = existing + "\n" + value
		}
	}
}

// ScanDirectory recursively scans a directory for annotations
func ScanDirectory(rootDir string, excludePatterns []string) ([]*model.Annotation, error) {
	var allAnnotations []*model.Annotation

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// Check if path should be excluded
		if shouldExclude(relPath, excludePatterns) {
			return nil
		}

		// Parse file for annotations
		annotations, err := ParseFile(path)
		if err != nil {
			// Skip files that can't be parsed
			return nil
		}

		// Update relative paths in annotations
		for _, ann := range annotations {
			ann.Location.File = relPath
		}

		allAnnotations = append(allAnnotations, annotations...)
		return nil
	})

	return allAnnotations, err
}

// shouldExclude checks if a path matches any exclude pattern
func shouldExclude(path string, patterns []string) bool {
	for _, pattern := range patterns {
		matched, err := doublestar.Match(pattern, path)
		if err == nil && matched {
			return true
		}
	}
	return false
}
