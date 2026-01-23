package cli

import (
	"fmt"

	"github.com/shakedlokits/stash/internal/convert"
	"github.com/shakedlokits/stash/internal/frontmatter"
	"github.com/shakedlokits/stash/internal/markdown"
	"github.com/shakedlokits/stash/internal/notes"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull <file>",
	Short: "Pull content from Apple Notes to a Markdown file",
	Long: `Pull content from Apple Notes to a Markdown file. The file must have an
apple_notes_id in its frontmatter. The note content will be converted to
Markdown and written to the file, preserving the frontmatter.`,
	Args: cobra.ExactArgs(1),
	RunE: runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Read markdown file
	fmt.Printf("Reading file: %s\n", filePath)
	markdownContent, err := markdown.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Error: File not found: %s", filePath)
	}

	// Extract ID from frontmatter (required for pull)
	noteID, err := frontmatter.GetAppleNotesID(markdownContent)
	if err != nil {
		return fmt.Errorf("Error: No apple_notes_id found in frontmatter")
	}

	// Initialize notes client
	notesClient := notes.NewClient()

	// Find note in Apple Notes
	fmt.Println("Searching for note...")
	found, err := notesClient.Find(noteID)
	if err != nil {
		return fmt.Errorf("Error: Failed to search for note: %v", err)
	}
	if !found {
		return fmt.Errorf("Error: Note not found in Apple Notes")
	}

	// Read note content
	fmt.Println("Reading note content...")
	htmlContent, err := notesClient.Read(noteID)
	if err != nil {
		return fmt.Errorf("Error: Failed to read note content: %v", err)
	}

	// Convert HTML to Markdown
	markdownBody, err := convert.HTMLToMarkdown(htmlContent)
	if err != nil {
		return fmt.Errorf("Error: Failed to convert HTML to markdown: %v", err)
	}

	// Extract existing frontmatter and rebuild file content
	existingFrontmatter := frontmatter.Extract(markdownContent)

	// Combine: frontmatter + empty line + markdown body
	var updatedContent string
	if existingFrontmatter != "" {
		updatedContent = existingFrontmatter + "\n" + markdownBody
	} else {
		updatedContent = markdownBody
	}

	// Write back to file
	err = markdown.WriteFile(filePath, updatedContent)
	if err != nil {
		return fmt.Errorf("Error: Failed to write to file: %v", err)
	}

	fmt.Printf("File updated: %s\n", filePath)

	return nil
}
