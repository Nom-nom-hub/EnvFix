# Contributing to envfix

First, thank you for your interest in contributing to envfix! We appreciate all contributions, from bug reports to new features.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## How to Contribute

### Reporting Bugs

Before creating a bug report, check the existing issues to avoid duplicates.

When creating a bug report:
- **Use a clear, descriptive title**
- **Provide a step-by-step reproduction**
- **Provide specific examples**
- **Include environment information** (OS, Go version, envfix version)
- **Include error messages and logs**
- **Describe the actual vs expected behavior**

Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md).

### Suggesting Enhancements

Use the [feature request template](.github/ISSUE_TEMPLATE/feature_request.md).

When creating a feature request:
- **Describe the problem** the feature solves
- **Provide concrete use cases**
- **Describe the expected behavior**
- **Consider alternative solutions**

### Pull Requests

1. **Fork the repository** and create your feature branch
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Commit your changes** using Conventional Commits
   ```bash
   git commit -m "feat: add new feature"
   git commit -m "fix: resolve issue with X"
   git commit -m "docs: update README"
   git commit -m "test: add tests for X"
   git commit -m "chore: update dependencies"
   ```

3. **Ensure code quality**
   ```bash
   # Format code
   go fmt ./...
   
   # Run linter
   golangci-lint run ./...
   
   # Run tests
   go test -v ./...
   
   # Check coverage
   go test -cover ./...
   ```

4. **Push to your fork** and open a Pull Request
   ```bash
   git push origin feature/my-feature
   ```

5. **Use the PR template** to describe your changes

## Development Setup

### Prerequisites
- Go 1.21 or later
- Git
- (Optional) golangci-lint for linting

### Clone and Setup
```bash
git clone https://github.com/yourusername/envfix.git
cd envfix
go mod tidy
go mod download
```

### Building
```bash
# Build binary
go build -o envfix

# Build for other platforms
GOOS=windows GOARCH=amd64 go build -o envfix.exe
GOOS=linux GOARCH=amd64 go build -o envfix
GOOS=darwin GOARCH=amd64 go build -o envfix
```

### Testing
```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -cover ./...

# Run specific test
go test -v -run TestName ./...

# Run with race detector
go test -race ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run ./...

# Run vet
go vet ./...

# Check for issues
golangci-lint run --timeout=5m ./...
```

## Branching Model

- **main**: Stable, production-ready code
- **develop**: Pre-release integration branch
- **feature/**: Individual feature branches (branch off from develop)
- **fix/**: Bug fix branches (branch off from develop)

## Commit Conventions

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types
- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that don't affect code meaning (formatting, etc)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Code change that improves performance
- **test**: Adding or updating tests
- **chore**: Changes to build process, dependencies, etc

### Examples
```bash
git commit -m "feat: add orphaned environment detection"
git commit -m "fix: resolve version parsing issue"
git commit -m "docs: update README with new features"
git commit -m "test: add integration tests for repair"
git commit -m "chore: update dependencies"
```

## Code Style

We follow standard Go conventions:

1. **Format code with gofmt**
   ```bash
   go fmt ./...
   ```

2. **Follow Go naming conventions**
   - Package names: lowercase, single word
   - Exported names: PascalCase
   - Unexported names: camelCase
   - Constants: UPPER_SNAKE_CASE

3. **Write clear comments**
   - Export comments for all exported types/functions
   - Comments explain WHY, not WHAT
   - Keep comments up-to-date with code

4. **Error handling**
   - Always handle errors explicitly
   - Never use panic for control flow
   - Wrap errors with context

5. **Testing**
   - Write tests for all new functionality
   - Keep tests focused and clear
   - Use table-driven tests for multiple cases

## Documentation

### Code Documentation
- Document all exported functions and types
- Include examples where helpful
- Keep godoc comments up-to-date

### User Documentation
- Update README.md for user-facing changes
- Add to docs/ directory for detailed guides
- Update CHANGELOG.md with notable changes
- Add examples to examples/ directory

## Testing Requirements

### Unit Tests
- All new code must have unit tests
- Target 85%+ code coverage
- Test both happy path and error cases
- Use table-driven tests for multiple scenarios

### Integration Tests
- Test feature interactions
- Test cross-language scenarios
- Test on multiple platforms (Windows, Linux, macOS)

### Running Tests
```bash
# Run all tests
go test -v ./...

# Run with coverage report
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestName ./...

# Run with race detection
go test -race ./...
```

## Performance Considerations

- Avoid unnecessary allocations
- Use sync.Pool for frequent allocations
- Profile with pprof before optimizing
- Document performance-critical sections
- Include benchmarks for critical paths

## Security

- Never commit secrets or credentials
- Use environment variables for configuration
- Validate all user input
- Keep dependencies up-to-date
- Report security issues privately to maintainers

## Release Process

1. **Update version**
   - Update version in code
   - Update CHANGELOG.md
   - Update documentation

2. **Create PR**
   - Submit changes for review
   - Wait for approval

3. **Tag release**
   - `git tag v1.0.1`
   - `git push origin v1.0.1`

4. **GitHub Actions**
   - Automated build and release
   - Creates release artifacts
   - Publishes documentation

## Getting Help

- **Questions?** Open an issue with [question] label
- **Need guidance?** Check DEVELOPMENT.md
- **Stuck?** Comment on related PR/issue

## Review Process

Pull requests are reviewed by maintainers. We aim to respond within 48 hours.

### Review Criteria
- ✅ Code quality (formatting, style, design)
- ✅ Tests (coverage, pass rate, edge cases)
- ✅ Documentation (updated, clear, complete)
- ✅ Performance (no regressions)
- ✅ Security (no vulnerabilities)
- ✅ Compatibility (no breaking changes without justification)

## Merging

- PR must have at least one approval
- All CI checks must pass
- Commits should be squashed or organized logically
- Use "Squash and merge" for single-feature PRs

## Questions?

Feel free to open an issue or reach out to maintainers. We're here to help!

---

Thank you for contributing to envfix! 🎉
