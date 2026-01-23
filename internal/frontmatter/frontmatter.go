package frontmatter

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// AppleNotesIDKey is the frontmatter field name for Apple Notes ID
	AppleNotesIDKey = "apple_notes_id"
)

var (
	// frontmatterRegex matches YAML frontmatter blocks (---\n...\n---)
	frontmatterRegex = regexp.MustCompile(`(?ms)^---\s*(.*?)\s*^---\n?`)

	// appleNotesIDRegex extracts the apple_notes_id value from frontmatter
	appleNotesIDRegex = regexp.MustCompile(`(?m)^apple_notes_id:\s*(.+)\s*$`)
)

// Strip removes YAML frontmatter from content, returning only the body
func Strip(content string) string {
	return frontmatterRegex.ReplaceAllString(content, "")
}

// Extract returns just the YAML frontmatter block (including --- delimiters)
func Extract(content string) string {
	matches := frontmatterRegex.FindString(content)
	return matches
}

// GetAppleNotesID extracts the apple_notes_id field from frontmatter
// Returns empty string and error if not found
func GetAppleNotesID(content string) (string, error) {
	fm := Extract(content)
	if fm == "" {
		return "", fmt.Errorf("no frontmatter found")
	}

	matches := appleNotesIDRegex.FindStringSubmatch(fm)
	if len(matches) < 2 {
		return "", fmt.Errorf("no %s found in frontmatter", AppleNotesIDKey)
	}

	return strings.TrimSpace(matches[1]), nil
}

// UpdateAppleNotesID adds or updates the apple_notes_id field in frontmatter
// If no frontmatter exists, creates new frontmatter with just the ID
// If frontmatter exists, updates or adds the apple_notes_id field
func UpdateAppleNotesID(content, id string) string {
	fm := Extract(content)
	body := Strip(content)

	if fm == "" {
		// No frontmatter, create new one
		return fmt.Sprintf("---\n%s: %s\n---\n\n%s", AppleNotesIDKey, id, body)
	}

	// Extract the inner YAML content (without --- delimiters)
	inner := strings.TrimPrefix(fm, "---\n")
	inner = strings.TrimPrefix(inner, "---")
	inner = strings.TrimSuffix(inner, "---\n")
	inner = strings.TrimSuffix(inner, "---")
	inner = strings.TrimSpace(inner)

	// Remove existing apple_notes_id lines
	lines := strings.Split(inner, "\n")
	var filtered []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), AppleNotesIDKey+":") {
			if line != "" {
				filtered = append(filtered, line)
			}
		}
	}

	// Build new frontmatter
	var newFrontmatter strings.Builder
	newFrontmatter.WriteString("---\n")
	for _, line := range filtered {
		newFrontmatter.WriteString(line + "\n")
	}
	newFrontmatter.WriteString(fmt.Sprintf("%s: %s\n", AppleNotesIDKey, id))
	newFrontmatter.WriteString("---\n\n")
	newFrontmatter.WriteString(body)

	return newFrontmatter.String()
}
