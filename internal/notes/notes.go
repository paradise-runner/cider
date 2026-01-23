package notes

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client handles Apple Notes operations via AppleScript
type Client struct{}

// NewClient creates a new Apple Notes client
func NewClient() *Client {
	return &Client{}
}

// Create creates a new note with HTML content and returns the note ID
func (c *Client) Create(htmlContent string) (string, error) {
	// Escape double quotes for AppleScript
	escapedContent := strings.ReplaceAll(htmlContent, `"`, `\"`)

	script := fmt.Sprintf(`
tell application "Notes"
  try
    set newNote to make new note with properties {body:"%s"}
    return id of newNote
  on error errMsg
    error errMsg
  end try
end tell
`, escapedContent)

	result, err := c.runAppleScript(script)
	if err != nil {
		return "", fmt.Errorf("failed to create note: %w", err)
	}

	result = strings.TrimSpace(result)

	// Verify it looks like a valid note ID
	if !strings.HasPrefix(result, "x-coredata://") {
		return "", fmt.Errorf("invalid note ID returned: %s", result)
	}

	return result, nil
}

// Read reads the HTML body content of a note by ID
func (c *Client) Read(noteID string) (string, error) {
	script := fmt.Sprintf(`
tell application "Notes"
  try
    set deletedNotesFolder to folder "Recently Deleted"
    set theNote to first note whose id is "%s"
    set theFolder to container of theNote
    
    if theFolder is equal to deletedNotesFolder then
      return ""
    end if
    
    return body of theNote
  on error
    return ""
  end try
end tell
`, noteID)

	result, err := c.runAppleScript(script)
	if err != nil {
		return "", fmt.Errorf("failed to read note: %w", err)
	}

	result = strings.TrimSpace(result)

	if result == "" {
		return "", fmt.Errorf("note not found or is in Recently Deleted")
	}

	return result, nil
}

// Update updates the HTML body content of an existing note
func (c *Client) Update(noteID, htmlContent string) error {
	// Escape double quotes for AppleScript
	escapedContent := strings.ReplaceAll(htmlContent, `"`, `\"`)

	script := fmt.Sprintf(`
tell application "Notes"
  try
    set deletedNotesFolder to folder "Recently Deleted"
    set theNote to first note whose id is "%s"
    set theFolder to container of theNote
    
    if theFolder is equal to deletedNotesFolder then
      error "Note is in Recently Deleted"
    end if
    
    set body of theNote to "%s"
    return "%s"
  on error errMsg
    error errMsg
  end try
end tell
`, noteID, escapedContent, noteID)

	result, err := c.runAppleScript(script)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	result = strings.TrimSpace(result)

	if result != noteID {
		return fmt.Errorf("update verification failed: expected %s, got %s", noteID, result)
	}

	return nil
}

// Delete moves a note to Recently Deleted folder
func (c *Client) Delete(noteID string) error {
	script := fmt.Sprintf(`
tell application "Notes"
  try
    set deletedNotesFolder to folder "Recently Deleted"
    set theNote to first note whose id is "%s"
    set theFolder to container of theNote
    
    if theFolder is equal to deletedNotesFolder then
      error "Note already in Recently Deleted"
    end if
    
    delete theNote
    return "%s"
  on error errMsg
    error errMsg
  end try
end tell
`, noteID, noteID)

	result, err := c.runAppleScript(script)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	result = strings.TrimSpace(result)

	if result != noteID {
		return fmt.Errorf("delete verification failed: expected %s, got %s", noteID, result)
	}

	return nil
}

// Find checks if a note exists (and is not in Recently Deleted)
func (c *Client) Find(noteID string) (bool, error) {
	script := fmt.Sprintf(`
tell application "Notes"
  try
    set deletedNotesFolder to folder "Recently Deleted"
    set theNote to first note whose id is "%s"
    set theFolder to container of theNote
    
    if theFolder is equal to deletedNotesFolder then
      return ""
    end if
    
    return id of theNote
  on error
    return ""
  end try
end tell
`, noteID)

	result, err := c.runAppleScript(script)
	if err != nil {
		return false, fmt.Errorf("failed to find note: %w", err)
	}

	result = strings.TrimSpace(result)

	return result != "", nil
}

// runAppleScript executes an AppleScript and returns the output
func (c *Client) runAppleScript(script string) (string, error) {
	cmd := exec.Command("osascript")
	cmd.Stdin = strings.NewReader(script)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		stderrStr := stderr.String()
		if stderrStr != "" {
			return "", fmt.Errorf("osascript error: %s", stderrStr)
		}
		return "", fmt.Errorf("osascript failed: %w", err)
	}

	return stdout.String(), nil
}
