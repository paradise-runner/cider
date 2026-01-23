package frontmatter

import (
	"os"
	"path/filepath"
	"testing"
)

func loadFixture(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join("../../test/fixtures", name)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", name, err)
	}
	return string(content)
}

func TestStrip(t *testing.T) {
	tests := []struct {
		name        string
		fixture     string
		contains    []string // strings that should be in output
		notContains []string // strings that should NOT be in output
	}{
		{
			name:        "with apple id",
			fixture:     "with_apple_id.md",
			contains:    []string{"# Synced Note", "This note has been synced"},
			notContains: []string{"---", "apple_notes_id:", "title:", "tags:"},
		},
		{
			name:        "with frontmatter no id",
			fixture:     "with_frontmatter_no_id.md",
			contains:    []string{"# My Cool Note"},
			notContains: []string{"---", "title:", "tags:", "created:"},
		},
		{
			name:        "no frontmatter",
			fixture:     "no_frontmatter.md",
			contains:    []string{"# My Cool Note", "## Features"},
			notContains: []string{"---"},
		},
		{
			name:        "empty file",
			fixture:     "empty.md",
			contains:    []string{},
			notContains: []string{"---"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := loadFixture(t, tt.fixture)
			result := Strip(content)

			for _, s := range tt.contains {
				if !contains(result, s) {
					t.Errorf("Strip() result should contain %q, got:\n%s", s, result)
				}
			}

			for _, s := range tt.notContains {
				if contains(result, s) {
					t.Errorf("Strip() result should NOT contain %q, got:\n%s", s, result)
				}
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		contains []string
		isEmpty  bool
	}{
		{
			name:     "with apple id",
			fixture:  "with_apple_id.md",
			contains: []string{"---", "apple_notes_id:", "title:", "tags:"},
		},
		{
			name:     "with frontmatter no id",
			fixture:  "with_frontmatter_no_id.md",
			contains: []string{"---", "title:", "tags:", "created:"},
		},
		{
			name:    "no frontmatter",
			fixture: "no_frontmatter.md",
			isEmpty: true,
		},
		{
			name:    "empty file",
			fixture: "empty.md",
			isEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := loadFixture(t, tt.fixture)
			result := Extract(content)

			if tt.isEmpty {
				if result != "" {
					t.Errorf("Extract() should return empty string, got: %s", result)
				}
				return
			}

			for _, s := range tt.contains {
				if !contains(result, s) {
					t.Errorf("Extract() result should contain %q, got:\n%s", s, result)
				}
			}
		})
	}
}

func TestGetAppleNotesID(t *testing.T) {
	tests := []struct {
		name      string
		fixture   string
		wantID    string
		wantError bool
	}{
		{
			name:      "with apple id",
			fixture:   "with_apple_id.md",
			wantID:    "x-coredata://0C0AFA7A-CCCF-4DEB-A67D-424B29956D25/ICNote/p1970",
			wantError: false,
		},
		{
			name:      "with frontmatter no id",
			fixture:   "with_frontmatter_no_id.md",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "no frontmatter",
			fixture:   "no_frontmatter.md",
			wantID:    "",
			wantError: true,
		},
		{
			name:      "empty file",
			fixture:   "empty.md",
			wantID:    "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := loadFixture(t, tt.fixture)
			id, err := GetAppleNotesID(content)

			if tt.wantError {
				if err == nil {
					t.Errorf("GetAppleNotesID() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("GetAppleNotesID() unexpected error: %v", err)
				}
			}

			if id != tt.wantID {
				t.Errorf("GetAppleNotesID() = %q, want %q", id, tt.wantID)
			}
		})
	}
}

func TestUpdateAppleNotesID(t *testing.T) {
	tests := []struct {
		name     string
		fixture  string
		newID    string
		contains []string
	}{
		{
			name:    "with apple id - update",
			fixture: "with_apple_id.md",
			newID:   "x-coredata://NEW-ID/ICNote/p999",
			contains: []string{
				"---",
				"apple_notes_id: x-coredata://NEW-ID/ICNote/p999",
				"title: Synced Note",
				"tags: synced",
				"# Synced Note",
			},
		},
		{
			name:    "with frontmatter no id - add",
			fixture: "with_frontmatter_no_id.md",
			newID:   "x-coredata://ADDED-ID/ICNote/p123",
			contains: []string{
				"---",
				"apple_notes_id: x-coredata://ADDED-ID/ICNote/p123",
				"title: My Note Title",
				"tags: personal, ideas",
				"# My Cool Note",
			},
		},
		{
			name:    "no frontmatter - create",
			fixture: "no_frontmatter.md",
			newID:   "x-coredata://CREATED-ID/ICNote/p456",
			contains: []string{
				"---",
				"apple_notes_id: x-coredata://CREATED-ID/ICNote/p456",
				"# My Cool Note",
				"## Features",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := loadFixture(t, tt.fixture)
			result := UpdateAppleNotesID(content, tt.newID)

			for _, s := range tt.contains {
				if !contains(result, s) {
					t.Errorf("UpdateAppleNotesID() result should contain %q, got:\n%s", s, result)
				}
			}

			// Verify the ID is actually extractable
			id, err := GetAppleNotesID(result)
			if err != nil {
				t.Errorf("GetAppleNotesID() after update failed: %v", err)
			}
			if id != tt.newID {
				t.Errorf("GetAppleNotesID() after update = %q, want %q", id, tt.newID)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
