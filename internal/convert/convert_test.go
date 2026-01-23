package convert

import (
	"os"
	"path/filepath"
	"strings"
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

func TestMarkdownToHTML(t *testing.T) {
	// Skip if pandoc not available
	if err := CheckPandocAvailable(); err != nil {
		t.Skip("pandoc not available:", err)
	}

	tests := []struct {
		name     string
		fixture  string
		contains []string
	}{
		{
			name:    "simple markdown",
			fixture: "simple.md",
			contains: []string{
				"<p>This is a simple paragraph.</p>",
				"<p>This is another paragraph.</p>",
			},
		},
		{
			name:    "formatted markdown",
			fixture: "formatted.md",
			contains: []string{
				"<strong>",
				"<em>",
				"<code>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			markdown := loadFixture(t, tt.fixture)
			html, err := MarkdownToHTML(markdown)

			if err != nil {
				t.Fatalf("MarkdownToHTML() error: %v", err)
			}

			for _, s := range tt.contains {
				if !strings.Contains(html, s) {
					t.Errorf("MarkdownToHTML() should contain %q, got:\n%s", s, html)
				}
			}
		})
	}
}

func TestHTMLToMarkdown(t *testing.T) {
	// Skip if pandoc not available
	if err := CheckPandocAvailable(); err != nil {
		t.Skip("pandoc not available:", err)
	}

	tests := []struct {
		name        string
		fixture     string
		contains    []string
		notContains []string
	}{
		{
			name:    "simple html",
			fixture: "simple.html",
			contains: []string{
				"This is a simple paragraph",
			},
			notContains: []string{"<p>", "<div>"},
		},
		{
			name:    "formatted html",
			fixture: "formatted.html",
			contains: []string{
				"**bold**",
				"*italic*",
				"`inline code`",
			},
			notContains: []string{"<strong>", "<em>", "<code>"},
		},
		{
			name:    "apple notes sample",
			fixture: "apple_notes_sample.html",
			contains: []string{
				"# ", // Should have H1
			},
			notContains: []string{
				`<div><b><span style="font-size: 24px">`, // Apple title pattern should be converted
				"&nbsp;",                                 // nbsp should be removed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := loadFixture(t, tt.fixture)
			markdown, err := HTMLToMarkdown(html)

			if err != nil {
				t.Fatalf("HTMLToMarkdown() error: %v", err)
			}

			for _, s := range tt.contains {
				if !strings.Contains(markdown, s) {
					t.Errorf("HTMLToMarkdown() should contain %q, got:\n%s", s, markdown)
				}
			}

			for _, s := range tt.notContains {
				if strings.Contains(markdown, s) {
					t.Errorf("HTMLToMarkdown() should NOT contain %q, got:\n%s", s, markdown)
				}
			}
		})
	}
}

func TestConvertAppleTitleToH1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "apple notes title",
			input: `<div><b><span style="font-size: 24px">My Title</span></b></div>`,
			want:  `<h1>My Title</h1>`,
		},
		{
			name:  "multiple titles",
			input: `<div><b><span style="font-size: 24px">Title 1</span></b></div><p>Content</p><div><b><span style="font-size: 24px">Title 2</span></b></div>`,
			want:  `<h1>Title 1</h1><p>Content</p><h1>Title 2</h1>`,
		},
		{
			name:  "no title",
			input: `<p>Just a paragraph</p>`,
			want:  `<p>Just a paragraph</p>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertAppleTitleToH1(tt.input)
			if got != tt.want {
				t.Errorf("convertAppleTitleToH1() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRemoveNbsp(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "remove nbsp line",
			input: "Line 1\n&nbsp;\nLine 2",
			want:  "Line 1\nLine 2",
		},
		{
			name:  "keep nbsp in text",
			input: "Line with&nbsp;nbsp inside",
			want:  "Line with&nbsp;nbsp inside",
		},
		{
			name:  "no nbsp",
			input: "Line 1\nLine 2",
			want:  "Line 1\nLine 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeNbsp(tt.input)
			if got != tt.want {
				t.Errorf("removeNbsp() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTrimTrailingWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "trim spaces",
			input: "Line 1   \nLine 2\t\t\nLine 3",
			want:  "Line 1\nLine 2\nLine 3",
		},
		{
			name:  "no trailing whitespace",
			input: "Line 1\nLine 2",
			want:  "Line 1\nLine 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimTrailingWhitespace(tt.input)
			if got != tt.want {
				t.Errorf("trimTrailingWhitespace() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCollapseBlankLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "collapse multiple blank lines",
			input: "Line 1\n\n\n\nLine 2",
			want:  "Line 1\n\nLine 2",
		},
		{
			name:  "keep single blank line",
			input: "Line 1\n\nLine 2",
			want:  "Line 1\n\nLine 2",
		},
		{
			name:  "no blank lines",
			input: "Line 1\nLine 2",
			want:  "Line 1\nLine 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := collapseBlankLines(tt.input)
			if got != tt.want {
				t.Errorf("collapseBlankLines() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCheckPandocAvailable(t *testing.T) {
	err := CheckPandocAvailable()
	// This test just checks if the function runs without panic
	// The actual availability depends on the system
	if err != nil {
		t.Logf("Pandoc not available: %v", err)
	}
}
