# Orphaned Environment Detection & Cleanup Guide

## Overview

Orphaned environments are virtual environments, package installations, and toolchain versions that exist on your system but are no longer connected to any active project. These can accumulate over time and waste significant disk space.

envfix v1.0.1 adds comprehensive detection and cleanup for:
- Orphaned Python virtual environments (.venv, venv, .virtualenvs)
- Stale node_modules without parent projects
- Multiple language toolchain versions (Go, Rust)
- Old backup directories from previous repairs

## Detection

### What Gets Detected

#### Python Virtual Environments
A venv is considered **orphaned** if:
- It exists in a standard location (.venv, venv, .virtualenvs, etc.)
- The parent directory has NO Python project files:
  - requirements.txt
  - pyproject.toml
  - setup.py
  - poetry.lock
  - Pipfile

#### Node.js Modules
node_modules are considered **orphaned** if:
- They exist in common project locations
- The parent directory has NO package.json file
- They're old backups (node_modules.backup, node_modules.old)

#### Toolchain Versions
Multiple versions detected for:
- **Go**: Via gimme or version managers
- **Rust**: Multiple rustup toolchains (>2 installed)
- **Python**: pyenv, conda, system Python all present

### Using Scan Command

```bash
# Detect orphaned environments
envfix scan

# The output will show:
# Issues: [ID: orphaned_venv, Title: "Orphaned Python virtual environment detected"]
# Warnings: ["Orphaned venv: /home/user/.venv (3 days old)"]
```

### Example Output

```
Environment Scan Results
========================

Platform: linux (amd64)

Python 3.11.0
├── Status: healthy
└── Issues: 1
    └── Orphaned venv: /home/user/old-project/.venv (14 days old)

Node.js v18.12.0
├── Status: healthy
└── Warnings: Stale node_modules at /home/user/archived-project/node_modules (2 months old)

Warnings (3):
- Orphaned venv: /home/user/.venv (3 days old)
- Multiple Go versions detected in /home/user/.gimme/versions (oldest: 7 days)
- Orphaned venv: /home/user/experiments/venv (45 days old)

Health: ⚠ Issues found (3 warnings)
Scan Time: 2.34s
```

## Cleanup

### Using Clean Command

```bash
# Run cleanup with all orphaned environments removed
envfix clean

# The output will show:
# Cleaning orphaned environments...
# ✓ Removed orphaned venv: old-project (125.3 MB)
# ✓ Cleanup complete (freed ~125.3 MB)
```

### Using Repair Command

```bash
# Repair will automatically remove orphaned environments if detected
envfix repair

# Output:
# [1/2] Remove orphaned virtual environments
# ✓ Completed: Remove orphaned virtual environments
```

### Using Doctor Command

```bash
# Full scan + repair (includes orphaned cleanup)
envfix doctor
```

## How It Works

### Detection Algorithm

1. **Scan Phase**:
   - Enumerate standard venv locations:
     - ~/.venv
     - ~/venv
     - ~/.virtualenvs
     - ~/.local/share/venvs
     - ~/.pyenv/versions
   
2. **Verification Phase**:
   - For each directory, check if it's a real venv:
     - Look for pyvenv.cfg file
   - Check parent directory for project markers
   
3. **Age Calculation**:
   - Get last modified time
   - Calculate age (hours, days, weeks, months, years)
   - Report to user

4. **Toolchain Scanning**:
   - Check ~/.gimme/versions for Go versions
   - Check ~/.rustup/toolchains for Rust versions
   - Check for multiple version managers (pyenv, conda, nvm, etc.)

### Cleanup Algorithm

1. **Validation**:
   - Confirm it's an orphaned venv (pyvenv.cfg exists, no parent project)
   - Calculate disk space to be freed
   
2. **Removal**:
   - Safely remove directory
   - Log removal action
   - Accumulate freed space
   
3. **Report**:
   - Show total disk space freed
   - List removed environments

## Configuration

### Local Config (.envfix.yaml)

```yaml
cleanup:
  remove_orphaned: true
  min_age_days: 7  # Only remove orphaned envs older than 7 days
  exclude_patterns:
    - "**/work/**"    # Don't remove work project venvs
    - "**/projects/**"
```

