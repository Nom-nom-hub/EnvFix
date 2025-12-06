# Developer Quick Start - envfix v1.0.1

## What's New in v1.0.1?

### Orphaned Environment Detection and Cleanup

Automatically detect and remove unused virtual environments to free up disk space.

## Quick Examples

### Detect Orphaned Environments

```bash
# Scan and report orphaned venvs
envfix scan

# Output will show:
# Issues:
#   [orphaned_venv] Orphaned Python virtual environment detected
#     /home/user/old-project/.venv (14 days old)
```

### Remove Orphaned Environments

```bash
# Clean up (removes orphaned venvs and old backups)
envfix clean

# Output:
# Cleaning orphaned environments...
# ✓ Removed orphaned venv: old-project (125.3 MB)
# ✓ Cleanup complete (freed ~125.3 MB)
```

### Full Repair (Including Orphaned Cleanup)

```bash
# Scan, fix issues, and clean up
envfix doctor

# Combines scan + repair + clean in one command
```

## For Developers

### Module Locations

```
Detector Module:       internal/detector/orphaned.go
Rebuilder Module:      internal/rebuilder/orphaned_repair.go
Cleaner Module:        internal/cleaner/cleaner.go (enhanced)
Tests:                 internal/detector/orphaned_test.go
Documentation:         ORPHANED_ENVIRONMENTS.md
```

### Key Functions

#### Detection (detector package)
```go
func (d *Detector) CheckOrphanedEnvironments(result *types.ScanResult)
func (d *Detector) checkOrphanedPythonVenvs(result *types.ScanResult)
func (d *Detector) checkOrphanedNodeModules(result *types.ScanResult)
func (d *Detector) checkOrphanedToolchains(result *types.ScanResult)
func (d *Detector) isOrphanedVenv(venvPath string) bool
func (d *Detector) getDirectoryAge(dirPath string) string
```

#### Cleanup (cleaner package)
```go
func (c *Cleaner) CleanOrphanedEnvironments() (int64, error)
func (c *Cleaner) isOrphanedVenv(venvPath string) bool
```

#### Repair (rebuilder package)
```go
func (r *Rebuilder) RemoveOrphanedEnvironments() error
func (r *Rebuilder) isOrphanedVenv(venvPath string) bool
func (r *Rebuilder) ValidateEnvironmentToolchains() error
```

### Adding Tests

```go
// Create test in internal/detector/orphaned_test.go
func TestNewFeature(t *testing.T) {
    d := &Detector{}
    // Your test here
}

// Run tests
go test ./internal/detector -v
```

### Integration Points

1. **Scan Workflow** (detector/detector.go:62-64)
   ```go
   // Check for orphaned environments
   d.CheckOrphanedEnvironments(result)
   ```

2. **Clean Workflow** (cleaner/cleaner.go:51-58)
   ```go
   size, err = c.CleanOrphanedEnvironments()
   if err != nil {
       utils.Warning("Orphaned environment cleanup failed: %v", err)
   } else {
       cleanedSize += size
   }
   ```

3. **Repair Workflow** (rebuilder/rebuilder.go:110-128)
   ```go
   if hasOrphanedIssues {
       r.actions = append(r.actions, types.RepairAction{
           ID:          "orphaned_envs",
           Description: "Remove orphaned virtual environments",
           Action:      r.RemoveOrphanedEnvironments,
           Reversible:  false,
       })
   }
   ```

## How It Works

### Detection Algorithm

1. **Scan standard venv locations:**
   - ~/.venv
   - ~/venv
   - ~/.virtualenvs
   - ~/.local/share/venvs
   - ~/.pyenv/versions

