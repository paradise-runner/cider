# Go Implementation of Stash

This directory contains a complete Go reimplementation of the `stash` CLI tool, which was originally built with bash and bashly.

## Status

✅ **Complete** - All features from the bash version have been implemented:

- ✅ Push markdown files to Apple Notes (create or update)
- ✅ Pull content from Apple Notes to markdown files
- ✅ Show diffs between local files and Apple Notes
- ✅ Shell completions (bash, zsh, fish, powershell)
- ✅ Frontmatter management (YAML with `apple_notes_id`)
- ✅ Markdown ↔ HTML conversion via pandoc
- ✅ AppleScript integration with Apple Notes

## Building

```bash
# Build with default version
make build

# Build with specific version
make build VERSION=1.0.0

# Clean build artifacts
make clean
```

The binary will be created at `dist/stash-go`.

## Testing

```bash
# Run unit tests
make test-go

# Run integration tests (requires Apple Notes)
make test-go-integration
```

## Project Structure

```
cmd/stash/            # Entry point
  main.go

internal/
  cli/                # Cobra CLI commands
    root.go           # Root command setup
    push.go           # Push command
    pull.go           # Pull command
    diff.go           # Diff command
  
  frontmatter/        # YAML frontmatter operations
    frontmatter.go
    frontmatter_test.go
  
  markdown/           # File I/O
    file.go
    file_test.go
  
  convert/            # Pandoc wrapper
    convert.go
    convert_test.go
  
  notes/              # Apple Notes integration
    notes.go
    notes_test.go
```

## Dependencies

### Runtime
- **pandoc** - Markdown ↔ HTML conversion
- **macOS with Apple Notes** - For AppleScript integration

### Build-time (Go modules)
- `github.com/spf13/cobra` - CLI framework

## Usage

The Go version has the exact same CLI interface as the bash version:

```bash
# Push a file to Apple Notes
./dist/stash-go push my-note.md

# Pull changes from Apple Notes
./dist/stash-go pull my-note.md

# Show differences
./dist/stash-go diff my-note.md

# Generate shell completions
./dist/stash-go completion bash
```

## Comparison to Bash Version

| Feature | Bash | Go |
|---------|------|-----|
| Single binary | ❌ (requires bash 5+) | ✅ |
| Fast startup | ❌ | ✅ |
| Cross-compilation | ❌ | ✅ (macOS only due to Apple Notes) |
| Type safety | ❌ | ✅ |
| Testing | Approval tests | Standard Go tests + integration |
| Dependencies | bashly, Docker | Just Go toolchain |

## Migration Notes

The Go implementation maintains 100% feature parity with the bash version:

1. **Same CLI interface** - All commands, arguments, and flags are identical
2. **Same frontmatter format** - Uses `apple_notes_id` field
3. **Same pandoc flags** - Identical conversion behavior
4. **Same AppleScript** - Identical Apple Notes integration
5. **Same user prompts** - Confirmation dialogs work the same way

## Development

### Adding a new command

1. Create `internal/cli/yourcommand.go`
2. Implement the command with Cobra
3. Register in `init()` function
4. Add tests

### Running locally

```bash
go run ./cmd/stash <command> <args>
```

### Version injection

Version is set via ldflags at build time:

```bash
go build -ldflags "-X 'main.version=1.0.0'" -o stash ./cmd/stash
```
