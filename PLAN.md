# Bashly to Go Migration Plan

This document outlines the plan to migrate the `stash` CLI from a bashly-based bash implementation to a native Go implementation.

## Overview

### Current Architecture (Bash/Bashly)

```
src/
  bashly.yml              # CLI configuration
  push_command.sh         # Push command implementation
  pull_command.sh         # Pull command implementation
  diff_command.sh         # Diff command implementation
  completions_command.sh  # Completions command
  lib/
    # Pure functions (frontmatter manipulation)
    strip_frontmatter.sh
    get_id_from_frontmatter.sh
    update_frontmatter.sh
    extract_frontmatter.sh
    
    # File I/O
    read_markdown_file.sh
    write_markdown_file.sh
    
    # Format conversion (pandoc wrappers)
    markdown_to_html.sh
    html_to_markdown.sh
    
    # Apple Notes integration (AppleScript)
    create_note.sh
    read_note.sh
    update_note.sh
    delete_note.sh
    find_note.sh
    
    # Shell completions
    send_completions.sh
```

### Target Architecture (Go)

```
cmd/
  stash/
    main.go               # Entry point

internal/
  cli/
    root.go               # Root command setup
    push.go               # Push command
    pull.go               # Pull command
    diff.go               # Diff command
    completion.go         # Shell completions (built-in to Cobra)

  frontmatter/
    frontmatter.go        # All frontmatter operations
    frontmatter_test.go
  
  markdown/
    file.go               # File I/O operations
    file_test.go
  
  convert/
    convert.go            # Pandoc wrapper for HTML<->Markdown
    convert_test.go
  
  notes/
    notes.go              # Apple Notes AppleScript integration
    notes_test.go

go.mod
go.sum
```

---

## Migration Phases

### Phase 1: Project Scaffolding

**Goal**: Set up Go project structure with Cobra CLI framework.

**Tasks**:
1. Initialize Go module: `go mod init github.com/paradise-runner/cider`
2. Install Cobra: `go get github.com/spf13/cobra@latest`
3. Create directory structure:
   - `cmd/stash/main.go` - entry point
   - `internal/cli/root.go` - root command with version flag
   - `internal/` packages for business logic
4. Implement root command with:
   - `--help` flag (automatic with Cobra)
   - `--version` flag (set via ldflags at build time)
5. Verify: `go run ./cmd/stash --help` works

**Estimated effort**: Small

---

### Phase 2: Frontmatter Package

**Goal**: Implement pure frontmatter manipulation functions.

**Source functions to port**:
- `strip_frontmatter.sh` → `frontmatter.Strip(content string) string`
- `get_id_from_frontmatter.sh` → `frontmatter.GetAppleNotesID(content string) (string, error)`
- `update_frontmatter.sh` → `frontmatter.UpdateAppleNotesID(content, id string) string`
- `extract_frontmatter.sh` → `frontmatter.Extract(content string) string`

**Implementation notes**:
- Use regex or a YAML frontmatter library (e.g., `github.com/adrg/frontmatter`)
- Frontmatter format: `---\n<yaml>\n---\n`
- Key field: `apple_notes_id`
- The current bash uses `pcregrep` with multiline regex; Go's `regexp` package supports this

**Tests to port** (from `test/cases/unit_*_spec.sh`):
- `unit_strip_frontmatter_spec.sh` - 5 test cases
- `unit_get_id_from_frontmatter_spec.sh` - test cases for with/without ID
- `unit_update_frontmatter_spec.sh` - test cases for add/update ID
- `unit_extract_frontmatter_spec.sh` - test cases for extraction

**Test approach**: 
- Table-driven tests with `testing` package
- Use test fixtures from `test/fixtures/*.md`
- Compare against current approval outputs in `test/approvals/`

**Estimated effort**: Medium

---

### Phase 3: Markdown File I/O Package

**Goal**: Implement file read/write operations.

**Source functions to port**:
- `read_markdown_file.sh` → `markdown.ReadFile(path string) (string, error)`
- `write_markdown_file.sh` → `markdown.WriteFile(path, content string) error`

**Implementation notes**:
- Simple file I/O using `os.ReadFile` / `os.WriteFile`
- Error handling for file not found, permission errors
- The bash version requires file to exist for writes (no create)

**Tests to port**:
- `unit_read_markdown_file_spec.sh`
- `unit_write_markdown_file_spec.sh`

**Estimated effort**: Small

---

### Phase 4: Pandoc Conversion Package

**Goal**: Wrap pandoc for Markdown ↔ HTML conversion.

**Source functions to port**:
- `markdown_to_html.sh` → `convert.MarkdownToHTML(markdown string) (string, error)`
- `html_to_markdown.sh` → `convert.HTMLToMarkdown(html string) (string, error)`

**Implementation notes**:
- Shell out to `pandoc` using `os/exec`
- Markdown → HTML: `pandoc -f gfm -t html --wrap=none`
- HTML → Markdown: `pandoc -f html -t gfm-raw_html --wrap=none`
- Post-processing for `html_to_markdown`:
  1. Convert Apple Notes title pattern to H1: `<div><b><span style="font-size: 24px">...</span></b></div>` → `<h1>...</h1>`
  2. Remove `&nbsp;` list separators
  3. Trim trailing whitespace
  4. Collapse multiple blank lines