2. **For each directory:**
   - Check for pyvenv.cfg (confirms it's a venv)
   - Check parent for project files (requirements.txt, pyproject.toml, etc.)
   - If no project files found = orphaned
   - Calculate age and report

3. **Also detect:**
   - Multiple Go versions
   - Multiple Rust toolchains
   - Stale node_modules without projects

### Cleanup Algorithm

1. **For each orphaned venv:**
   - Verify it has pyvenv.cfg
   - Verify no project files exist
   - Calculate disk space
   - Remove directory
   - Log removal

2. **Also clean:**
   - node_modules.backup
   - node_modules.old
   - .node_modules.backup

## Safety Features

### Verification Before Removal
- Must have pyvenv.cfg
- Must have NO project files
- Prevents false positives

### Age-Based Safeguards
- Reports age to user
- Avoids very new venvs

### Project File Detection
Detects these markers (not orphaned if present):
- requirements.txt (pip)
- pyproject.toml (poetry)
- setup.py (setuptools)
- poetry.lock (poetry)
- Pipfile (pipenv)

## Configuration

### Local Config (.envfix.yaml)

```yaml
cleanup:
  remove_orphaned: true
  min_age_days: 7
  exclude_patterns:
    - "**/work/**"
```

### Global Config (~/.envfix/config.yaml)

```yaml
cleanup:
  aggressive: true
  interactive: false
```

## Common Development Tasks

### Running Tests

```bash
# All tests
go test ./...

# Orphaned module tests only
go test ./internal/detector -v -run Orphaned

# With coverage
go test ./internal/detector -v -cover
```

### Building

```bash
# Clean build
go build -o envfix

# With version info
go build -ldflags="-X main.version=1.0.1" -o envfix

# Cross-platform
GOOS=linux GOARCH=amd64 go build -o envfix-linux
GOOS=darwin GOARCH=amd64 go build -o envfix-darwin
```

### Running Commands

```bash
# Build and test
go build && ./envfix scan
./envfix scan --json | jq '.warnings'
./envfix clean --dry-run
./envfix repair
```

## Extending the Feature

### Add New Venv Location

Edit `internal/detector/orphaned.go`:
```go
venvLocations := []string{
    filepath.Join(home, ".venv"),
    filepath.Join(home, "venv"),
    // Add new location:
    filepath.Join(home, "custom-venvs"),  // <-- New
}
```

### Add New Project Marker

Edit both files where `projectFiles` is defined:
```go
projectFiles := []string{
    filepath.Join(parentDir, "requirements.txt"),
    // Add new marker:
    filepath.Join(parentDir, ".python-version"),  // <-- New
}
```

### Add New Toolchain Check

Edit `internal/detector/orphaned.go`:
```go
// In checkOrphanedToolchains()
// Add new check for your language
pythonVersionDirs := []string{
    filepath.Join(home, ".pyenv", "versions"),
    // Add new location:
    filepath.Join(home, ".custom-versions"),  // <-- New
}
```

## Performance Tips

- Detection is fast (~100ms) - scans are quick
- Cleanup is proportional to size (~1s per 100MB)
- Use `--dry-run` to preview before cleaning
- Consider excluding large directories if needed

## Debugging

### Enable Verbose Logging

```bash
envfix scan --verbose
envfix clean --verbose
envfix repair --verbose
```

### Check What Would Be Removed

```bash
envfix clean --dry-run
```

### Manual Verification

```bash
# Find orphaned venvs manually
find ~/.venv* -type d -name "pyvenv.cfg" 2>/dev/null

# Check for projects
ls -la ~/.venv/..
```

## Related Documentation

- [ORPHANED_ENVIRONMENTS.md](./ORPHANED_ENVIRONMENTS.md) - Full user guide
- [PHASE_3_SUMMARY.md](./PHASE_3_SUMMARY.md) - Implementation details
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Project architecture
- [DEVELOPMENT.md](./DEVELOPMENT.md) - Development guide

## FAQ for Developers

**Q: Where is the main detection logic?**  
A: `internal/detector/orphaned.go` - `CheckOrphanedEnvironments()` function

**Q: How do I test changes?**  
A: `go test ./internal/detector -v -run TestYourTest`

**Q: Can I add a new language?**  
A: Yes! Extend `checkOrphanedToolchains()` in orphaned.go

**Q: What if I want to preserve a specific venv?**  
A: Add a marker file (requirements.txt, pyproject.toml) to its parent directory

**Q: How is age calculated?**  
A: Using directory modification time via `os.Stat()` and `time.Since()`

## Next Steps

1. Build and test the project
2. Run `envfix scan` to see detection in action
3. Run `envfix clean --dry-run` to preview cleanups
4. Review ORPHANED_ENVIRONMENTS.md for user perspective
5. Proceed to Phase 4 (Package Reinstallation)

---

**Version**: 1.0.1 (Development)  
**Status**: Ready for testing and integration