### Global Config (~/.envfix/config.yaml)

```yaml
cleanup:
  aggressive: true        # Remove all orphaned environments
  interactive: false      # Don't ask for confirmation
  exclude_venv_names:
    - "test_venv"
    - "experimental"
```

## Safety Features

### Backups and Reversibility
- Orphaned environment removal is NOT reversible
- No backup is created (orphaned envs are unused)
- Clean command shows what will be removed before execution

### Safeguards
- Must have pyvenv.cfg to be considered a venv
- Must have NO project files to be considered orphaned
- Age checking prevents removing newly created venvs
- Symlink checking prevents removing active venv links

### Manual Review
```bash
# Preview what would be removed (dry-run mode)
envfix clean --dry-run

# Shows:
# Would remove: /home/user/old-project/.venv (125.3 MB)
# Would remove: /home/user/.venv (45.2 MB)
# Total would free: ~170.5 MB
```

## Common Scenarios

### Scenario 1: Archived Projects

**Problem**: You have 5 old project directories with venvs taking up 500 MB

```bash
$ envfix scan
Warnings:
- Orphaned venv: /home/user/archived/project1/.venv (8 months old)
- Orphaned venv: /home/user/archived/project2/.venv (6 months old)
...

$ envfix clean
Cleaning orphaned environments...
✓ Removed orphaned venv: project1 (125.3 MB)
✓ Removed orphaned venv: project2 (145.2 MB)
...
✓ Cleanup complete (freed ~500 MB)
```

### Scenario 2: Experiment Directories

**Problem**: Many test/experiment venvs from development work

```bash
$ find ~/.venv* -type d -mtime +30  # Find venvs older than 30 days

$ envfix clean  # Automatically removes them
```

### Scenario 3: Multiple Versions

**Problem**: 3 versions of Go and 5 Rust toolchains installed

```bash
$ envfix scan
Warnings:
- Multiple Go versions detected in ~/.gimme/versions (oldest: 1 year)
- Multiple Rust toolchains detected (oldest: 8 months)

# Manual cleanup needed for toolchains (envfix doesn't auto-remove versions)
# Recommend: rustup uninstall <toolchain>
# Recommend: Remove old Go versions manually
```

## Troubleshooting

### "Orphaned venv detected for active project"

**Solution**: The venv is in a non-standard location or the project files were deleted.
```bash
# Manually check the directory
cd /path/to/venv/..
ls -la  # Look for requirements.txt, pyproject.toml, etc.

# If project is truly orphaned, it's safe to remove
rm -rf /path/to/venv
```

### "False positive - marked as orphaned but I need it"

**Solution**: Move venv back to standard location or add project files:
```bash
# Create a minimal marker file
touch pyproject.toml  # Now it won't be considered orphaned
```

### "envfix won't remove a venv I want deleted"

**Solution**: Check it has pyvenv.cfg and no project files:
```bash
# Manually verify
cat /path/to/venv/pyvenv.cfg  # Should exist
ls -la /path/to/venv/..       # Should have no requirements.txt, setup.py, etc.

# If conditions met, file a bug
# Otherwise, remove manually:
rm -rf /path/to/venv
```

## Performance Impact

- **Detection**: ~100ms for home directory scan
- **Cleanup**: ~1s per 100 MB of orphaned environments
- **Memory**: <5 MB overhead
- **Disk I/O**: Minimal (directory listing only during detection)

## Future Enhancements (v1.1+)

- [ ] Persistent "do not remove" list for specific venvs
- [ ] Age-based removal policies (auto-remove after 6 months)
- [ ] Dry-run report generation
- [ ] Integration with system package managers
- [ ] Orphaned dependency cache cleanup
- [ ] Node.js node_modules optimization (not full removal)

## Related Commands

```bash
envfix scan      # Detect orphaned environments
envfix clean     # Remove orphaned environments
envfix repair    # Fix issues (includes orphaned removal)
envfix doctor    # Full diagnostics (scan + repair)
envfix explain   # Get detailed info on specific issues
```

## See Also

- [DEVELOPMENT.md](./DEVELOPMENT.md) - Development guide
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Technical architecture
- [README.md](./README.md) - Main documentation
- [TESTING.md](./TESTING.md) - Testing guide
