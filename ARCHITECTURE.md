# envfix Architecture

## Project Structure

```
envfix/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── README.md               # User documentation
├── ARCHITECTURE.md         # This file
├── Makefile                # Build and development tasks
├── cmd/                    # CLI command implementations
│   ├── root.go            # Root command setup
│   ├── scan.go            # Scan command
│   ├── repair.go          # Repair command
│   ├── clean.go           # Clean command
│   ├── lock.go            # Lock/manifest command
│   ├── doctor.go          # Doctor command (scan + repair)
│   └── explain.go         # Explain command
├── internal/
│   ├── types/             # Core data structures
│   │   └── types.go       # Type definitions
│   ├── utils/             # Utility functions
│   │   ├── system.go      # System utilities (PATH, exec, etc)
│   │   ├── logger.go      # Logging system
│   │   └── *_test.go      # Unit tests
│   ├── detector/          # Environment detection
│   │   ├── detector.go    # Main detector orchestrator
│   │   ├── python.go      # Python detection
│   │   ├── node.go        # Node.js detection
│   │   ├── rust.go        # Rust detection
│   │   ├── go_env.go      # Go detection
│   │   └── *_test.go      # Unit tests
│   ├── rebuilder/         # Environment repair
│   │   ├── rebuilder.go   # Main rebuilder orchestrator
│   │   ├── python_repair.go # Python repair operations
│   │   └── node_repair.go   # Node.js repair operations
│   ├── cleaner/           # Cache and artifact cleanup
│   │   └── cleaner.go     # Cleanup operations
│   └── spec/              # Environment manifest generation
│       └── manifest.go    # env.yaml generation
└── .gitignore             # Git ignore rules
```

## Core Modules

### types (internal/types/types.go)

Defines all data structures:
- `ScanResult` - Complete environment scan results
- `PlatformInfo` - System information (OS, arch, home dir)
- `PythonEnv`, `NodeEnv`, `RustEnv`, `GoEnv` - Language environment states
- `Issue` - Detected problems with severity levels
- `EnvironmentManifest` - env.yaml structure
- `RepairAction` - Repair operations with reversibility metadata

### utils (internal/utils/)

Cross-platform utility functions:

**system.go:**
- `GetSystemInfo()` - Platform detection
- `FindExecutable()` - Locate programs in PATH
- `RunCommand()` - Execute external commands
- `PathExists()` - File system checks
- `IsWindows()`, `IsMacOS()`, `IsLinux()` - Platform detection
- `SplitPath()` - Parse PATH environment variable
- `ExpandPath()` - Handle ~ and environment variables
- `DirectorySize()` - Calculate directory sizes
- `CreateBackup()` - Safe file backup

**logger.go:**
- `Debug()`, `Info()`, `Warning()`, `Error()` - Logging functions
- File logging support
- Timestamp and level prefixes
- `Fatal()` - Exit with error code

### detector (internal/detector/)

**detector.go** - Main orchestrator:
- `NewDetector()` - Create detector instance
- `Scan()` - Full environment scan
- `CheckIssues()` - Issue analysis
- `CheckPackageManagerConflicts()` - Detect conflicts
- `CheckPath()` - Validate PATH entries

**python.go:**
- `DetectPython()` - Detect Python installation
- `findPythonExecutable()` - Locate python binary
- `getPythonPaths()` - Check common installation locations
- `isVenvActive()` - Check for active virtual environment

**node.go:**
- `DetectNode()` - Detect Node.js installation
- `findNodeExecutable()` - Locate node binary
- `getNodePaths()` - Check common installation locations
- `validateNodeEnv()` - Check for package manager conflicts

**rust.go:**
- `DetectRust()` - Detect Rust installation
- `findRustcExecutable()` - Locate rustc binary

**go_env.go:**
- `DetectGo()` - Detect Go installation
- `findGoExecutable()` - Locate go binary

### rebuilder (internal/rebuilder/)

**rebuilder.go** - Main orchestrator:
- `NewRebuilder()` - Create rebuilder instance
- `Repair()` - Execute environment repairs
- `planRepairs()` - Create repair action plan
- `executeRepairs()` - Run repair operations
- `SetDryRun()` - Enable preview mode

**python_repair.go:**
- `RepairPythonVenv()` - Rebuild virtual environment
- `InstallPythonDependencies()` - Install from requirements.txt
- `CheckPythonPath()` - Validate Python in PATH

**node_repair.go:**
- `RepairNodeModules()` - Rebuild node_modules
- `ResolveNodePackageManagerConflict()` - Fix multiple lock files
- `CleanNodeCache()` - Clear npm/yarn/pnpm caches

### cleaner (internal/cleaner/)

**cleaner.go:**
- `NewCleaner()` - Create cleaner instance
- `Clean()` - Execute cleanup operations
- `CleanPython()` - Remove pip/uv caches
- `CleanNode()` - Remove npm/yarn/pnpm caches
- `CleanBrokenSymlinks()` - Remove invalid symlinks
- `CleanCaches()` - Clear all language caches
- `RemoveOldVenvs()` - Delete backup environments

### spec (internal/spec/)

**manifest.go:**
- `NewGenerator()` - Create generator instance
- `Generate()` - Create EnvironmentManifest
- `GenerateYAML()` - Export as YAML bytes
- `GenerateString()` - Export as YAML string

## Data Flow

### Scan Flow

```
User runs: envfix scan
    ↓
cmd/scan.go: scanCmd.RunE()
    ↓
detector.NewDetector()
    ↓
detector.Scan()
    ├── GetSystemInfo() → PlatformInfo
    ├── DetectPython() → PythonEnv
    ├── DetectNode() → NodeEnv
    ├── DetectRust() → RustEnv
    ├── DetectGo() → GoEnv
    ├── CheckIssues()
    │   ├── CheckPackageManagerConflicts()
    │   └── CheckPath()
    └── Return ScanResult
    ↓
Display results to user
```

