package cli

import (
	"fmt"

	"github.com/paradise-runner/cider/internal/convert"
	"github.com/paradise-runner/cider/internal/frontmatter"
	"github.com/paradise-runner/cider/internal/markdown"
	"github.com/paradise-runner/cider/internal/notes"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull <file|directory>",
	Short: "Pull content from Apple Notes to Markdown file(s)",
	Long: `Pull content from Apple Notes to Markdown file(s). The file(s) must have an
apple_notes_id in their frontmatter. The note content will be converted to
Markdown and written to the file, preserving the frontmatter.

If a directory is provided, all Markdown files (.md, .markdown) will be
processed recursively.`,
	Args: cobra.ExactArgs(1),
	RunE: runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	path := args[0]

	// Resolve path to list of markdown files
	files, err := resolvePaths(path)
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("Error: No markdown files found in %s", path)
	}

	// Initialize notes client
	notesClient := notes.NewClient()

	// Process each file
	for _, filePath := range files {
		if err := pullSingleFile(filePath, notesClient); err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
			// Continue processing other files
		}
	}

	return nil
}

func pullSingleFile(filePath string, notesClient *notes.Client) error {
	// Read markdown file
	fmt.Printf("Reading file: %s\n", filePath)
	markdownContent, err := markdown.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Extract ID from frontmatter (required for pull)
	noteID, err := frontmatter.GetAppleNotesID(markdownContent)
	if err != nil {
		return fmt.Errorf("no apple_notes_id found in frontmatter")
	}

	// Find note in Apple Notes
	fmt.Println("Searching for note...")
	found, err := notesClient.Find(noteID)
	if err != nil {
		return fmt.Errorf("failed to search for note: %v", err)
	}
	if !found {
		return fmt.Errorf("note not found in Apple Notes")
	}

	// Read note content
	fmt.Println("Reading note content...")
	htmlContent, err := notesClient.Read(noteID)
	if err != nil {
		return fmt.Errorf("failed to read note content: %v", err)
	}

	// Convert HTML to Markdown
	markdownBody, err := convert.HTMLToMarkdown(htmlContent)
	if err != nil {
		return fmt.Errorf("failed to convert HTML to markdown: %v", err)
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
		return fmt.Errorf("failed to write to file: %v", err)
	}

	fmt.Printf("File updated: %s\n", filePath)

	return nil
}