**Tests to port**:
- `unit_markdown_to_html_spec.sh`
- `unit_html_to_markdown_spec.sh`

**Dependencies**:
- Runtime dependency on `pandoc` binary
- Should check for pandoc availability and provide clear error message

**Estimated effort**: Medium

---

### Phase 5: Apple Notes Integration Package

**Goal**: Implement AppleScript-based Apple Notes operations.

**Source functions to port**:
- `create_note.sh` → `notes.Create(htmlContent string) (noteID string, error)`
- `read_note.sh` → `notes.Read(noteID string) (htmlContent string, error)`
- `update_note.sh` → `notes.Update(noteID, htmlContent string) error`
- `delete_note.sh` → `notes.Delete(noteID string) error`
- `find_note.sh` → `notes.Find(noteID string) (bool, error)`

**Implementation notes**:
- Execute AppleScript via `osascript` command using `os/exec`
- Note ID format: `x-coredata://...` (Core Data URL)
- Must handle "Recently Deleted" folder filtering
- Escape double quotes in content before passing to AppleScript
- Error handling for note not found, AppleScript errors

**AppleScript templates** (embed as string constants):
```applescript
tell application "Notes"
  try
    set newNote to make new note with properties {body:"{{content}}"}
    return id of newNote
  on error errMsg
    error errMsg
  end try
end tell
```

**Tests**:
- Integration tests require actual Apple Notes access
- Unit tests can mock the `exec.Command` interface
- Consider interface abstraction: `type NotesClient interface { Create, Read, Update, Delete, Find }`

**Estimated effort**: Medium-Large

---

### Phase 6: Push Command

**Goal**: Implement the `push` command.

**Current behavior** (from `push_command.sh`):
1. Read markdown file from path argument
2. Extract `apple_notes_id` from frontmatter (may not exist)
3. If ID exists, check if note exists in Apple Notes via `find_note`
4. If note not found:
   - Prompt user: "Create new note? (y/n)"
   - If yes: strip frontmatter, convert to HTML, create note, update frontmatter with new ID
   - If no: exit
5. If note found:
   - Strip frontmatter, convert to HTML, update note

**Cobra command setup**:
```go
var pushCmd = &cobra.Command{
    Use:   "push <file>",
    Short: "Push a Markdown file to Apple Notes (create or update)",
    Args:  cobra.ExactArgs(1),
    RunE:  runPush,
}
```

**Tests to port**:
- `e2e_push_command_spec.sh` - 6 test scenarios:
  - New file, user confirms → creates note
  - New file, user cancels → exits
  - Existing note → updates
  - Orphaned ID, user confirms → creates new
  - Orphaned ID, user cancels → exits
  - Non-existent file → error

**Estimated effort**: Medium

---

### Phase 7: Pull Command

**Goal**: Implement the `pull` command.

**Current behavior** (from `pull_command.sh`):
1. Read markdown file from path argument
2. Extract `apple_notes_id` from frontmatter (required)
3. Find note in Apple Notes (error if not found)
4. Read note content (HTML)
5. Convert HTML to Markdown
6. Extract existing frontmatter from local file
7. Combine: frontmatter + blank line + new markdown body
8. Write back to file

**Cobra command setup**:
```go
var pullCmd = &cobra.Command{
    Use:   "pull <file>",
    Short: "Pull content from Apple Notes to a Markdown file",
    Args:  cobra.ExactArgs(1),
    RunE:  runPull,
}
```

**Tests to port**:
- `e2e_pull_command_spec.sh` - test scenarios:
  - Valid ID → pulls and updates file
  - No ID in frontmatter → error
  - Note not found → error
  - Non-existent file → error

**Estimated effort**: Medium

---

### Phase 8: Diff Command

**Goal**: Implement the `diff` command.

**Current behavior** (from `diff_command.sh`):
1. Read markdown file
2. Extract and validate `apple_notes_id`
3. Find note in Apple Notes
4. Read note content, convert to markdown
5. Strip frontmatter from local content
6. Show unified diff: remote (Apple Notes) vs local (file)

**Implementation notes**:
- Use Go diff library (e.g., `github.com/sergi/go-diff`) or shell out to `diff`
- Output format: unified diff with labels "Apple Notes" and file path

**Cobra command setup**:
```go
var diffCmd = &cobra.Command{
    Use:   "diff <file>",
    Short: "Show differences between a Markdown file and its linked Apple Note",
    Args:  cobra.ExactArgs(1),
    RunE:  runDiff,
}
```

**Tests to port**:
- `e2e_diff_command_spec.sh`

**Estimated effort**: Small-Medium

---

### Phase 9: Shell Completions

**Goal**: Implement shell completion support.

**Implementation notes**:
- Cobra has built-in completion support
- Add completion command: `stash completion bash`, `stash completion zsh`, etc.
- File path completion for the `<file>` argument

