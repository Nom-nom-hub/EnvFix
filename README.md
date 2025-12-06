# envfix

[![CI](https://github.com/Nom-nom-hub/EnvFix/actions/workflows/ci.yml/badge.svg)](https://github.com/Nom-nom-hub/EnvFix/actions/workflows/ci.yml)
[![Release](https://github.com/Nom-nom-hub/EnvFix/actions/workflows/release.yml/badge.svg)](https://github.com/Nom-nom-hub/EnvFix/releases)
[![GitHub Release](https://img.shields.io/github/v/release/Nom-nom-hub/EnvFix?label=release)](https://github.com/Nom-nom-hub/EnvFix/releases)
[![License](https://img.shields.io/github/license/Nom-nom-hub/EnvFix)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/doc/devel/release)

A cross-language environment and dependency doctor that diagnoses, repairs, standardizes, and cleans developer environments across Python, Node.js, Rust, Go, and system-level toolchains.

## Features

### Core Commands

- **`envfix scan`** - Deep diagnostic scan of your environment
  - Detects missing interpreters, broken venvs, PATH issues
  - Identifies package manager conflicts
  - Scans for orphaned toolchains and corrupted caches
  - Completes in under 5 seconds

- **`envfix repair`** - Automatically fix detected issues
  - Rebuilds virtual environments (Python)
  - Resolves package manager conflicts (Node.js)
  - Reinstalls corrupted dependencies
  - All operations are reversible and logged

- **`envfix clean`** - Remove orphaned and unused artifacts
  - Clears pip, npm, pnpm, yarn caches
  - Removes unused venv directories
  - Deletes broken symlinks
  - Cleans stale lockfiles

- **`envfix lock`** - Generate reproducible environment manifest
  - Creates `env.yaml` with all environment details
  - Captures language versions, toolchain versions, dependencies
  - Includes checksums for reproducibility
  - Enables cross-machine environment sync (Pro feature)

- **`envfix doctor`** - One-shot scan + repair
  - Convenience command for full environment diagnosis and fixing
  - Equivalent to running `scan` then `repair`

- **`envfix explain <issue>`** - Human-readable explanations
  - Explains what each issue means
  - Provides installation and fix instructions
  - Links to relevant documentation

## Installation

### From Source

```bash
git clone https://github.com/aicode-dev/envfix
cd envfix
go build -o envfix
```

### Prebuilt Binaries

Download from releases (coming soon)

## Quick Start

```bash
# Scan your environment for issues
envfix scan

# See if there are any problems
envfix scan --json  # Output as JSON

# Fix all detected issues automatically
envfix repair

# Preview changes without applying
envfix repair --dry-run

# Clean up unused caches and artifacts
envfix clean

# Generate environment manifest
envfix lock --output env.yaml

# Run full diagnosis and repair in one command
envfix doctor

# Get detailed explanation of an issue
envfix explain python_missing
```

## Supported Languages

### v1 (Current)

- **Python** - Virtual environments, pip, requirements.txt
- **Node.js** - npm, yarn, pnpm, node_modules
- **System tools** - PATH verification, symlink checking

### v2 (Planned)

- **Rust** - Cargo, rustup, toolchain management
- **Go** - Go modules, GOPATH, version management

## Environment Manifest (env.yaml)

The lock command generates a reproducible environment manifest:

```yaml
version: "1.0"
timestamp: 2025-12-06T12:30:45Z

platform:
  os: windows
  arch: x64
  version: "Windows 11"

python:
  version: "3.11.0"
  venv_path: ".venv"
  package_manager: "pip"

nodejs:
  version: "18.16.0"
  package_manager: "npm"

dependencies:
  requests: "2.31.0"
  express: "4.18.2"

checksums:
  requests: "sha256:abc123..."
  express: "sha256:def456..."

environment:
  PYTHON_VERSION: "3.11.0"
  NODE_ENV: "development"
```

## How It Works

### Scan

1. Detects system platform and architecture
2. Checks for Python, Node.js, Rust, Go installations
3. Verifies virtual environments and package managers
4. Scans for common issues:
   - Missing interpreters
   - Broken venvs
   - Multiple package managers
   - Corrupted caches
   - Broken symlinks in PATH

### Repair

1. Analyzes detected issues
2. Creates a reversible repair plan
3. Rebuilds virtual environments
4. Reinstalls corrupted dependencies
5. Resolves package manager conflicts
6. Logs all changes for auditing

### Clean

1. Scans for orphaned artifacts
2. Calculates freed space
3. Safely removes:
   - Cache directories
   - Backup venv directories
   - Broken symlinks
   - Stale lockfiles
4. Reports cleanup summary

## Command Reference

### Global Flags

```
-v, --verbose       Verbose logging output
-q, --quiet         Suppress non-error output
--log-file <path>   Log to file
--json              Output as JSON
--yaml              Output as YAML
```

### scan

```bash
envfix scan [flags]

Flags:
  --json              Output as JSON
  --yaml              Output as YAML
  --detailed          Show detailed scan results
  --python            Only scan Python environment
  --nodejs            Only scan Node.js environment
  --rust              Only scan Rust environment
  --go                Only scan Go environment
```

### repair

```bash
envfix repair [flags]

Flags:
  --yes               Skip confirmation prompts
  --dry-run           Preview changes without applying
  --python            Only repair Python environment
  --nodejs            Only repair Node.js environment
  --backup            Create backups before repair
```

### clean

```bash
envfix clean [flags]

Flags:
  --yes               Skip confirmation prompts
  --dry-run           Preview what would be deleted
  --python            Only clean Python caches
  --nodejs            Only clean Node.js caches
  --all-versions      Remove all old language versions
```

### lock

```bash
envfix lock [flags]

Flags:
  --output <path>     Output file (default: env.yaml)
  --json              Output as JSON instead of YAML
  --include-hashes    Include dependency hashes
  --env-vars          Include environment variables
```

## Examples

### Basic Workflow

```bash
# 1. Scan for issues
envfix scan

# 2. Review any issues found
envfix explain python_missing

# 3. Fix everything automatically
envfix repair

# 4. Clean up caches
envfix clean

# 5. Generate manifest for reproducibility
envfix lock
```

### Continuous Integration

```bash
# In your CI pipeline
envfix scan --json > scan-report.json

if [[ $? -ne 0 ]]; then
  echo "Environment check failed"
  exit 1
fi

# Then proceed with build/test
```

### Docker/Container Setup

```dockerfile
FROM ubuntu:22.04

# Install envfix
RUN go install envfix@latest

# Scan and prepare environment
RUN envfix scan && envfix repair && envfix clean

# ... rest of Dockerfile
```

## Orphaned Environments

envfix can detect and remove orphaned virtual environments and stale artifacts that consume disk space.

See [ORPHANED_ENVIRONMENTS.md](./ORPHANED_ENVIRONMENTS.md) for:
- How orphaned detection works
- Cleanup procedures
- Safety safeguards
- Configuration options
- Troubleshooting guide

**Quick example:**
```bash
# Detect orphaned venvs
envfix scan

# Remove them
envfix clean
```

## Dependency Management

envfix automatically manages project dependencies across Python, Node.js, Rust, and Go.

See [DEPENDENCIES.md](./DEPENDENCIES.md) for:
- Python: pip install and requirements.txt management
- Node.js: npm ci for clean installs, lock file validation
- Rust: cargo update and Cargo.lock management
- Go: go mod tidy and go.sum verification
- Dependency verification and validation
- Lock file management and regeneration
- Cache cleanup and optimization

**Quick example:**
```bash
# Install and verify all dependencies
envfix repair

# Verify dependency consistency
envfix repair --dry-run

# Clean dependency caches
envfix clean
```

## Version Management

envfix detects and reports version conflicts, EOL versions, and compatibility issues.

See [VERSION_MANAGEMENT.md](./VERSION_MANAGEMENT.md) for:
- Version conflict detection (multiple version managers)
- Minimum version requirement checking
- EOL version detection (Python 2, Node.js 12)
- Outdated version flagging
- Version compatibility validation
- Upgrade recommendations

**Quick example:**
```bash
# Detect version issues
envfix scan

# Validate version compatibility
envfix repair --dry-run

# See detailed explanations
envfix explain python_eol
```

## Troubleshooting

### envfix scan finds no Python

1. Python may not be installed
   ```bash
   envfix explain python_missing
   ```

2. Or it's not in your PATH
   ```bash
   which python3  # or where python3 on Windows
   ```

### Node.js package managers in conflict

Multiple lock files (package-lock.json, yarn.lock, pnpm-lock.yaml)

```bash
# Let envfix resolve it
envfix repair
```

### Repair failed

1. Check the error message
2. Run with verbose output: `envfix repair -v`
3. Check logs: `envfix repair --log-file repair.log`

### Orphaned environment issues

See [ORPHANED_ENVIRONMENTS.md](./ORPHANED_ENVIRONMENTS.md#troubleshooting)

## FAQ

**Q: Is my data safe?**
A: Yes. All repairs are reversible and logged. Backups are created before major changes.

**Q: Can I undo a repair?**
A: Yes. Renamed directories end with `.backup` or `.old` and can be restored.

**Q: Why is my scan slow?**
A: Large caches or many installed versions can slow scanning. Use language-specific flags to speed up: `envfix scan --python`

**Q: Can I use this in CI/CD?**
A: Yes. Use `--yes` flag to skip prompts and `--json` for machine-readable output.

## Contributing

Contributions welcome! Please see CONTRIBUTING.md

## License

MIT - See LICENSE file

## Support

- Issue Tracker: GitHub Issues

