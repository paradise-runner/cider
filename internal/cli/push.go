package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/paradise-runner/cider/internal/convert"
	"github.com/paradise-runner/cider/internal/frontmatter"
	"github.com/paradise-runner/cider/internal/markdown"
	"github.com/paradise-runner/cider/internal/notes"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push <file|directory>",
	Short: "Push Markdown file(s) to Apple Notes (create or update)",
	Long: `Push Markdown file(s) to Apple Notes. If the file has an apple_notes_id in its
frontmatter, the corresponding note will be updated. Otherwise, a new note
will be created and the frontmatter will be updated with the note ID.

If a directory is provided, all Markdown files (.md, .markdown) will be
processed recursively.`,
	Args: cobra.ExactArgs(1),
	RunE: runPush,
}

func init() {
	rootCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) error {
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
		if err := pushSingleFile(filePath, notesClient); err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
			// Continue processing other files
		}
	}

	return nil
}

func pushSingleFile(filePath string, notesClient *notes.Client) error {
	// Read markdown content from file
	fmt.Printf("Reading file: %s\n", filePath)
	markdownContent, err := markdown.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Extract ID from frontmatter (may not exist)
	noteID, _ := frontmatter.GetAppleNotesID(markdownContent)

	// Check if note exists in Apple Notes
	noteFound := false
	if noteID != "" {
		fmt.Println("Searching for note...")
		found, err := notesClient.Find(noteID)
		if err != nil {
			return fmt.Errorf("failed to search for note: %v", err)
		}
		noteFound = found
	}

	// If note not found, prompt user to create new
	if !noteFound {
		fmt.Println("Note not found in Apple Notes.")
		fmt.Print("Create new note? (y/n) ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user input: %v", err)
		}

		response = strings.TrimSpace(response)
		if !strings.EqualFold(response, "y") && !strings.EqualFold(response, "yes") {
			fmt.Println("Operation cancelled")
			return nil
		}

		// Strip frontmatter and convert to HTML
		fmt.Println("Creating note...")
		body := frontmatter.Strip(markdownContent)
		htmlContent, err := convert.MarkdownToHTML(body)
		if err != nil {
			return fmt.Errorf("failed to convert markdown to HTML: %v", err)
		}

		// Create new note
		newNoteID, err := notesClient.Create(htmlContent)
		if err != nil {
			return fmt.Errorf("failed to create note: %v", err)
		}

		// Update frontmatter with new ID
		updatedContent := frontmatter.UpdateAppleNotesID(markdownContent, newNoteID)
		err = markdown.WriteFile(filePath, updatedContent)
		if err != nil {
			return fmt.Errorf("failed to update file: %v", err)
		}

		fmt.Printf("Note created: %s\n", newNoteID)
	} else {
		// Note exists, update it
		fmt.Println("Updating note...")
		body := frontmatter.Strip(markdownContent)
		htmlContent, err := convert.MarkdownToHTML(body)
		if err != nil {
			return fmt.Errorf("failed to convert markdown to HTML: %v", err)
		}

		err = notesClient.Update(noteID, htmlContent)
		if err != nil {
			return fmt.Errorf("failed to update note: %v", err)
		}

		fmt.Printf("Note updated: %s\n", noteID)
	}

	return nil
}