```go
var completionCmd = &cobra.Command{
    Use:   "completion [bash|zsh|fish|powershell]",
    Short: "Generate shell completion script",
    // ... Cobra's built-in completion generation
}
```

**Estimated effort**: Small (mostly built-in to Cobra)

---

### Phase 10: Integration & Polish

**Goal**: Final integration, error handling, and polish.

**Tasks**:
1. Dependency checking at startup (pandoc availability)
2. Consistent error message formatting
3. Progress messages ("Reading file...", "Searching for note...", etc.)
4. Exit codes (0 = success, 1 = error, 2 = usage error for diff)
5. Build configuration:
   - Version injection via ldflags: `-ldflags "-X main.version=1.0.0"`
   - Build for macOS only (Apple Notes dependency)
6. Update Makefile with Go build targets
7. Full end-to-end testing

**Estimated effort**: Medium

---

## Implementation Order Summary

| Phase | Package/Command | Depends On | Effort |
|-------|-----------------|------------|--------|
| 1 | Project scaffolding | - | Small |
| 2 | `internal/frontmatter` | - | Medium |
| 3 | `internal/markdown` | - | Small |
| 4 | `internal/convert` | - | Medium |
| 5 | `internal/notes` | - | Medium-Large |
| 6 | `internal/cli/push.go` | 2, 3, 4, 5 | Medium |
| 7 | `internal/cli/pull.go` | 2, 3, 4, 5 | Medium |
| 8 | `internal/cli/diff.go` | 2, 3, 4, 5 | Small-Medium |
| 9 | `internal/cli/completion.go` | 1 | Small |
| 10 | Integration & polish | All | Medium |

**Total estimated effort**: ~2-3 days of focused work

---

## Key Differences from Bash Version

### Advantages of Go
- Single binary distribution (no bash version concerns)
- Better error handling with typed errors
- Easier testing with interfaces and mocking
- Cross-compilation support (though limited to macOS for Apple Notes)
- Better performance (though not critical for this use case)

### Maintained Parity
- Same CLI interface: `stash push|pull|diff <file>`
- Same user prompts and messages
- Same frontmatter format (`apple_notes_id`)
- Same pandoc dependency and conversion flags
- Same AppleScript integration approach

### Potential Improvements (Optional)
- `--yes` / `-y` flag to skip confirmation prompts
- `--quiet` / `-q` flag for scripting
- `--dry-run` flag to preview changes
- Better diff output with color support

---

## Testing Strategy

### Unit Tests
- Test pure functions (frontmatter, conversion helpers)
- Use table-driven tests
- Port existing fixtures from `test/fixtures/`
- Compare against approval baselines

### Integration Tests
- Test Apple Notes operations
- Require actual macOS with Notes.app
- Create → Test → Delete pattern (cleanup)
- Skip in CI with build tag: `//go:build integration`

### E2E Tests
- Test full CLI workflows
- Mock the `notes` package interface
- Verify file mutations and output messages

### Running Tests
```bash
# Unit tests only
go test ./...

# Include integration tests (requires Apple Notes)
go test -tags=integration ./...
```

---

## File Mapping Reference

| Bash Source | Go Target |
|-------------|-----------|
| `src/bashly.yml` | `internal/cli/root.go` (Cobra config) |
| `src/push_command.sh` | `internal/cli/push.go` |
| `src/pull_command.sh` | `internal/cli/pull.go` |
| `src/diff_command.sh` | `internal/cli/diff.go` |
| `src/completions_command.sh` | `internal/cli/completion.go` |
| `src/lib/strip_frontmatter.sh` | `internal/frontmatter/frontmatter.go` |
| `src/lib/get_id_from_frontmatter.sh` | `internal/frontmatter/frontmatter.go` |
| `src/lib/update_frontmatter.sh` | `internal/frontmatter/frontmatter.go` |
| `src/lib/extract_frontmatter.sh` | `internal/frontmatter/frontmatter.go` |
| `src/lib/read_markdown_file.sh` | `internal/markdown/file.go` |
| `src/lib/write_markdown_file.sh` | `internal/markdown/file.go` |
| `src/lib/markdown_to_html.sh` | `internal/convert/convert.go` |
| `src/lib/html_to_markdown.sh` | `internal/convert/convert.go` |
| `src/lib/create_note.sh` | `internal/notes/notes.go` |
| `src/lib/read_note.sh` | `internal/notes/notes.go` |
| `src/lib/update_note.sh` | `internal/notes/notes.go` |
| `src/lib/delete_note.sh` | `internal/notes/notes.go` |
| `src/lib/find_note.sh` | `internal/notes/notes.go` |

---

## Dependencies

### Runtime
- `pandoc` - Markdown ↔ HTML conversion (external binary)
- macOS with Apple Notes (for AppleScript integration)

### Build-time (Go modules)
- `github.com/spf13/cobra` - CLI framework
- `github.com/sergi/go-diff` - Diff generation (optional, could shell out to `diff`)
- `github.com/stretchr/testify` - Test assertions (optional but recommended)
