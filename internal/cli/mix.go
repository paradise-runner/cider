package cli

import (
	"fmt"
	"strings"

	"github.com/paradise-runner/cider/internal/convert"
	"github.com/paradise-runner/cider/internal/frontmatter"
	"github.com/paradise-runner/cider/internal/markdown"
	"github.com/paradise-runner/cider/internal/notes"
	"github.com/spf13/cobra"
)

// DiffType represents the type of differences between two markdown contents
type DiffType int

const (
	DiffTypeNone       DiffType = iota // No differences
	DiffTypeFormatOnly                 // Only formatting/whitespace differences
	DiffTypeSubstantive                // Substantive content differences
)

var mixCmd = &cobra.Command{
	Use:   "mix <file>",
	Short: "Intelligently merge changes between a Markdown file and its linked Apple Note",
	Long: `Intelligently merge changes between a Markdown file and its linked Apple Note.
The file must have an apple_notes_id in its frontmatter.

This command:
- Detects and resolves simple semantic differences (whitespace, formatting)
- If only one side changed, updates the other
- If both sides have conflicting changes, shows a diff and exits without updating
- If neither side changed, does nothing`,
	Args: cobra.ExactArgs(1),
	RunE: runMix,
}

func init() {
	rootCmd.AddCommand(mixCmd)
}

func runMix(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Read markdown file
	markdownContent, err := markdown.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error: File not found: %s", filePath)
	}

	// Extract ID from frontmatter (required for mix)
	noteID, err := frontmatter.GetAppleNotesID(markdownContent)
	if err != nil {
		return fmt.Errorf("Error: No apple_notes_id found in frontmatter")
	}

	// Initialize notes client
	notesClient := notes.NewClient()

	// Find note in Apple Notes
	found, err := notesClient.Find(noteID)
	if err != nil {
		return fmt.Errorf("Error: Failed to search for note: %v", err)
	}
	if !found {
		return fmt.Errorf("Error: Note not found in Apple Notes")
	}

	// Read note content from Apple Notes
	htmlContent, err := notesClient.Read(noteID)
	if err != nil {
		return fmt.Errorf("Error: Failed to read note content: %v", err)
	}

	// Convert remote HTML to markdown
	remoteMarkdown, err := convert.HTMLToMarkdown(htmlContent)
	if err != nil {
		return fmt.Errorf("Error: Failed to convert HTML to markdown: %v", err)
	}

	// Strip frontmatter from local content for comparison
	localMarkdown := frontmatter.Strip(markdownContent)

	// Normalize both for comparison
	normalizedRemote := normalizeContent(remoteMarkdown)
	normalizedLocal := normalizeContent(localMarkdown)

	// Determine the type of differences
	diffResult := analyzeDifferences(remoteMarkdown, localMarkdown, normalizedRemote, normalizedLocal)

	switch diffResult {
	case DiffTypeNone:
		fmt.Println("No changes detected - files are already in sync")
		return nil

	case DiffTypeFormatOnly:
		// Only formatting differences - use local version as canonical
		fmt.Println("Only formatting differences detected - synchronizing both sides with local version")
		
		// Update Apple Note with local content
		htmlToUpdate, err := convert.MarkdownToHTML(localMarkdown)
		if err != nil {
			return fmt.Errorf("Error: Failed to convert markdown to HTML: %v", err)
		}
		
		err = notesClient.Update(noteID, htmlToUpdate)
		if err != nil {
			return fmt.Errorf("Error: Failed to update Apple Note: %v", err)
		}
		
		fmt.Println("✓ Apple Note updated successfully")
		fmt.Println("✓ Files synchronized")
		return nil

	case DiffTypeSubstantive:
		// Substantive differences that can't be auto-resolved
		fmt.Println("Conflicting changes detected between Apple Note and Markdown file.")
		fmt.Println("Showing differences (Apple Note -> Markdown file):\n")
		
		err = showDiff("Apple Notes", filePath, remoteMarkdown, localMarkdown)
		if err != nil {
			if exitErr, ok := err.(*exitError); ok && exitErr.code == 1 {
				// This is expected when there are differences
				fmt.Println("\nCannot automatically resolve conflicts.")
				fmt.Println("Please review the differences and manually update one side, then run mix again.")
				return exitErr
			}
			return fmt.Errorf("Error: Failed to generate diff: %v", err)
		}
		return nil

	default:
		return fmt.Errorf("Error: Unknown diff type")
	}
}

// normalizeContent normalizes whitespace and line endings for comparison
func normalizeContent(content string) string {
	// Normalize line endings
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	
	// Trim trailing whitespace from each line
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	content = strings.Join(lines, "\n")
	
	// Trim leading/trailing whitespace from entire content
	content = strings.TrimSpace(content)
	
	return content
}

// analyzeDifferences determines the type of differences between two contents
func analyzeDifferences(raw1, raw2, normalized1, normalized2 string) DiffType {
	// If normalized versions are identical, no substantive differences
	if normalized1 == normalized2 {
		return DiffTypeNone
	}
	
	// If raw versions differ but normalized versions are same, it's formatting only
	if raw1 != raw2 && normalized1 == normalized2 {
		return DiffTypeFormatOnly
	}
	
	// Otherwise, substantive differences exist
	return DiffTypeSubstantive
}
