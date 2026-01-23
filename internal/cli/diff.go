package cli

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/shakedlokits/stash/internal/convert"
	"github.com/shakedlokits/stash/internal/frontmatter"
	"github.com/shakedlokits/stash/internal/markdown"
	"github.com/shakedlokits/stash/internal/notes"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff <file>",
	Short: "Show differences between a Markdown file and its linked Apple Note",
	Long: `Show differences between a Markdown file and its linked Apple Note.
The file must have an apple_notes_id in its frontmatter. Shows a unified
diff with the Apple Note content as the old version and the file as new.`,
	Args:          cobra.ExactArgs(1),
	SilenceUsage:  true, // Don't show usage on errors (diff returns 1 when files differ)
	SilenceErrors: true, // We handle errors ourselves
	RunE:          runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Read markdown file
	markdownContent, err := markdown.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: File not found: %s\n", filePath)
		return &exitError{code: 2}
	}

	// Extract ID from frontmatter (required for diff)
	noteID, err := frontmatter.GetAppleNotesID(markdownContent)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: No apple_notes_id found in frontmatter\n")
		return &exitError{code: 2}
	}

	// Initialize notes client
	notesClient := notes.NewClient()

	// Find note in Apple Notes
	found, err := notesClient.Find(noteID)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: Failed to search for note: %v\n", err)
		return &exitError{code: 2}
	}
	if !found {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: Note not found in Apple Notes\n")
		return &exitError{code: 2}
	}

	// Read note content from Apple Notes
	htmlContent, err := notesClient.Read(noteID)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: Failed to read note content: %v\n", err)
		return &exitError{code: 2}
	}

	// Convert remote HTML to markdown
	remoteMarkdown, err := convert.HTMLToMarkdown(htmlContent)
	if err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error: Failed to convert HTML to markdown: %v\n", err)
		return &exitError{code: 2}
	}

	// Strip frontmatter from local content for fair comparison
	localMarkdown := frontmatter.Strip(markdownContent)

	// Show diff: remote (old) vs local (new)
	err = showDiff("Apple Notes", filePath, remoteMarkdown, localMarkdown)
	if err != nil {
		// diff returns exit code 1 if files differ, which is expected
		if exitErr, ok := err.(*exitError); ok && exitErr.code == 1 {
			return exitErr
		}
		return fmt.Errorf("Error: Failed to generate diff: %v", err)
	}

	return nil
}

// showDiff uses the system diff command to show differences
func showDiff(label1, label2, content1, content2 string) error {
	// Write content to temp files
	tmpFile1, err := writeTempFile(content1)
	if err != nil {
		return err
	}
	tmpFile2, err := writeTempFile(content2)
	if err != nil {
		return err
	}

	cmd := exec.Command("diff", "-u", "--label", label1, "--label", label2, tmpFile1, tmpFile2)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()

	// Print the diff output
	fmt.Print(stdout.String())

	// diff returns exit code 0 if same, 1 if different, 2 if error
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			if exitCode == 1 {
				// Files differ (expected)
				return &exitError{code: 1}
			}
			return fmt.Errorf("diff failed with exit code %d", exitCode)
		}
		return err
	}

	return nil
}

func writeTempFile(content string) (string, error) {
	tmpFile, err := exec.Command("mktemp").Output()
	if err != nil {
		return "", err
	}
	tmpPath := string(bytes.TrimSpace(tmpFile))

	writeCmd := exec.Command("sh", "-c", fmt.Sprintf("cat > %s", tmpPath))
	writeCmd.Stdin = bytes.NewBufferString(content)
	err = writeCmd.Run()
	if err != nil {
		return "", err
	}

	return tmpPath, nil
}

// exitError is used to return a specific exit code
type exitError struct {
	code int
}

func (e *exitError) Error() string {
	return fmt.Sprintf("exit code %d", e.code)
}
