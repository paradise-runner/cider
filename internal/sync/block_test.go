package sync

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   []Block
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "whitespace only",
			input: "   \n\n  \n",
			want:  nil,
		},
		{
			name:  "single heading",
			input: "# Hello",
			want: []Block{
				{Type: BlockHeading, Content: "# Hello"},
			},
		},
		{
			name:  "heading levels",
			input: "# H1\n\n## H2\n\n### H3\n\n#### H4\n\n##### H5\n\n###### H6",
			want: []Block{
				{Type: BlockHeading, Content: "# H1"},
				{Type: BlockHeading, Content: "## H2"},
				{Type: BlockHeading, Content: "### H3"},
				{Type: BlockHeading, Content: "#### H4"},
				{Type: BlockHeading, Content: "##### H5"},
				{Type: BlockHeading, Content: "###### H6"},
			},
		},
		{
			name:  "single paragraph",
			input: "This is a paragraph.",
			want: []Block{
				{Type: BlockParagraph, Content: "This is a paragraph."},
			},
		},
		{
			name:  "multi-line paragraph",
			input: "Line one\nLine two\nLine three",
			want: []Block{
				{Type: BlockParagraph, Content: "Line one\nLine two\nLine three"},
			},
		},
		{
			name:  "two paragraphs",
			input: "First paragraph.\n\nSecond paragraph.",
			want: []Block{
				{Type: BlockParagraph, Content: "First paragraph."},
				{Type: BlockParagraph, Content: "Second paragraph."},
			},
		},
		{
			name:  "heading and paragraph",
			input: "# Title\n\nSome text here.",
			want: []Block{
				{Type: BlockHeading, Content: "# Title"},
				{Type: BlockParagraph, Content: "Some text here."},
			},
		},
		{
			name:  "fenced code block with backticks",
			input: "```go\nfmt.Println(\"hello\")\n```",
			want: []Block{
				{Type: BlockFencedCode, Content: "```go\nfmt.Println(\"hello\")\n```"},
			},
		},
		{
			name:  "fenced code block with tildes",
			input: "~~~\nsome code\n~~~",
			want: []Block{
				{Type: BlockFencedCode, Content: "~~~\nsome code\n~~~"},
			},
		},
		{
			name:  "fenced code with blank lines inside",
			input: "```\nline 1\n\nline 3\n```",
			want: []Block{
				{Type: BlockFencedCode, Content: "```\nline 1\n\nline 3\n```"},
			},
		},
		{
			name:  "unordered list with dashes",
			input: "- Item 1\n- Item 2\n- Item 3",
			want: []Block{
				{Type: BlockList, Content: "- Item 1\n- Item 2\n- Item 3"},
			},
		},
		{
			name:  "unordered list with asterisks",
			input: "* Item 1\n* Item 2",
			want: []Block{
				{Type: BlockList, Content: "* Item 1\n* Item 2"},
			},
		},
		{
			name:  "ordered list",
			input: "1. First\n2. Second\n3. Third",
			want: []Block{
				{Type: BlockList, Content: "1. First\n2. Second\n3. Third"},
			},
		},
		{
			name:  "list with continuation",
			input: "- Item 1\n  continued\n- Item 2",
			want: []Block{
				{Type: BlockList, Content: "- Item 1\n  continued\n- Item 2"},
			},
		},
		{
			name:  "blockquote single line",
			input: "> A quote",
			want: []Block{
				{Type: BlockBlockquote, Content: "> A quote"},
			},
		},
		{
			name:  "blockquote multiple lines",
			input: "> Line 1\n> Line 2\n> Line 3",
			want: []Block{
				{Type: BlockBlockquote, Content: "> Line 1\n> Line 2\n> Line 3"},
			},
		},
		{
			name:  "horizontal rule dashes",
			input: "---",
			want: []Block{
				{Type: BlockHorizontalRule, Content: "---"},
			},
		},
		{
			name:  "horizontal rule asterisks",
			input: "***",
			want: []Block{
				{Type: BlockHorizontalRule, Content: "***"},
			},
		},
		{
			name:  "horizontal rule underscores",
			input: "___",
			want: []Block{
				{Type: BlockHorizontalRule, Content: "___"},
			},
		},
		{
			name: "full document",
			input: `# My Document

This is the introduction paragraph.

## Section One

Some content here.

- Item A
- Item B
- Item C

## Section Two

> A wise quote

Final thoughts.`,
			want: []Block{
				{Type: BlockHeading, Content: "# My Document"},
				{Type: BlockParagraph, Content: "This is the introduction paragraph."},
				{Type: BlockHeading, Content: "## Section One"},
				{Type: BlockParagraph, Content: "Some content here."},
				{Type: BlockList, Content: "- Item A\n- Item B\n- Item C"},
				{Type: BlockHeading, Content: "## Section Two"},
				{Type: BlockBlockquote, Content: "> A wise quote"},
				{Type: BlockParagraph, Content: "Final thoughts."},
			},
		},
		{
			name:  "code block between paragraphs",
			input: "Before code.\n\n```\ncode here\n```\n\nAfter code.",
			want: []Block{
				{Type: BlockParagraph, Content: "Before code."},
				{Type: BlockFencedCode, Content: "```\ncode here\n```"},
				{Type: BlockParagraph, Content: "After code."},
			},
		},
		{
			name:  "loose list with blank lines",
			input: "- Item 1\n\n- Item 2\n\n- Item 3",
			want: []Block{
				{Type: BlockList, Content: "- Item 1\n\n- Item 2\n\n- Item 3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.input)

			if len(got) != len(tt.want) {
				t.Fatalf("Parse() returned %d blocks, want %d\nGot: %v", len(got), len(tt.want), got)
			}

			for i := range got {
				if got[i].Type != tt.want[i].Type {
					t.Errorf("block[%d].Type = %d, want %d", i, got[i].Type, tt.want[i].Type)
				}
				if got[i].Content != tt.want[i].Content {
					t.Errorf("block[%d].Content = %q, want %q", i, got[i].Content, tt.want[i].Content)
				}
			}
		})
	}
}

