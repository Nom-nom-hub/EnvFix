# Version Management Guide

## Overview

envfix v1.0.1 includes comprehensive version detection and conflict management to ensure your development environment uses compatible and supported language versions.

## Features

### Version Conflict Detection

#### Multiple Version Managers
Detects when multiple version managers are installed:

```bash
# Python version managers
envfix scan

# Reports:
# Warning: Multiple Python version managers detected: pyenv, conda, system
```

Manages:
- Python: pyenv, conda, system Python
- Node.js: nvm, n, nodeenv
- Rust: rustup (default manager)
- Go: gimme or other managers

#### Minimum Version Requirements

Validates that all interpreters meet minimum requirements:

| Language | Minimum | Reason |
|----------|---------|--------|
| Python | 3.7 | Security + features |
| Node.js | 14 | LTS support |
| Rust | 1.56 | Editions stability |
| Go | 1.16 | Module support |

#### Outdated & EOL Detection

Detects end-of-life or deprecated versions:

```bash
envfix scan

# Example outputs:
# Critical: Python 2 is end-of-life (EOL: Jan 1, 2020)
# Critical: Node.js 12 is end-of-life (EOL: Apr 30, 2022)
# Warning: Python 3.6 is below minimum (3.7+)
# Warning: Rust 1.50.0 is outdated (recommend 1.56+)
```

### Version Compatibility Validation

```bash
# Validate all installed versions are compatible
envfix repair

# Shows:
# [x] Python: /usr/bin/python3
# [x] Node.js: /usr/bin/node
# [x] npm: /usr/bin/npm
# [x] Rust: /usr/bin/cargo
# [x] Go: /usr/bin/go
# ✓ All installed versions are compatible
```

---

## Common Issues & Solutions

### Python 2 Detected

**Problem**:
```
Critical: Python 2 is end-of-life
Python 2 reached end-of-life on January 1, 2020. Upgrade to Python 3 immediately.
```

**Solution**:
```bash
# Install Python 3
sudo apt-get install python3  # Linux
brew install python3           # macOS
# Or download from https://www.python.org/

# Or use a version manager
pyenv install 3.11.0
pyenv global 3.11.0
```

### Multiple Python Version Managers

**Problem**:
```
Warning: Multiple Python version managers detected: pyenv, conda, system
```

**Solutions**:
1. **Pick one**: Choose pyenv, conda, or system Python
2. **Configure PATH**: Ensure only one is in PATH
3. **Set global**: Use `pyenv global 3.11.0` or `conda env list`

**Example with pyenv**:
```bash
# Remove conda from PATH
# Edit ~/.bashrc and remove conda initialization

# Or use pyenv exclusively
pyenv install 3.11.0
pyenv global 3.11.0
```

### Node.js 12 Detected

**Problem**:
```
Critical: Node.js 12 is end-of-life
Node.js 12 reached end-of-life on April 30, 2022. Upgrade to 16+ immediately.
```

**Solution**:
```bash
# Using nvm
nvm install 18
nvm alias default 18
nvm use 18

# Or from nodejs.org
# Download from https://nodejs.org/ (LTS version)
```

### Version Below Minimum

**Problem**:
```
Warning: Python version below minimum requirement
Python 3.7+ is required. Current: 3.6.0
```

**Solution**:
```bash
# Upgrade Python
pyenv install 3.9.0
pyenv global 3.9.0

# Or system package manager
sudo apt-get install python3.9  # Linux
brew install python@3.9         # macOS
```

### Multiple Rust Toolchains

**Problem**:
```
Warning: Multiple Rust toolchains detected: 5 versions
```

**Solution**:
```bash
# List installed toolchains
rustup toolchain list

# Remove old ones
rustup toolchain remove 1.65-gnu
rustup toolchain remove 1.66-gnu

# Keep only latest stable
rustup update
```

### Multiple Go Versions

**Problem**:
```
Warning: Multiple Go versions detected in ~/.gimme/versions
```

**Solution**:
```bash
# Remove old versions manually
rm -rf ~/.gimme/versions/go/go1.18.0
rm -rf ~/.gimme/versions/go/go1.19.0

# Or reinstall Go
go install golang.org/dl/go1.21.0@latest
~/go/bin/go1.21.0 download
```

---

## Configuration

### Local Configuration (.envfix.yaml)

```yaml
versions:
  # Minimum required versions
  python_min: "3.8"      # Default: 3.7
  node_min: "16"         # Default: 14
  rust_min: "1.56"       # Default: 1.56
  go_min: "1.18"         # Default: 1.16

  # Check for EOL versions
  check_eol: true        # Default: true

  # Check for outdated versions
  check_outdated: true   # Default: true

  # Report on multiple version managers
  check_conflicts: true  # Default: true
```

### Global Configuration (~/.envfix/config.yaml)

