# AGENTS.md - AI Coding Agent Instructions

This file provides context for AI coding agents working in the Cider codebase.

## Project Overview

**Cider** is a Go CLI tool for bidirectionally syncing Markdown files with Apple Notes.
- **Language:** Go 1.21+
- **CLI Framework:** Cobra (github.com/spf13/cobra)
- **Platform:** macOS only (uses AppleScript via `osascript`)
- **External Dependency:** `pandoc` for Markdown/HTML conversion

## Build/Test/Lint Commands

| Command | Description |
|---------|-------------|
| `make build` | Build binary to `dist/cider` |
| `make test-go` | Run unit tests |
| `make test-go-integration` | Run integration tests (requires Apple Notes) |
| `make clean` | Remove build artifacts |

### Running a Single Test

```bash
# Run a specific test function
go test -v ./internal/frontmatter -run TestGetAppleNotesID

# Run tests in a specific package
go test -v ./internal/convert/...

# Run with race detection
go test -race ./internal/...
```

### Development

```bash
# Run without building
go run ./cmd/cider <command> <args>

# Build with version
make build VERSION=1.0.0
```

## Directory Structure

```
cmd/cider/main.go        # Entry point (minimal - just calls cli.Execute())
internal/
  cli/                   # Cobra commands (push.go, pull.go, diff.go, root.go)
  convert/               # Pandoc wrapper (MarkdownToHTML, HTMLToMarkdown)
  frontmatter/           # YAML frontmatter operations
  markdown/              # File I/O utilities
  notes/                 # Apple Notes client (AppleScript integration)
test/fixtures/           # Test fixture files
```

## Code Style Guidelines

### Import Organization

Group imports in two blocks separated by a blank line:
1. Standard library (alphabetical)
2. Third-party and internal packages (alphabetical)

```go
import (
    "fmt"
    "os"
    "strings"

    "github.com/paradise-runner/cider/internal/convert"
    "github.com/spf13/cobra"
)
```

### Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Exported functions | PascalCase | `MarkdownToHTML`, `NewClient` |
| Unexported functions | camelCase | `runPush`, `pushSingleFile` |
| Command handlers | `run<Action>` | `runPush`, `runPull`, `runDiff` |
| Local variables | Short camelCase | `err`, `noteID`, `htmlContent` |
| Package-level regex | camelCase var | `frontmatterRegex` |
| Exported constants | PascalCase | `AppleNotesIDKey` |
| Files | lowercase, underscores | `notes_test.go` |
| Packages | Short, single word | `cli`, `notes`, `convert` |

### Error Handling

- Always wrap errors with context using `%w`:
```go
return "", fmt.Errorf("failed to read note: %w", err)
```

- Start error messages lowercase
- Return early on errors
- Check specific error types when needed:
```go
if os.IsNotExist(err) {
    return "", fmt.Errorf("file not found: %s", path)
}
```

### Function Patterns

**Command structure** (each command in its own file):
```go
var pushCmd = &cobra.Command{
    Use:   "push <file|directory>",
    Short: "Push Markdown file(s) to Apple Notes",
    Args:  cobra.ExactArgs(1),
    RunE:  runPush,
}

func init() {
    rootCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) error {
    // Implementation
}
```

**Client pattern:**
```go
type Client struct{}

func NewClient() *Client {
    return &Client{}
}

func (c *Client) Create(content string) (string, error) {
    // Implementation
}
```

### Testing Patterns

Use table-driven tests:
```go
tests := []struct {
    name      string
    input     string
    want      string
    wantError bool
}{
    {
        name:  "valid input",
        input: "test",
        want:  "expected",
    },
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := Function(tt.input)
        if (err != nil) != tt.wantError {
            t.Errorf("unexpected error: %v", err)
        }
        if got != tt.want {
            t.Errorf("got %q, want %q", got, tt.want)
        }
    })
}
```

**Test helpers:**
```go
func loadFixture(t *testing.T, name string) string {
    t.Helper()
    path := filepath.Join("../../test/fixtures", name)
    content, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("failed to load fixture %s: %v", name, err)
    }
    return string(content)
}
```

### Comments

- Document all exported functions with a single-line comment
- Use inline comments only for non-obvious logic
- No comments for simple/obvious code

```go
// ReadFile reads a markdown file and returns its content
func ReadFile(path string) (string, error) {
```

### Package Structure

- One primary responsibility per package
- Each package has a main file (`notes.go`, `convert.go`)
- Tests co-located: `*_test.go`
- All library code in `internal/` (Go idiom for private packages)
- `cmd/cider/main.go` should remain minimal

## Key Domain Concepts

- **Frontmatter:** YAML block at start of Markdown files, delimited by `---`
- **apple_notes_id:** Frontmatter field linking local file to Apple Note
- Notes are stored as HTML in Apple Notes; conversion uses pandoc

## Integration Test Notes

Integration tests (`//go:build integration`) require:
- macOS with Apple Notes app
- Tests create/delete real notes - use with caution
- Run with: `make test-go-integration`
