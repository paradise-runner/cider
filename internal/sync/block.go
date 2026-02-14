package sync

import (
	"regexp"
	"strings"
)

// BlockType identifies the kind of markdown structural element
type BlockType int

const (
	BlockParagraph      BlockType = iota // Default text block
	BlockHeading                         // ATX heading (# through ######)
	BlockFencedCode                      // Fenced code block (``` or ~~~)
	BlockBlockquote                      // Blockquote (> prefix)
	BlockList                            // Ordered or unordered list
	BlockHorizontalRule                  // Thematic break (---, ***, ___)
)

// Block represents a single structural element of a markdown document
type Block struct {
	Type    BlockType
	Content string // raw markdown text of the block (without surrounding blank lines)
}

var (
	headingRegex = regexp.MustCompile(`^#{1,6}(\s|$)`)
	hrDashRegex  = regexp.MustCompile(`^\s{0,3}(-\s*){3,}$`)
	hrStarRegex  = regexp.MustCompile(`^\s{0,3}(\*\s*){3,}$`)
	hrUnderRegex = regexp.MustCompile(`^\s{0,3}(_\s*){3,}$`)
	ulRegex      = regexp.MustCompile(`^\s{0,3}[-*+]\s`)
	olRegex      = regexp.MustCompile(`^\s{0,3}\d{1,9}[.)]\s`)
	fenceRegex   = regexp.MustCompile("^\\s{0,3}(`{3,}|~{3,})")
)

// Parse splits markdown content into a sequence of blocks.
// It recognises headings, fenced code blocks, blockquotes, lists,
// horizontal rules, and paragraphs. Blank lines between blocks are
// consumed as separators and not represented in the output.
func Parse(content string) []Block {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	lines := strings.Split(content, "\n")
	var blocks []Block
	i := 0

	for i < len(lines) {
		// Skip blank lines (separators between blocks)
		if isBlankLine(lines[i]) {
			i++
			continue
		}

		// Fenced code block
		if m := fenceRegex.FindString(lines[i]); m != "" {
			block, next := parseFencedCode(lines, i, strings.TrimSpace(m))
			blocks = append(blocks, block)
			i = next
			continue
		}

		// ATX heading
		if headingRegex.MatchString(lines[i]) {
			blocks = append(blocks, Block{Type: BlockHeading, Content: lines[i]})
			i++
			continue
		}

		// Horizontal rule (must be checked before list, since --- could be either)
		if isHorizontalRule(lines[i]) {
			blocks = append(blocks, Block{Type: BlockHorizontalRule, Content: lines[i]})
			i++
			continue
		}

		// Blockquote
		if strings.HasPrefix(strings.TrimLeft(lines[i], " "), ">") {
			block, next := parseBlockquote(lines, i)
			blocks = append(blocks, block)
			i = next
			continue
		}

		// List
		if isListItem(lines[i]) {
			block, next := parseList(lines, i)
			blocks = append(blocks, block)
			i = next
			continue
		}

		// Default: paragraph
		block, next := parseParagraph(lines, i)
		blocks = append(blocks, block)
		i = next
	}

	return blocks
}

// Render reassembles a sequence of blocks into markdown text,
// separating each block with a blank line.
func Render(blocks []Block) string {
	if len(blocks) == 0 {
		return ""
	}
	parts := make([]string, len(blocks))
	for i, b := range blocks {
		parts[i] = b.Content
	}
	return strings.Join(parts, "\n\n") + "\n"
}

// --- helpers -----------------------------------------------------------------

func isBlankLine(line string) bool {
	return strings.TrimSpace(line) == ""
}

func isHorizontalRule(line string) bool {
	return hrDashRegex.MatchString(line) || hrStarRegex.MatchString(line) || hrUnderRegex.MatchString(line)
}

func isListItem(line string) bool {
	return ulRegex.MatchString(line) || olRegex.MatchString(line)
}

func parseFencedCode(lines []string, start int, fence string) (Block, int) {
	// fence is the trimmed opening marker (e.g. "```" or "~~~")
	fenceChar := fence[0]
	fenceLen := len(fence)
	var collected []string
	collected = append(collected, lines[start])

	i := start + 1
	for i < len(lines) {
		collected = append(collected, lines[i])
		trimmed := strings.TrimSpace(lines[i])
		// Closing fence: same character, at least as many, and nothing else
		if len(trimmed) >= fenceLen && trimmed == strings.Repeat(string(fenceChar), len(trimmed)) {
			i++
			break
		}
		i++
	}

	return Block{
		Type:    BlockFencedCode,
		Content: strings.Join(collected, "\n"),
	}, i
}

func parseBlockquote(lines []string, start int) (Block, int) {
	var collected []string
	i := start
	for i < len(lines) {
		if strings.HasPrefix(strings.TrimLeft(lines[i], " "), ">") {
			collected = append(collected, lines[i])
			i++
		} else if isBlankLine(lines[i]) {
			// A blank line ends the blockquote
			break
		} else {
			// Lazy continuation: a non-blank, non-quote line continues the
			// blockquote only if the previous line was a quote line. This is
			// a simplification; CommonMark has more nuanced rules.
			break
		}
	}
	return Block{
		Type:    BlockBlockquote,
		Content: strings.Join(collected, "\n"),
	}, i
}

func parseList(lines []string, start int) (Block, int) {
	var collected []string
	i := start

	for i < len(lines) {
		if isListItem(lines[i]) {
			collected = append(collected, lines[i])
			i++
			continue
		}

		// Indented continuation line (part of the current list item)
		if !isBlankLine(lines[i]) && (strings.HasPrefix(lines[i], "  ") || strings.HasPrefix(lines[i], "\t")) {
			collected = append(collected, lines[i])
			i++
			continue
		}

		// Blank line: could be inside a loose list or the end of the list
		if isBlankLine(lines[i]) {
			// Peek ahead: if the next non-blank line is a list item or
			// indented continuation, the blank line is part of the list.
			j := i
			for j < len(lines) && isBlankLine(lines[j]) {
				j++
			}
			if j < len(lines) && (isListItem(lines[j]) || strings.HasPrefix(lines[j], "  ") || strings.HasPrefix(lines[j], "\t")) {
				// Include the blank line(s) as part of the list
				for i < j {
					collected = append(collected, lines[i])
					i++
				}
				continue
			}
			// End of list
			break
		}

		// Non-list, non-indented, non-blank → end of list
		break
	}

	return Block{
		Type:    BlockList,
		Content: strings.Join(collected, "\n"),
	}, i
}

func parseParagraph(lines []string, start int) (Block, int) {
	var collected []string
	i := start

	for i < len(lines) {
		if isBlankLine(lines[i]) {
			break
		}
		// If the line starts a different block type, stop the paragraph
		if i > start {
			if headingRegex.MatchString(lines[i]) {
				break
			}
			if fenceRegex.MatchString(lines[i]) {
				break
			}
			if isHorizontalRule(lines[i]) {
				break
			}
			if strings.HasPrefix(strings.TrimLeft(lines[i], " "), ">") {
				break
			}
			if isListItem(lines[i]) {
				break
			}
		}
		collected = append(collected, lines[i])
		i++
	}

	return Block{
		Type:    BlockParagraph,
		Content: strings.Join(collected, "\n"),
	}, i
}
