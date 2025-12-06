# envfix Testing Guide

## Unit Tests

Run all unit tests:
```bash
go test ./...
```

Run specific package tests:
```bash
go test ./internal/detector/...
go test ./internal/utils/...
go test ./internal/parser/...
```

Run with coverage:
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Run verbose:
```bash
go test -v ./...
```

## Integration Tests

Integration tests are marked with `// +build integration` and test real system interactions.

Run integration tests:
```bash
go test -tags integration ./...
```

Run all tests including integration:
```bash
go test -tags integration -v ./...
```

## Manual Testing

### 1. Basic Scan

```bash
# Scan current environment
./envfix.exe scan

# With verbose output
./envfix.exe scan --verbose

# Export as JSON (future)
./envfix.exe scan --json
```

Expected output:
- Detected language versions
- Platform information
- Any warnings or issues
- Scan duration

### 2. Environment Manifest

```bash
# Generate env.yaml
./envfix.exe lock

# Generate and specify output
./envfix.exe lock --output my-env.yaml
```

Expected output:
- env.yaml created with version info
- All detected languages captured

### 3. Cache Cleanup

```bash
# Dry run (preview changes)
./envfix.exe clean --dry-run

# Actual cleanup
./envfix.exe clean
```

Expected output:
- List of caches to be removed
- Space freed in MB

### 4. Help and Explanations

```bash
# Get help on specific issue
./envfix.exe explain python_missing
./envfix.exe explain package_manager_conflict
./envfix.exe explain broken_symlink
```

### 5. Configuration

```bash
# Create local config
./envfix.exe init

# Create global config
./envfix.exe init --global

# View generated config
cat .envfix.yaml
```

## Test Scenarios

### Scenario 1: Healthy Environment (Current)

**Setup:** System with Python, Node.js, Rust, Go installed

**Run:**
```bash
./envfix.exe scan
```

**Expected:**
- ✓ All languages detected
- ✓ All versions displayed
- ✓ Status: healthy
- ✓ No issues reported

### Scenario 2: Missing Language

**Setup:** Remove or hide Python from PATH

**Run:**
```bash
./envfix.exe scan
```

**Expected:**
- ✗ Python shows as NOT FOUND
- ⚠ Issue reported
- ✓ Other languages detected

### Scenario 3: Clean Caches

**Setup:** Any system with cache directories

**Run:**
```bash
./envfix.exe clean
```

**Expected:**
- ✓ Caches identified
- ✓ Space freed calculated
- ✓ Report generated

### Scenario 4: Configuration Management

**Setup:** Fresh directory

**Run:**
```bash
./envfix.exe init
cat .envfix.yaml
./envfix.exe init --global
```

**Expected:**
- ✓ .envfix.yaml created locally
- ✓ Global config in ~/.envfix/config.yaml
- ✓ Properly formatted YAML

### Scenario 5: Dependency Detection

**Setup:** Directory with package.json and/or requirements.txt

**Run:**
```bash
./envfix.exe scan
```

**Expected:**
- ⚠ Missing lock file warnings
- ✓ Dependency counts reported

## Performance Testing

### Scan Performance

Target: < 5 seconds for typical environments

```bash
time ./envfix.exe scan
```

Current system (with all languages):
- ~2-3 seconds on Windows 11

### Memory Usage

Monitor memory during scan:

**Windows:**
```bash
tasklist /fi "IMAGENAME eq envfix.exe" /v
```

Expected: < 50 MB

## Regression Testing

Test compatibility after changes:

```bash
# Before changes
./envfix.exe scan > before.txt

# After changes
./envfix.exe scan > after.txt

# Compare
diff before.txt after.txt
```

## Continuous Testing

### Build Verification

```bash
# Build for Windows
go build -o envfix.exe

# Build for Linux
GOOS=linux go build -o envfix

# Build for macOS
GOOS=darwin go build -o envfix
```

### Cross-Platform Testing

1. **Windows:** Test full suite
2. **Linux:** Test in WSL or VM
3. **macOS:** Test on Mac hardware or CI

## Known Issues

### Potential Test Failures

1. **PATH modification side effects** - Some tests modify PATH; ensure cleanup
2. **Permission issues** - Cleanup tests may fail without proper permissions
3. **System-dependent** - Tests may behave differently on Windows vs Unix

### Workarounds

1. Run tests in isolated directory
2. Use `--dry-run` for safe testing
3. Don't run as admin unless necessary

## Test Coverage Goals

- **Utils**: > 90% coverage
- **Detector**: > 80% coverage
- **Parser**: > 85% coverage
- **Types**: > 70% coverage
- **Overall**: > 80% coverage

## CI/CD Testing

When setting up CI/CD:

```yaml
# GitHub Actions example
- name: Run Tests
  run: go test -v -race -coverprofile=coverage.out ./...

- name: Integration Tests
  run: go test -v -tags integration ./...

- name: Build
  run: |
    go build -o envfix
    GOOS=linux go build -o envfix-linux
    GOOS=darwin go build -o envfix-darwin
```

## Debugging Tests

Enable verbose logging:

```bash
# Test with verbose output
go test -v -run TestName ./...

# Test with race detector
go test -race ./...

# Test with coverage and logging
go test -v -coverprofile=coverage.out ./... 2>&1 | grep -E "RUN|PASS|FAIL"
```

## Adding New Tests

1. Create `*_test.go` file in package
2. Follow naming: `TestFeatureName`
3. Use `t.Error()` or `t.Fatalf()` for failures
4. Include setup and teardown
5. Document expected behavior

Example:
```go
func TestNewFeature(t *testing.T) {
	// Setup
	expected := "result"
	
	// Execute
	actual := NewFunction()
	
	// Verify
	if actual != expected {
		t.Errorf("NewFunction() = %q; want %q", actual, expected)
	}
}
```

## Test Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Cobra Testing](https://github.com/spf13/cobra/blob/main/cmd/cmd_test.go)
- [YAML Testing](https://pkg.go.dev/gopkg.in/yaml.v3)

## Troubleshooting

**Tests hang:**
- Check for deadlocks in concurrent code
- Verify file cleanup in tests
- Check for infinite loops

**Tests fail intermittently:**
- Could be timing issues
- Could be file system race conditions
- Add small delays with `time.Sleep()`

**Coverage gaps:**
- Run `go test -cover` to see summary
- Run `go tool cover -html=coverage.out` for visual view
- Add tests for uncovered branches
