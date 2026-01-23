package markdown

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	tests := []struct {
		name      string
		fixture   string
		wantError bool
		contains  string
	}{
		{
			name:      "existing file",
			fixture:   "simple.md",
			wantError: false,
			contains:  "This is a simple paragraph",
		},
		{
			name:      "nonexistent file",
			fixture:   "does-not-exist.md",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("../../test/fixtures", tt.fixture)
			content, err := ReadFile(path)

			if tt.wantError {
				if err == nil {
					t.Errorf("ReadFile() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ReadFile() unexpected error: %v", err)
				}
				if tt.contains != "" && !contains(content, tt.contains) {
					t.Errorf("ReadFile() content should contain %q", tt.contains)
				}
			}
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		fileExists bool
		wantError  bool
	}{
		{
			name:       "write to existing file",
			content:    "# Updated Content\n\nThis is updated.",
			fileExists: true,
			wantError:  false,
		},
		{
			name:       "write to nonexistent file",
			content:    "# New Content",
			fileExists: false,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file if needed
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.md")

			if tt.fileExists {
				err := os.WriteFile(tmpFile, []byte("# Original"), 0644)
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
			}

			// Test WriteFile
			err := WriteFile(tmpFile, tt.content)

			if tt.wantError {
				if err == nil {
					t.Errorf("WriteFile() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("WriteFile() unexpected error: %v", err)
				}

				// Verify content was written
				written, err := os.ReadFile(tmpFile)
				if err != nil {
					t.Fatalf("failed to read written file: %v", err)
				}
				if string(written) != tt.content {
					t.Errorf("WriteFile() wrote %q, want %q", string(written), tt.content)
				}
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
