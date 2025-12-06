# Dependency Management Guide

## Overview

envfix v1.0.1 includes comprehensive dependency management features that help keep your project dependencies clean, up-to-date, and reproducible across environments.

## Features

### Python Dependencies

#### Install & Upgrade
```bash
# Upgrade Python dependencies from requirements.txt
envfix repair

# The repair action will:
# - Find pip in venv or system PATH
# - Run: pip install --upgrade -r requirements.txt
# - Fall back to: pip install -r requirements.txt if upgrade fails
```

#### Verification
```bash
# Verify Python dependencies are properly installed
envfix repair

# Runs: pip check
# Detects conflicting or missing packages
```

#### Lock File Management
```bash
# Regenerate requirements.txt from current venv
envfix repair

# Captures current installed versions for reproducibility
```

### Node.js Dependencies

#### Clean Install (npm ci)
```bash
# Install dependencies using npm ci
envfix repair

# Automatically uses:
# - npm ci (if package-lock.json exists)
# - npm install (if no lock file)
# - pnpm install (for pnpm projects)
# - yarn install (for yarn projects)
```

#### Why npm ci?
- **Reproducible**: Uses exact versions from package-lock.json
- **Faster**: Skips dependency resolution
- **Safer**: Fails if lock file is out of sync
- **Recommended**: Best practice for CI/CD pipelines

#### Verification
```bash
# Verify Node.js dependencies
envfix repair

# Runs: npm list --depth=0
# Shows installed packages with any issues
```

#### Lock File Consistency
```bash
# Check lock file validity
envfix repair

# Detects:
# - Missing lock files with package.json
# - Multiple lock files (npm, pnpm, yarn)
# - Conflicting package managers
```

### Rust Dependencies

#### Update Dependencies
```bash
# Update Rust dependencies via Cargo
envfix repair

# Runs: cargo update
# Updates Cargo.lock to latest compatible versions
```

#### Verification
```bash
# Verify Rust dependencies
envfix repair

# Runs: cargo check
# Compiles code without producing binary
# Detects any dependency issues
```

#### Cache Management
```bash
# Clear cargo registry cache
envfix clean

# Frees disk space from downloaded crates
```

### Go Dependencies

#### Tidy Modules
```bash
# Clean up Go module dependencies
envfix repair

# Runs: go mod tidy
# Removes unused dependencies
# Adds missing dependencies
```

#### Verification
```bash
# Verify Go module integrity
envfix repair

# Runs: go mod verify
# Checks that modules are valid
```

#### Cache Management
```bash
# Clear Go module cache
envfix clean

# Frees disk space from module cache
```

## Dependency Verification

### What Gets Checked

envfix verifies:
1. **Installation**: All dependencies are installed
2. **Consistency**: Lock files match package configurations
3. **Completeness**: All required dependencies exist
4. **Validity**: No conflicting versions

### Check Dependency Lock Files

```bash
envfix repair

# Validates:
# - Python: requirements.txt exists for projects
# - Node.js: No multiple lock files (npm/yarn/pnpm)
# - Rust: Cargo.lock exists for binaries
# - Go: go.sum matches go.mod
```

## Lock File Management

### Regenerating Lock Files

Lock files ensure reproducible installs across environments.

```bash
# Regenerate all lock files
envfix repair

# Python:
#   pip freeze > requirements.txt

# Node.js:
#   npm install  # Regenerates package-lock.json

# Rust:
#   cargo update  # Regenerates Cargo.lock

# Go:
#   go mod tidy   # Updates go.sum
```

### Best Practices

1. **Always commit lock files** to version control
2. **Use lock files in CI/CD** for reproducible builds
3. **Regenerate when updating** dependencies
4. **Keep lock files in sync** with package manifests

## Cache Management

### Dependency Caches

envfix can clean up caches to free disk space:

```bash
# Clean all dependency caches
envfix clean

# Clears:
# - Python pip cache
# - Node.js npm cache
# - Rust cargo registry cache
# - Go module cache
```

### Cache Locations

| Language | Cache Location | Size |
|----------|---|---|
| Python | ~/.cache/pip | Often 500MB+ |
| Python (Windows) | %APPDATA%\pip\Cache | Often 500MB+ |
| Node.js | ~/.npm | Often 1GB+ |
| Rust | ~/.cargo/registry/cache | Often 500MB+ |
| Go | $GOPATH/pkg/mod/cache | Often 100MB+ |

## Dependency Issues

### Common Problems & Solutions

