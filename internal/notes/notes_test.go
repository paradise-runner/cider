//go:build integration
// +build integration

package notes

import (
	"strings"
	"testing"
)

// These tests require macOS with Apple Notes and will create/delete actual notes
// Run with: go test -tags=integration ./internal/notes

func TestCreate(t *testing.T) {
	client := NewClient()

	htmlContent := "<p>Test note created by Go test</p>"
	noteID, err := client.Create(htmlContent)

	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	if !strings.HasPrefix(noteID, "x-coredata://") {
		t.Errorf("Create() returned invalid note ID: %s", noteID)
	}

	// Cleanup: delete the note
	defer func() {
		if err := client.Delete(noteID); err != nil {
			t.Logf("Warning: failed to cleanup test note: %v", err)
		}
	}()

	t.Logf("Created note: %s", noteID)
}

func TestRead(t *testing.T) {
	client := NewClient()

	// Create a test note
	htmlContent := "<p>Test content for reading</p>"
	noteID, err := client.Create(htmlContent)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	defer func() {
		if err := client.Delete(noteID); err != nil {
			t.Logf("Warning: failed to cleanup test note: %v", err)
		}
	}()

	// Read the note
	content, err := client.Read(noteID)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if !strings.Contains(content, "Test content for reading") {
		t.Errorf("Read() content does not match, got: %s", content)
	}

	t.Logf("Read note content: %s", content)
}

func TestUpdate(t *testing.T) {
	client := NewClient()

	// Create a test note
	htmlContent := "<p>Original content</p>"
	noteID, err := client.Create(htmlContent)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	defer func() {
		if err := client.Delete(noteID); err != nil {
			t.Logf("Warning: failed to cleanup test note: %v", err)
		}
	}()

	// Update the note
	updatedContent := "<p>Updated content</p>"
	err = client.Update(noteID, updatedContent)
	if err != nil {
		t.Fatalf("Update() error: %v", err)
	}

	// Verify the update
	content, err := client.Read(noteID)
	if err != nil {
		t.Fatalf("Read() after update error: %v", err)
	}

	if !strings.Contains(content, "Updated content") {
		t.Errorf("Update() did not update content, got: %s", content)
	}

	t.Logf("Updated note content: %s", content)
}

func TestDelete(t *testing.T) {
	client := NewClient()

	// Create a test note
	htmlContent := "<p>Test note to delete</p>"
	noteID, err := client.Create(htmlContent)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	// Delete the note
	err = client.Delete(noteID)
	if err != nil {
		t.Fatalf("Delete() error: %v", err)
	}

	// Verify it's deleted (in Recently Deleted, so Find should return false)
	found, err := client.Find(noteID)
	if err != nil {
		t.Fatalf("Find() after delete error: %v", err)
	}

	if found {
		t.Errorf("Find() should return false after delete, got true")
	}

	t.Logf("Deleted note: %s", noteID)
}

func TestFind(t *testing.T) {
	client := NewClient()

	// Test 1: Find existing note
	htmlContent := "<p>Test note for finding</p>"
	noteID, err := client.Create(htmlContent)
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}

	defer func() {
		if err := client.Delete(noteID); err != nil {
			t.Logf("Warning: failed to cleanup test note: %v", err)
		}
	}()

	found, err := client.Find(noteID)
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}

	if !found {
		t.Errorf("Find() should return true for existing note")
	}

	// Test 2: Find non-existent note
	fakeID := "x-coredata://fake/ICNote/p999999"
	found, err = client.Find(fakeID)
	if err != nil {
		t.Fatalf("Find() error for fake ID: %v", err)
	}

	if found {
		t.Errorf("Find() should return false for non-existent note")
	}

	t.Logf("Find tests passed")
}

func TestReadNonExistent(t *testing.T) {
	client := NewClient()

	fakeID := "x-coredata://fake/ICNote/p999999"
	_, err := client.Read(fakeID)

	if err == nil {
		t.Errorf("Read() should return error for non-existent note")
	}

	t.Logf("Read non-existent error (expected): %v", err)
}
