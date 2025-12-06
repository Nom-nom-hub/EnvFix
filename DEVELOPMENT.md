# envfix Development Guide

## Setup Development Environment

### Prerequisites
- Go 1.21 or later
- Git
- Make (optional, for using Makefile)

### Installation

```bash
# Clone repository
git clone https://github.com/aicode-dev/envfix
cd envfix

# Install dependencies
go mod tidy

# Build
go build -o envfix

# Run
./envfix --help
```

## Project Structure

```
envfix/
├── main.go                 # Entry point
├── go.mod / go.sum         # Dependencies
├── Makefile                # Build automation
├── README.md               # User documentation
├── ARCHITECTURE.md         # Technical design
├── TESTING.md              # Test guide
├── DEVELOPMENT.md          # This file
├── CHANGELOG.md            # Version history
├── cmd/                    # CLI commands
│   ├── root.go            # Command setup
│   ├── scan.go            # Scan implementation
│   ├── repair.go          # Repair implementation
│   ├── clean.go           # Clean implementation
│   ├── lock.go            # Lock/manifest
│   ├── doctor.go          # Doctor (scan+repair)
│   ├── explain.go         # Explanation command
│   └── init.go            # Init config
├── internal/
│   ├── config/            # Configuration system
│   │   └── config.go
│   ├── types/             # Data structures
│   │   └── types.go
│   ├── utils/             # Utilities
│   │   ├── system.go      # System utilities
│   │   ├── logger.go      # Logging
│   │   └── *_test.go      # Tests
│   ├── parser/            # Dependency parsing
│   │   ├── dependencies.go
│   │   └── *_test.go
│   ├── detector/          # Environment detection
│   │   ├── detector.go
│   │   ├── python.go
│   │   ├── node.go
│   │   ├── rust.go
│   │   ├── go_env.go
│   │   ├── dependencies.go
│   │   └── *_test.go
│   ├── rebuilder/         # Repair operations
│   │   ├── rebuilder.go
│   │   ├── python_repair.go
│   │   └── node_repair.go
├── spec/                  # Manifest generation
│   └── manifest.go
└── cleaner/               # Cleanup operations
    └── cleaner.go
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature
```

### 2. Make Changes

Follow the code style and structure of the project.

### 3. Test Your Changes

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...

# Integration tests
go test -tags integration ./...

# Specific package
go test -v ./internal/detector/
```

### 4. Build and Verify

```bash
# Build
go build -o envfix

# Test manually
./envfix scan
./envfix lock
./envfix clean
```

### 5. Commit and Push

```bash
git add .
git commit -m "feat: add your feature"
git push origin feature/your-feature
```

## Code Style

### Naming Conventions

- **Packages**: lowercase, single word (e.g., `detector`, `rebuilder`)
- **Functions**: CamelCase, exported start with capital (e.g., `Scan()`, `DetectPython()`)
- **Constants**: UPPER_SNAKE_CASE
- **Variables**: camelCase

### File Organization

1. Package declaration
2. Imports (grouped: stdlib, external, local)
3. Type definitions
4. Interface definitions
5. Function implementations
6. Helper functions

### Comments

- Exported functions have doc comments: `// FunctionName does something`
- Complex logic has inline comments
- Avoid redundant comments

### Error Handling

```go
if err != nil {
    utils.Error("descriptive message: %v", err)
    return err
}
```

## Adding New Commands

### 1. Create command file in `cmd/`

