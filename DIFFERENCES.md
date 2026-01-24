# Syntactical Differences Between example.md and original_example.md

This document outlines the syntactical differences between the translated `example.md` and the source `original_example.md`.

## 1. Headings

| Original | Translated |
|----------|------------|
| `## Preliminaries` | `**Preliminaries**` (bold text) |
| `### Backslash Escapes` | `**Backslash Escapes**` (bold text) |
| `## H2 Heading` through `###### H6 Heading` | `**H2 Heading**` through `**H6 Heading**` (bold text) |
| ATX heading syntax (`#`, `##`, `###`, etc.) | Converted to bold (`**...**`) for H2-H6 |

Most headings below H1 are converted from ATX heading syntax to bold text.

## 2. Fenced Code Blocks

| Original | Translated |
|----------|------------|
| Triple backticks with content on multiple lines | Single backtick inline code spans |
| ``````` blocks | Inline code with single backticks |
| Language identifiers (e.g., `python`, `javascript`) | Removed |
| Multi-line code blocks | Split into multiple single-line code spans |

**Example:**
```
# Original
```python
def hello():
    print("world")
```

# Translated
`def hello():`
` print(&quotworld")`
```

## 3. Blockquotes

| Original | Translated |
|----------|------------|
| `> ` prefix syntax | Removed entirely |
| Nested blockquotes (`> > `) | Flattened or removed |
| Block-level quoting | Content appears as regular paragraphs |

## 4. Entity Encoding

| Original | Translated |
|----------|------------|
| `&nbsp;` | Rendered as space character |
| `&amp;` | `&` |
| `&copy;` | `©` |
| `&AElig;` | `Æ` |
| `&frac34;` | `¾` |
| `&#35;` | `\#` (escaped) |
| `&#1234;` | `Ӓ` |
| `<a title="...">` | `&lta title=&quota lot` |

HTML entities are resolved to their character equivalents, and some HTML tags are entity-encoded.

## 5. Backslash Escapes

| Original | Translated |
|----------|------------|
| `\!\"\#\$\%\&\'\(\)\*\+\,\-\.\/\:\;\<\=\>\?\@\[\\\]\^\_\`\{\|\}\~` | `!"#\$%&'()\*+,-./:;\<=\>?@\[\\^\_\`{\|}~` |

Many escape sequences are partially preserved or resolved inconsistently.

## 6. Links and Images

| Original | Translated |
|----------|------------|
| `[Link with **bold** and *italic*](/url "Title")` | `*Link with* ***bold*** *and* **italic**` |
| `![Image with **alt** text](url "Title")` | `![](data:image/png;base64,(null))` |
| `[Link with \`code\` in text](/url)` | `*Link with* *code* *in text*` |
| Reference-style link definitions | Removed entirely |

Links are converted to emphasized text, URLs are stripped, and images lose their source/alt text.

## 7. Lists

| Original | Translated |
|----------|------------|
| `+` and `*` list markers | Converted to `-` |
| `1.` ordered list markers | Preserved but with extra spacing (`1.  `) |
| Loose lists (with blank lines between items) | Collapsed to tight lists |

## 8. Tables

| Original | Translated |
|----------|------------|
| Full pipe-delimited table syntax | `[TABLE]` placeholder |

Tables are replaced with a `[TABLE]` placeholder text.

## 9. Thematic Breaks

| Original | Translated |
|----------|------------|
| `***`, `---`, `------` | Removed or converted to other elements |

Horizontal rules/thematic breaks are not preserved.

## 10. Setext Headings

| Original | Translated |
|----------|------------|
| `Foo *bar baz*` + `====` underline | Text with arrow character (`→`) |
| Setext H1/H2 underline syntax | Not preserved |

## 11. Footnotes

| Original | Translated |
|----------|------------|
| `[^1]` footnote reference | `*1*` (italicized) |
| `[^1]: The footnote content.` | `1.  The footnote content.*↩︎*` (as list item) |

Footnote syntax is converted to inline text with return arrow.

## 12. Inline Code

| Original | Translated |
|----------|------------|
| `` `code span with *emphasis*` `` | `This is a code span with \*emphasis\* inside.` |
| Double backticks ``` `` ``` | Preserved but spacing differs |

Code span contents may have their special characters escaped.

## 13. HTML Blocks

| Original | Translated |
|----------|------------|
| Raw HTML like `<a title="...">` | Entity-encoded: `&lta title=&quot...` |

HTML is entity-encoded rather than preserved raw.

## 14. Whitespace and Formatting

| Original | Translated |
|----------|------------|
| Hard line breaks (`\` at end of line) | Preserved as code blocks |
| Indented content | Indentation often removed |
| Tab characters | Converted to spaces or removed |

## 15. Special Emphasis Patterns

| Original | Translated |
|----------|------------|
| `___foo__ bar__` | `***foo*** *bar*\_` |
| Complex nested emphasis | Often parsed/rendered differently |

Underscore-based emphasis is converted to asterisk-based.

## Summary

The `example.md` file represents a lossy translation from `original_example.md` with these primary transformations:

1. **Structure loss**: Block-level elements (headings, blockquotes, code blocks) converted to inline equivalents
2. **Link/Image degradation**: URLs stripped, images replaced with base64 placeholders
3. **Entity resolution**: HTML entities converted to literal characters
4. **Syntax normalization**: Various list markers unified to `-`, emphasis unified to `*`
5. **Feature removal**: Tables, reference links, footnotes, and thematic breaks lost or replaced with placeholders