### Repair Flow

```
User runs: envfix repair
    ↓
cmd/repair.go: repairCmd.RunE()
    ↓
rebuilder.NewRebuilder()
    ↓
rebuilder.Repair()
    ├── detector.Scan() → ScanResult
    ├── planRepairs(ScanResult)
    │   └── Append RepairActions based on issues
    ├── executeRepairs()
    │   └── For each RepairAction:
    │       ├── Create backup
    │       ├── Execute action.Action()
    │       └── Log result
    └── Return success/error
    ↓
Log results to user
```

### Clean Flow

```
User runs: envfix clean
    ↓
cmd/clean.go: cleanCmd.RunE()
    ↓
cleaner.NewCleaner()
    ↓
cleaner.Clean()
    ├── CleanPython()
    │   ├── Remove pip cache
    │   └── Remove uv cache
    ├── CleanNode()
    │   ├── Remove npm cache
    │   ├── Remove pnpm cache
    │   └── Remove yarn cache
    ├── CleanBrokenSymlinks()
    ├── CleanCaches()
    └── Return total freed space
    ↓
Display cleanup summary
```

### Lock Flow

```
User runs: envfix lock
    ↓
cmd/lock.go: lockCmd.RunE()
    ↓
spec.NewGenerator()
    ↓
generator.GenerateYAML()
    ├── generator.Generate()
    │   ├── detector.Scan() → ScanResult
    │   └── Create EnvironmentManifest
    ├── yaml.Marshal()
    └── Return YAML bytes
    ↓
Write to env.yaml
Display to user
```

## Key Design Decisions

### 1. Modular Architecture

Each language (Python, Node, Rust, Go) has separate detection and repair modules:
- Easy to add new languages
- Each module can be developed independently
- Simple to test individual languages

### 2. Type Safety

All data structures are strongly typed:
- `ScanResult` contains typed language environments
- `Issue` objects have severity and language fields
- `RepairAction` includes reversibility metadata

### 3. Reversibility

Core operations are designed to be reversible:
- Backups created before modifications
- Rename strategy allows rollback (.old, .backup extensions)
- All actions logged for auditing

### 4. Dry-Run Support

Repair operations support preview mode:
- Plans can be generated without execution
- Users can review changes before applying
- Useful for CI/CD pipelines

### 5. Cross-Platform

Utilities handle OS differences:
- `IsWindows()`, `IsMacOS()`, `IsLinux()` checks
- Path separator handling (`;` vs `:`)
- Platform-specific executable names and paths

### 6. Error Handling

Consistent error reporting:
- Issues collected during scan
- Severity levels (critical, warning, info)
- Detailed messages with fix suggestions

## Extension Points

### Adding a New Language (e.g., Ruby)

1. Create `internal/detector/ruby.go`:
   ```go
   func (d *Detector) DetectRuby() types.RubyEnv {
       // Detection logic
   }
   ```

2. Create `internal/rebuilder/ruby_repair.go`:
   ```go
   func (r *Rebuilder) RepairRubyEnv() error {
       // Repair logic
   }
   ```

3. Add types to `internal/types/types.go`:
   ```go
   type RubyEnv struct {
       Version    string
       Executable string
       Issues     []string
   }
   ```

4. Update `Detector.Scan()` to call `DetectRuby()`

5. Update `Rebuilder.planRepairs()` to handle Ruby issues

### Adding a New Repair Operation

1. Implement method on `Rebuilder`:
   ```go
   func (r *Rebuilder) FixNewIssue() error {
       // Implementation
   }
   ```

2. Add to `planRepairs()`:
   ```go
   r.actions = append(r.actions, types.RepairAction{
       ID:          "new_issue",
       Description: "Fix new issue",
       Action:      r.FixNewIssue,
       Reversible:  true,
   })
   ```

### Adding a New Cleanup Operation

1. Implement method on `Cleaner`:
   ```go
   func (c *Cleaner) CleanNewCache() (int64, error) {
       // Return bytes freed
   }
   ```

2. Call from `Clean()`:
   ```go
   size, err := c.CleanNewCache()
   cleanedSize += size
   ```

## Testing Strategy

### Unit Tests

Located in same package with `_test.go` suffix:
- `internal/utils/system_test.go` - System utility tests
- `internal/detector/detector_test.go` - Detector tests
- Test coverage > 80% target

### Test Categories

1. **Platform Tests** - Verify cross-platform compatibility
2. **Detection Tests** - Verify environment detection
3. **Repair Tests** - Verify repair operations (reversible)
4. **Utility Tests** - Verify utility functions

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...

# Specific package
go test ./internal/detector/...
```

## Performance Targets

- Scan: < 5 seconds on standard machines
- Repair: < 30 seconds
- Clean: < 10 seconds
- Lock: < 2 seconds

Achieved through:
- Parallel language detection
- Minimal file I/O
- Efficient caching
- Direct binary execution vs parsing

## Security Considerations

1. **No elevated privileges** - Runs as regular user unless necessary
2. **Backups before changes** - All modifications reversible
3. **Audit logging** - All changes logged
4. **Safe cleanup** - Only removes known caches/backups
5. **Input validation** - Command outputs sanitized
6. **No remote calls** - Offline-first design

## Future Enhancements

### v2
- Rust and Go full support
- Plugin architecture
- Environment history and rollback

### v3
- Cloud manifest sync
- Team collaboration
- AI-powered explanations

### v4
- IDE integrations
- Real-time monitoring
- Advanced analytics