```go
package cmd

import "github.com/spf13/cobra"

var newCmd = &cobra.Command{
    Use:   "newcommand",
    Short: "Short description",
    Long:  "Longer description",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

### 2. Add flags if needed

```go
func init() {
    newCmd.Flags().BoolVar(&flag, "flag", false, "Description")
    rootCmd.AddCommand(newCmd)
}
```

### 3. Register in root.go

Already handled if you add to init()

## Adding Language Support

### 1. Create detector in `internal/detector/newlang.go`

```go
func (d *Detector) DetectNewLang() types.NewLangEnv {
    env := types.NewLangEnv{
        Status: "missing",
        Issues: []string{},
    }
    
    // Detection logic
    
    return env
}
```

### 2. Add type to `internal/types/types.go`

```go
type NewLangEnv struct {
    Version    string
    Executable string
    Status     string
    Issues     []string
}
```

### 3. Update scan to detect new language

In `detector.go`:
```go
result.NewLang = d.DetectNewLang()
```

### 4. Add repair if applicable

Create `internal/rebuilder/newlang_repair.go`:
```go
func (r *Rebuilder) RepairNewLang() error {
    // Implementation
    return nil
}
```

## Testing Strategy

### Unit Tests

Located in same package:
```go
func TestFunction(t *testing.T) {
    expected := "result"
    actual := Function()
    if actual != expected {
        t.Errorf("Function() = %q; want %q", actual, expected)
    }
}
```

### Integration Tests

With build tag:
```go
// +build integration

func TestFullScan(t *testing.T) {
    // Tests actual system interactions
}
```

### Coverage Goals

- Utils: > 90%
- Detector: > 80%
- Parser: > 85%
- Overall: > 80%

## Performance Considerations

### Scan Optimization

- Detect languages in parallel (future)
- Cache executable paths
- Minimize file I/O
- Use fast path for PATH lookup

### Memory

- Stream large files when possible
- Reuse buffers
- Avoid storing entire caches in memory

## Security Checklist

- [ ] No arbitrary code execution
- [ ] Validate all external input
- [ ] No hardcoded credentials
- [ ] Proper error messages (no paths in production)
- [ ] Safe file operations with backups
- [ ] Audit logging of changes

## Debugging

### Enable Verbose Output

```bash
./envfix scan --verbose --log-file debug.log
```

### Check Logs

```bash
tail -f debug.log
```

### Debug a Specific Function

Add logging:
```go
utils.Debug("Variable value: %v", variable)
```

### Run with Race Detector

```bash
go test -race ./...
```

## Building for Distribution

### Windows

```bash
go build -o envfix.exe
# For signed binary
signtool sign /f cert.pfx envfix.exe
```

### macOS

```bash
GOOS=darwin GOARCH=amd64 go build -o envfix
# For Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o envfix-m1
```

### Linux

```bash
GOOS=linux GOARCH=amd64 go build -o envfix
```

## Release Checklist

- [ ] All tests passing
- [ ] Coverage > 80%
- [ ] README updated
- [ ] CHANGELOG updated
- [ ] Version bumped (semantic versioning)
- [ ] Cross-platform builds successful
- [ ] Manual testing on all platforms
- [ ] Documentation updated
- [ ] Git tags created

## Common Tasks

### Add a new environment check

1. Add check to `Detector.Scan()`
2. Populate `ScanResult` with findings
3. Add issue if critical
4. Add warning if informational

### Add a new repair operation

1. Implement in `Rebuilder`
2. Add to `planRepairs()`
3. Make reversible with backup
4. Support dry-run mode

### Add configuration option

1. Add field to `Config` struct
2. Update YAML marshaling
3. Document in README
4. Add tests

## Troubleshooting Development

### Build Fails

```bash
go clean
go mod tidy
go mod verify
go build
```

### Tests Timeout

Increase timeout:
```bash
go test -timeout 2m ./...
```

### Race Condition Detected

```bash
go test -race ./...
```

Review concurrent code sections.

### Import Errors

```bash
go mod tidy
go mod download
```

## Contributing

See CONTRIBUTING.md for contribution guidelines.

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Documentation](https://cobra.dev/)
- [YAML v3 Documentation](https://pkg.go.dev/gopkg.in/yaml.v3)
- [Effective Go](https://golang.org/doc/effective_go)

## Getting Help

- Check existing issues on GitHub
- Review ARCHITECTURE.md for design decisions
- Check test files for usage examples
- Review similar functions in codebase
