package cli

import (
	"testing"
)

func TestNormalizeContent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no changes needed",
			input: "Hello\nWorld",
			want:  "Hello\nWorld",
		},
		{
			name:  "trailing whitespace",
			input: "Hello   \nWorld  ",
			want:  "Hello\nWorld",
		},
		{
			name:  "leading and trailing blank lines",
			input: "\n\nHello\nWorld\n\n",
			want:  "Hello\nWorld",
		},
		{
			name:  "windows line endings",
			input: "Hello\r\nWorld\r\n",
			want:  "Hello\nWorld",
		},
		{
			name:  "mac line endings",
			input: "Hello\rWorld\r",
			want:  "Hello\nWorld",
		},
		{
			name:  "mixed whitespace",
			input: "  Hello   \t\n  World\t  \n\n",
			want:  "Hello\n  World",
		},
		{
			name:  "empty content",
			input: "",
			want:  "",
		},
		{
			name:  "only whitespace",
			input: "   \n\t\n   ",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeContent(tt.input)
			if got != tt.want {
				t.Errorf("normalizeContent() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAnalyzeDifferences(t *testing.T) {
	tests := []struct {
		name         string
		raw1         string
		raw2         string
		normalized1  string
		normalized2  string
		want         DiffType
	}{
		{
			name:         "identical content",
			raw1:         "Hello\nWorld",
			raw2:         "Hello\nWorld",
			normalized1:  "Hello\nWorld",
			normalized2:  "Hello\nWorld",
			want:         DiffTypeNone,
		},
		{
			name:         "formatting only - trailing whitespace",
			raw1:         "Hello  \nWorld",
			raw2:         "Hello\nWorld",
			normalized1:  "Hello\nWorld",
			normalized2:  "Hello\nWorld",
			want:         DiffTypeFormatOnly,
		},
		{
			name:         "formatting only - line endings",
			raw1:         "Hello\r\nWorld",
			raw2:         "Hello\nWorld",
			normalized1:  "Hello\nWorld",
			normalized2:  "Hello\nWorld",
			want:         DiffTypeFormatOnly,
		},
		{
			name:         "substantive difference",
			raw1:         "Hello\nWorld",
			raw2:         "Hello\nUniverse",
			normalized1:  "Hello\nWorld",
			normalized2:  "Hello\nUniverse",
			want:         DiffTypeSubstantive,
		},
		{
			name:         "added content",
			raw1:         "Hello",
			raw2:         "Hello\nWorld",
			normalized1:  "Hello",
			normalized2:  "Hello\nWorld",
			want:         DiffTypeSubstantive,
		},
		{
			name:         "removed content",
			raw1:         "Hello\nWorld",
			raw2:         "Hello",
			normalized1:  "Hello\nWorld",
			normalized2:  "Hello",
			want:         DiffTypeSubstantive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := analyzeDifferences(tt.raw1, tt.raw2, tt.normalized1, tt.normalized2)
			if got != tt.want {
				t.Errorf("analyzeDifferences() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test that DiffType constants have expected values
func TestDiffTypeValues(t *testing.T) {
	if DiffTypeNone != 0 {
		t.Errorf("DiffTypeNone should be 0, got %d", DiffTypeNone)
	}
	if DiffTypeFormatOnly != 1 {
		t.Errorf("DiffTypeFormatOnly should be 1, got %d", DiffTypeFormatOnly)
	}
	if DiffTypeSubstantive != 2 {
		t.Errorf("DiffTypeSubstantive should be 2, got %d", DiffTypeSubstantive)
	}
}
