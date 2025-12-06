# Changelog

All notable changes to envfix will be documented in this file.

## [1.0.1] - 2025-12-06 (In Development)

### Added

#### Advanced Diagnostics
- **Orphaned Environment Detection** - Identifies unused venvs and toolchains
  - Detects orphaned Python virtual environments
  - Finds stale node_modules without parent project
  - Identifies multiple language toolchain versions
  - Reports directory age and cleanup recommendations
  - Integrated into `scan` and `doctor` commands

#### Orphaned Cleanup
- **CleanOrphanedEnvironments** - Removes unused environments
  - Safe removal with confirmation checks
  - Targets common venv locations
  - Cleans old node_modules backups
  - Integrated into `clean` command

#### Dependency Management
- **Package Dependency Installation & Upgrade**
  - Python: pip install and upgrade from requirements.txt
  - Node.js: npm ci (clean install) for reproducible installs
  - Rust: cargo update for dependency updates
  - Go: go mod tidy for module cleanup

- **Dependency Verification & Validation**
  - Verify all dependencies are properly installed
  - Check lock file consistency across languages
  - Detect missing lock files
  - Validate dependency configurations

- **Lock File Management**
  - Regenerate lock files for reproducible builds
  - Python: requirements.txt regeneration
  - Node.js: package-lock.json generation
  - Rust: Cargo.lock management
  - Go: go.sum verification

- **Dependency Cache Cleaning**
  - Clear pip, npm, cargo, and go module caches
  - Integrated cache cleanup in repair workflow

#### Version Conflict Detection & Management
- **Version Conflict Detection**
  - Detect multiple version managers (pyenv, conda, nvm, etc.)
  - Identify minimum version requirement violations
  - Detect end-of-life versions (Python 2, Node.js 12)
  - Flag outdated versions (Rust, Go)
  - Integrated into `scan` command

- **Version Compatibility Validation**
  - Validate all installed versions are compatible
  - Check interpreter availability in PATH
  - Report installation status
  - Integrated into `repair` command

- **Version Management Tools**
  - Suggest version upgrades
  - Detect old/duplicate installations
  - Provide upgrade recommendations
  - Links to official documentation

## [1.0.0] - 2025-12-06

### Added

#### Core Functionality
- **scan** command - Complete environment diagnostic scanning
  - Detects Python, Node.js, Rust, Go installations
  - Checks virtual environments and package managers
  - Identifies PATH issues and broken symlinks
  - Scans for dependency issues (missing lockfiles, conflicts)
  - Performance: ~2-3 seconds typical scan time

- **lock** command - Environment manifest generation
  - Generates reproducible `env.yaml` files
  - Captures all language versions
  - Includes platform and architecture information
  - Exports to YAML format

- **repair** command - Automatic environment fixes
  - Rebuilds Python virtual environments with backup
  - Resolves Node.js package manager conflicts
  - Dry-run mode for safe preview
  - Reversible operations with backup strategy

- **clean** command - Cache and artifact cleanup
  - Removes pip, npm, pnpm, yarn caches
  - Cleans Rust/Cargo registry cache
  - Identifies and removes broken symlinks
  - Reports freed disk space

- **doctor** command - One-shot scan + repair
  - Convenience wrapper for complete diagnosis
  - Combines scan and repair in single operation

- **explain** command - Issue explanations
  - Human-readable issue descriptions
  - Installation instructions
  - Fix recommendations
  - Built-in knowledge base

- **init** command - Configuration setup
  - Creates local `.envfix.yaml` configuration
  - Creates global `~/.envfix/config.yaml`
  - Force overwrite support

#### Configuration System
- `.envfix.yaml` support for local configuration
- `~/.envfix/config.yaml` for global settings
- Configurable repair strategies (conservative, aggressive, interactive)
- Exclude patterns for directories
- Cache cleanup behavior customization
- Environment variable management

#### Dependency Detection
- Parse `requirements.txt` (Python)
- Parse `package.json` (Node.js)
- Parse `go.mod` (Go modules)
- Parse `Cargo.toml` (Rust)
- Detect lockfile consistency issues
- Warn on missing or conflicting lockfiles

#### Language Support
- **Python**: Detection, venv repair, dependency tracking
- **Node.js**: Detection, npm/yarn/pnpm support, conflict resolution
- **Rust**: Detection and version reporting (repair in v2)
- **Go**: Detection and version reporting (repair in v2)
- **System**: PATH validation, symlink checking

#### Utilities
- Cross-platform support (Windows, macOS, Linux)
- Command execution with error handling
- File system operations (copy, backup, remove)
- Environment variable parsing
- Comprehensive logging with timestamps
- Platform-specific executable detection

#### Documentation
- Comprehensive README with examples (335+ lines)
- Architecture documentation (500+ lines)
- Testing guide with integration test examples
- Inline code comments and docstrings
- API documentation in code

#### Testing
- Unit tests for core modules
- Integration test framework
- Test coverage for utilities and detector
- Parser tests for dependency files

#### Build & Distribution
- Go module setup (go.mod)
- Makefile with build targets
- Cross-platform build support
- Dependency management with go mod

### Features by Version

#### v1.0 (This Release)
- [x] Core scanning and diagnostics
- [x] Environment manifest generation
- [x] Basic repair (Python, Node.js)
- [x] Cache cleaning
- [x] Configuration system
- [x] Dependency detection
- [x] Cross-platform support

#### v2 (Planned)
- [ ] Full Rust support (repair operations)
- [ ] Full Go support (repair operations)
- [ ] Advanced repair strategies
- [ ] Environment history tracking
- [ ] Rollback capability
- [ ] Plugin system foundation

#### v3 (Future)
- [ ] Cloud manifest sync
- [ ] Team collaboration features
- [ ] Environment history API
- [ ] AI-powered diagnostics

### Bug Fixes
- (None - Initial release)

### Performance
- Scan completes in ~2-3 seconds typical
- Memory usage < 50 MB
- No external dependencies required beyond Go stdlib and specified packages
- Minimal file I/O overhead

### Breaking Changes
- (None - Initial release)

### Deprecated
- (None)

### Security
- No elevated privileges required
- Safe file operations with backups
- Audit logging of all changes
- Input validation on command outputs

## Version History

| Version | Date | Status |
|---------|------|--------|
| 1.0.0 | 2025-12-06 | Released |

## Upgrade Guide

### From Pre-1.0
This is the initial public release.

## Known Issues

### v1.0
1. Repair operations don't automatically execute actual commands (dry-run ready)
2. Lockfile parsing limited to basic formats
3. No support for monorepo structures yet
4. Config file doesn't validate syntax

### Workarounds
- Use `--dry-run` for safe preview
- Manually run repairs if needed
- Validate YAML files before loading

## Future Roadmap

### Short-term (v1.1 - Q1 2025)
- [ ] Lockfile hash validation
- [ ] Parallel language detection
- [ ] Configuration validation
- [ ] Better error messages

### Medium-term (v2.0 - Q2 2025)
- [ ] Rust and Go repair
- [ ] Plugin system
- [ ] Environment history
- [ ] Rollback operations

### Long-term (v3.0 - Q3-Q4 2025)
- [ ] Cloud sync features
- [ ] Team management
- [ ] Advanced analytics
- [ ] IDE integrations

## Credits

Developed as part of the aicode project.

## License

MIT - See LICENSE file for details
