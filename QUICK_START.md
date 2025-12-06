# envfix Quick Start Guide

## Installation

### Windows
```bash
# Build from source
go build -o envfix.exe

# Or download from releases (coming soon)
```

### macOS/Linux
```bash
# Build from source
go build -o envfix

# Or download from releases (coming soon)
```

## Basic Usage

### 1. Scan Your Environment
```bash
# Quick scan
envfix scan

# With verbose output
envfix scan --verbose

# Save to log file
envfix scan --log-file scan.log
```

**Output:** Shows all detected languages, versions, and any issues.

### 2. Generate Environment Manifest
```bash
# Create env.yaml
envfix lock

# Save to different file
envfix lock --output my-environment.yaml
```

**Output:** Creates `env.yaml` with all environment details for reproducibility.

### 3. Clean Caches and Artifacts
```bash
# Preview what will be cleaned
envfix clean --dry-run

# Actually clean
envfix clean
```

**Output:** Removes old caches, freed space reported.

### 4. Fix Environment Issues
```bash
# Preview repairs (dry-run)
envfix repair --dry-run

# Actually repair (with backups)
envfix repair
```

**Output:** Rebuilds broken environments, creates backups.

### 5. One-Shot Diagnosis
```bash
# Scan + repair in one command
envfix doctor
```

**Output:** Complete diagnosis and repair.

### 6. Get Help on Issues
```bash
envfix explain python_missing
envfix explain package_manager_conflict
envfix explain broken_symlink
```

**Output:** Detailed explanation and how to fix.

### 7. Setup Configuration
```bash
# Create local .envfix.yaml
envfix init

# Create global ~/.envfix/config.yaml
envfix init --global
```

**Output:** Configuration files for customization.

## Common Scenarios

### Scenario: Fresh Setup
```bash
# 1. Scan environment
envfix scan

# 2. Generate manifest
envfix lock

# 3. Save manifest for team
git add env.yaml
git commit -m "Add environment manifest"
```

### Scenario: Troubleshooting Environment
```bash
# 1. Scan for issues
envfix scan

# 2. Understand issues
envfix explain <issue_name>

# 3. Fix automatically
envfix repair --dry-run  # Preview first
envfix repair            # Apply fixes

# 4. Clean up
envfix clean
```

### Scenario: CI/CD Integration
```bash
# In your CI pipeline
envfix scan || exit 1  # Fail if issues
envfix clean --yes     # Clean without prompts
# Continue with build/test
```

### Scenario: Team Synchronization
```bash
# Person 1: Generate manifest
envfix lock

# Commit and push
git add env.yaml
git commit -m "Update environment manifest"
git push

# Person 2: Check if environment matches
envfix scan  # Compare with env.yaml
```

## Commands Reference

| Command | Purpose | Typical Use |
|---------|---------|------------|
| `scan` | Diagnose environment | Find issues |
| `repair` | Fix detected issues | Auto-fix broken setup |
| `clean` | Remove caches/artifacts | Free disk space |
| `lock` | Generate env.yaml | Save environment snapshot |
| `doctor` | Scan + repair | One-shot diagnosis |
| `explain` | Get help on issue | Understand problems |
| `init` | Setup configuration | Customize behavior |

## Global Flags

All commands support:
```bash
--verbose        # Detailed output
--log-file FILE  # Log to file
--help          # Show help
```

## Tips & Tricks

### Speed Up Scan
```bash
# Only scan Python
envfix scan --python

# Only scan Node.js
envfix scan --node
```

### Dry Run Before Changes
```bash
# Always preview repairs first
envfix repair --dry-run
```

### Save Output to File
```bash
# Log everything
envfix scan --log-file scan.log

# Check later
cat scan.log
```

### Custom Configuration
```bash
# Create and edit config
envfix init
vim .envfix.yaml  # or your editor

# Customize:
# - repair strategy
# - excluded directories
# - cache cleanup
# - environment variables
```

### Keep Manifests for History
```bash
# Generate dated manifests
envfix lock --output env.2025-12-06.yaml

# Compare over time
diff env.2025-12-06.yaml env.2025-12-07.yaml
```

## Troubleshooting

### "Python not found"
```bash
envfix explain python_missing
# Follow installation instructions
```

### "Package manager conflict"
```bash
envfix repair --dry-run  # Preview
envfix repair            # Fix (removes conflicting files)
```

### "Broken symlink"
```bash
envfix clean  # Removes broken symlinks
```

### Need detailed help?
```bash
envfix <command> --help
# Or see README.md and ARCHITECTURE.md
```

## Next Steps

1. **Read README.md** - Full documentation
2. **Check TESTING.md** - How to run tests
3. **See DEVELOPMENT.md** - If contributing
4. **Review ARCHITECTURE.md** - How it works

## Support

- GitHub Issues: Report bugs
- GitHub Discussions: Ask questions
- Check existing docs first

## Installation Issues?

If you can't build:
1. Verify Go 1.21+ installed: `go version`
2. Check Go path: `go env GOPATH`
3. Clean and rebuild: `go clean && go build`
4. See DEVELOPMENT.md for troubleshooting

## Common Commands

```bash
# Show all available commands
envfix --help

# Get specific command help
envfix scan --help
envfix repair --help
envfix clean --help

# Quick version check
envfix --version  # (coming in v1.1)
```

## Examples

### Full Workflow
```bash
# 1. Initial setup
envfix init
envfix scan

# 2. Fix issues
envfix repair --dry-run
envfix repair

# 3. Clean and document
envfix clean
envfix lock

# 4. Share setup
git add .envfix.yaml env.yaml
git commit -m "Setup reproducible environment"
```

### Continuous Monitoring
```bash
# Check daily
envfix scan > daily-$(date +%Y-%m-%d).txt

# Monitor changes
diff daily-2025-12-06.txt daily-2025-12-07.txt
```

### Team Collaboration
```bash
# Setup shared environment
envfix init --global

# Generate team manifest
envfix lock --output team-env.yaml

# Share with team
git add team-env.yaml
```

---

**For more details:** See README.md, ARCHITECTURE.md, and TESTING.md
