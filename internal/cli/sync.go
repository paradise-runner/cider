package cli

import (
	"fmt"

	"github.com/paradise-runner/cider/internal/convert"
	"github.com/paradise-runner/cider/internal/frontmatter"
	"github.com/paradise-runner/cider/internal/markdown"
	"github.com/paradise-runner/cider/internal/notes"
	syncpkg "github.com/paradise-runner/cider/internal/sync"
	"github.com/spf13/cobra"
)

var syncStrategy string

var syncCmd = &cobra.Command{
	Use:   "sync <file|directory>",
	Short: "Sync Markdown file(s) with their linked Apple Notes",
	Long: `Sync Markdown file(s) with their linked Apple Notes by performing a
block-level merge. Blocks are structural markdown elements such as headings,
paragraphs, lists, code blocks, and blockquotes.

The default strategy is "union", which keeps blocks from both sides. Use
--strategy to change the merge behavior:

  union         Include blocks from both local and remote (default)
  prefer-local  Keep local blocks, discard remote-only blocks
  prefer-remote Keep remote blocks, discard local-only blocks

The file(s) must have an apple_notes_id in their frontmatter.
If a directory is provided, all Markdown files (.md, .markdown) will be
processed recursively.`,
	Args: cobra.ExactArgs(1),
	RunE: runSync,
}

func init() {
	syncCmd.Flags().StringVar(&syncStrategy, "strategy", "union",
		`merge strategy: "union", "prefer-local", or "prefer-remote"`)
	rootCmd.AddCommand(syncCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	path := args[0]

	strategy, err := parseStrategy(syncStrategy)
	if err != nil {
		return err
	}

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
		if err := syncSingleFile(filePath, notesClient, strategy); err != nil {
			fmt.Printf("Error processing %s: %v\n", filePath, err)
			// Continue processing other files
		}
	}

	return nil
}

func syncSingleFile(filePath string, notesClient *notes.Client, strategy syncpkg.MergeStrategy) error {
	// Read markdown file
	fmt.Printf("Reading file: %s\n", filePath)
	markdownContent, err := markdown.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Extract ID from frontmatter (required for sync)
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

	// Convert remote HTML to Markdown
	remoteMarkdown, err := convert.HTMLToMarkdown(htmlContent)
	if err != nil {
		return fmt.Errorf("failed to convert HTML to markdown: %v", err)
	}

	// Strip frontmatter from local content for comparison
	localBody := frontmatter.Strip(markdownContent)

	// Run block-level sync
	result := syncpkg.Sync(localBody, remoteMarkdown, strategy)

	if !syncpkg.HasChanges(result.Changes) {
		fmt.Println("Already in sync.")
		return nil
	}

	fmt.Printf("Changes: %s\n", result.Summary)

	// Rebuild local file: existing frontmatter + merged body
	existingFrontmatter := frontmatter.Extract(markdownContent)
	var updatedContent string
	if existingFrontmatter != "" {
		updatedContent = existingFrontmatter + "\n" + result.Merged
	} else {
		updatedContent = result.Merged
	}

	// Write merged content back to local file
	err = markdown.WriteFile(filePath, updatedContent)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	fmt.Printf("File updated: %s\n", filePath)

	// Push merged content to Apple Notes
	fmt.Println("Updating Apple Note...")
	mergedHTML, err := convert.MarkdownToHTML(result.Merged)
	if err != nil {
		return fmt.Errorf("failed to convert merged markdown to HTML: %v", err)
	}

	err = notesClient.Update(noteID, mergedHTML)
	if err != nil {
		return fmt.Errorf("failed to update note: %v", err)
	}

	fmt.Printf("Note updated: %s\n", noteID)
	return nil
}

func parseStrategy(s string) (syncpkg.MergeStrategy, error) {
	switch s {
	case "union":
		return syncpkg.MergeUnion, nil
	case "prefer-local":
		return syncpkg.MergePreferLocal, nil
	case "prefer-remote":
		return syncpkg.MergePreferRemote, nil
	default:
		return 0, fmt.Errorf("unknown strategy %q: use \"union\", \"prefer-local\", or \"prefer-remote\"", s)
	}
}