func TestRender(t *testing.T) {
	tests := []struct {
		name   string
		blocks []Block
		want   string
	}{
		{
			name:   "empty",
			blocks: nil,
			want:   "",
		},
		{
			name: "single block",
			blocks: []Block{
				{Type: BlockHeading, Content: "# Title"},
			},
			want: "# Title\n",
		},
		{
			name: "multiple blocks",
			blocks: []Block{
				{Type: BlockHeading, Content: "# Title"},
				{Type: BlockParagraph, Content: "A paragraph."},
				{Type: BlockList, Content: "- Item 1\n- Item 2"},
			},
			want: "# Title\n\nA paragraph.\n\n- Item 1\n- Item 2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Render(tt.blocks)
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseRenderRoundTrip(t *testing.T) {
	// A well-formed document should survive parse→render without meaningful
	// content loss (though exact whitespace may differ).
	input := `# Title

First paragraph with some text.

## Subsection

- Item 1
- Item 2
- Item 3

> A blockquote line
> continues here

` + "```" + `python
def hello():
    print("world")
` + "```" + `

Final paragraph.
`

	blocks := Parse(input)
	output := Render(blocks)
	reparsed := Parse(output)

	if len(blocks) != len(reparsed) {
		t.Fatalf("round-trip block count: got %d, original %d", len(reparsed), len(blocks))
	}

	for i := range blocks {
		if blocks[i].Type != reparsed[i].Type {
			t.Errorf("round-trip block[%d].Type = %d, want %d", i, reparsed[i].Type, blocks[i].Type)
		}
		if blocks[i].Content != reparsed[i].Content {
			t.Errorf("round-trip block[%d].Content = %q, want %q", i, reparsed[i].Content, blocks[i].Content)
		}
	}
}