#### Python: Missing requirements.txt
**Problem**: Project has Python code but no requirements.txt

**Solution**:
```bash
# Create requirements.txt from current venv
pip freeze > requirements.txt

# Or use envfix to detect and fix
envfix repair
```

#### Node.js: Multiple lock files
**Problem**: package-lock.json, yarn.lock, pnpm-lock.yaml all present

**Solution**:
```bash
# envfix will detect this
envfix scan

# Remove conflicting lock files
rm yarn.lock pnpm-lock.yaml
npm install  # Regenerate package-lock.json

# Or use repair
envfix repair  # Automatically resolves
```

#### Rust: Cargo.lock conflicts
**Problem**: Binary project missing Cargo.lock

**Solution**:
```bash
# Add Cargo.lock to version control for binaries
cargo update  # Generates Cargo.lock

# Or use envfix
envfix repair
```

#### Go: Out of sync go.sum
**Problem**: go.mod and go.sum don't match

**Solution**:
```bash
# Sync them
go mod tidy

# Or use envfix
envfix repair
```

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup dependencies
        run: |
          # Use envfix to verify and install
          envfix repair
          
      - name: Run tests
        run: npm test
```

### GitLab CI Example

```yaml
stages:
  - setup
  - test

setup:
  stage: setup
  script:
    - envfix repair  # Install and verify dependencies
    
test:
  stage: test
  script:
    - npm test
```

## Configuration

### Local Configuration (.envfix.yaml)

```yaml
dependencies:
  auto_upgrade: true
  verify_lockfiles: true
  clean_caches: true
  
repair:
  install_python: true
  install_node: true
  update_rust: true
  tidy_go: true
```

### Global Configuration (~/.envfix/config.yaml)

```yaml
dependencies:
  # Use npm ci for deterministic installs
  npm_ci: true
  
  # Auto-clean caches
  auto_clean_caches: true
  
  # Regenerate lock files
  regenerate_lockfiles: false  # Manual only
```

## Advanced Usage

### Verify Only (No Changes)

```bash
# Check dependency status without modifying
envfix repair --dry-run

# Shows what would be done:
# [1/4] Verify all dependencies are properly installed
# [2/4] Check dependency lock file consistency
# [3/4] Clean dependency caches
```

### Force Regenerate Lock Files

```bash
# Manually regenerate all lock files
envfix clean
# Then for Python:
pip freeze > requirements.txt
# For Node.js:
npm install
```

### Check Specific Language

```bash
# Check only Node.js dependencies
envfix explain nodejs_missing

# Check only Python dependencies
envfix explain python_missing
```

## Performance

| Operation | Time | Impact |
|-----------|------|--------|
| Verify dependencies | 1-3s | None (read-only) |
| Check lock files | <1s | None (read-only) |
| Install Python deps | 5-60s | Installs packages |
| Install Node.js deps | 5-120s | Installs packages |
| Update Rust deps | 5-30s | Updates Cargo.lock |
| Tidy Go modules | 1-5s | Updates go.sum |
| Clean caches | 1-10s | Frees disk space |

## Troubleshooting

### "pip not found"
```bash
# Install pip
python -m ensurepip

# Or reinstall Python with pip
```

### "npm ci fails with error"
```bash
# Delete node_modules and retry
rm -rf node_modules package-lock.json
envfix repair
```

### "cargo update takes too long"
```bash
# Clear cargo cache first
envfix clean

# Then update
envfix repair
```

### "go mod tidy causes merge conflicts"
```bash
# Resolve conflicts in go.mod, then:
go mod tidy
git add go.mod go.sum
```

## FAQ

**Q: Should I commit lock files?**  
A: Yes! Lock files ensure reproducible builds across machines.

**Q: What's the difference between install and update?**  
A: Install uses exact versions from lock files. Update gets compatible newer versions.

**Q: Can I use different package managers?**  
A: For a single project, use one. envfix will warn about conflicts.

**Q: Does envfix handle monorepos?**  
A: Basic support. Full monorepo support coming in v1.1.

**Q: How often should I update dependencies?**  
A: Regularly (weekly/monthly) but use lock files to control versions.

## Related Commands

```bash
envfix scan      # Detect dependency issues
envfix repair    # Install and verify dependencies
envfix clean     # Clean dependency caches
envfix doctor    # Full diagnostic including dependencies
```

## See Also

- [README.md](./README.md) - Main documentation
- [ORPHANED_ENVIRONMENTS.md](./ORPHANED_ENVIRONMENTS.md) - Orphaned cleanup
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Technical architecture