```yaml
versions:
  # Strict mode: fail if versions below minimum
  strict_mode: false

  # Auto-suggest upgrades
  suggest_upgrades: true

  # Check compatibility between versions
  validate_compatibility: true
```

---

## Version Compatibility Matrix

Recommended version combinations for stability:

| Python | Node.js | Rust | Go | Status |
|--------|---------|------|-----|--------|
| 3.11+ | 18+ | 1.70+ | 1.21+ | ✅ Latest |
| 3.10 | 16+ | 1.65+ | 1.20 | ✅ LTS |
| 3.9 | 16+ | 1.60+ | 1.19 | ⚠️ Aging |
| 3.8 | 14+ | 1.56+ | 1.18 | ⚠️ Old |
| 3.7 | 12 | 1.50 | 1.16 | ❌ EOL Risk |
| <3.7 | <12 | <1.50 | <1.16 | ❌ EOL |

---

## Update Recommendations

### Python
```bash
# Check installed version
python3 --version

# Upgrade
pyenv install 3.11.0 && pyenv global 3.11.0

# Or with system package manager
sudo apt-get install python3.11  # Linux
brew install python@3.11          # macOS
```

### Node.js
```bash
# Check installed version
node --version

# Upgrade with nvm
nvm install 18
nvm alias default 18

# Or from nodejs.org
# Download LTS version from https://nodejs.org/
```

### Rust
```bash
# Check installed version
rustc --version

# Upgrade
rustup update
rustup update stable
```

### Go
```bash
# Check installed version
go version

# Upgrade
go install golang.org/dl/go1.21.0@latest
~/go/bin/go1.21.0 download
```

---

## Performance Impact

| Operation | Time | Memory |
|-----------|------|--------|
| Check versions | <1s | <5MB |
| Parse versions | <100ms | <1MB |
| Validate compatibility | <500ms | <5MB |
| Suggest upgrades | 1-2s | <10MB |

---

## Advanced Usage

### Get Version Information

```bash
# Get detailed version info
envfix explain python_version_low
envfix explain node_eol
envfix explain rust_outdated
```

### Check Compatibility

```bash
# Dry-run to see issues without fixing
envfix repair --dry-run

# Shows all version-related issues
```

### Custom Minimum Versions

```yaml
# .envfix.yaml
versions:
  python_min: "3.10"   # Require Python 3.10+
  node_min: "18"       # Require Node 18+
  strict_mode: true    # Fail if not met
```

---

## Integration with CI/CD

### GitHub Actions

```yaml
name: Validate Environment

on: [push, pull_request]

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/setup-python@v4
        with:
          python-version: '3.11'
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Validate versions
        run: |
          envfix repair --dry-run
          
          # Or fail if versions incorrect
          envfix scan --json | jq '.issues[] | select(.id == "python_eol")'
```

### GitLab CI

```yaml
validate_versions:
  stage: test
  image: ubuntu:22.04
  script:
    - apt-get update && apt-get install -y golang-go
    - go build && ./envfix repair --dry-run
    - ./envfix scan --json
```

---

## FAQ

**Q: Should I use pyenv, conda, or system Python?**  
A: Pick one. Recommend pyenv for flexibility, conda for data science, system Python for simplicity.

**Q: Can I use Python 3.6?**  
A: Not recommended. Python 3.6 EOL was Dec 2021. Upgrade to 3.8+.

**Q: Is Node.js 12 still supported?**  
A: No. Node.js 12 EOL was Apr 2022. Upgrade to 16+ (LTS) or 18+ (current).

**Q: How do I remove old Go versions?**  
A: Run `rm -rf ~/.gimme/versions/go/go<old-version>` for each old version.

**Q: What if my project requires Python 3.6?**  
A: This is a breaking change. Consider:
1. Migrating to Python 3.8+
2. Using Docker for isolation
3. Maintaining separate environments

**Q: Can envfix auto-upgrade versions?**  
A: No, it detects and recommends. Manual upgrade is safer. v1.1 may add assisted upgrades.

---

## Troubleshooting

### Version detection not working
1. Check interpreter is in PATH: `which python3`
2. Verify version output: `python3 --version`
3. Check for corruption: `python3 -m venv /tmp/test`

### Multiple version managers detected
1. Check PATH: `echo $PATH`
2. Remove one from PATH
3. Verify with: `which python`

### Compatibility validation fails
1. Check all interpreters installed: `which python3 node python cargo go`
2. Verify PATH not corrupted
3. Try reinstalling problematic tool

---

## Related Commands

```bash
envfix scan      # Detect version issues
envfix repair    # Validate versions
envfix clean     # Clean caches
envfix doctor    # Full diagnostic
```

## See Also

- [README.md](./README.md) - Main documentation
- [DEPENDENCIES.md](./DEPENDENCIES.md) - Dependency management
- [ORPHANED_ENVIRONMENTS.md](./ORPHANED_ENVIRONMENTS.md) - Orphaned cleanup
- [QUICK_START.md](./QUICK_START.md) - Getting started
