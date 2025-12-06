package detector

import (
	"envfix/internal/types"
	"envfix/internal/utils"
	"fmt"
	"time"
)

type Detector struct{}

func NewDetector() *Detector {
	return &Detector{}
}

// Scan performs system and project checks for environment issues
func (d *Detector) Scan() (*types.ScanResult, error) {
	startTime := time.Now()

	result := &types.ScanResult{
		Timestamp: time.Now(),
		Issues:    []types.Issue{},
		Warnings:  []string{},
	}

	// Detect platform
	osName, arch, homeDir := utils.GetSystemInfo()
	result.Platform = types.PlatformInfo{
		OS:      osName,
		Arch:    arch,
		HomeDir: homeDir,
		PathSep: utils.GetPathSeparator(),
	}

	// Detect Python
	result.Python = d.DetectPython()

	// Detect Node.js
	result.Node = d.DetectNode()

	// Detect Rust
	result.Rust = d.DetectRust()

	// Detect Go
	result.Go = d.DetectGo()

	// Check for issues
	d.CheckIssues(result)

	// Check for dependency issues
	d.DetectDependencyIssues(result)

	// Check for corruption issues
	d.CheckCorruptionIssues(result)

	// Check for version conflicts
	d.CheckVersionConflicts(result)

	// Check for version mismatches
	versionIssues := d.DetectVersionMismatches(result)
	result.Issues = append(result.Issues, versionIssues...)

	// Check for orphaned environments
	d.CheckOrphanedEnvironments(result)

	// Calculate health status
	result.IsHealthy = len(result.Issues) == 0
	result.ScanDuration = time.Since(startTime).Seconds()

	return result, nil
}

// CheckCorruptionIssues detects corrupted environments
func (d *Detector) CheckCorruptionIssues(result *types.ScanResult) {
	// Check Python venv corruption
	if result.Python.Status == "healthy" && result.Python.VenvFound {
		venvIssues := d.CheckVenvCorruption(result.Python.VenvPath)
		for _, issue := range venvIssues {
			result.Python.Issues = append(result.Python.Issues, issue)
			result.Issues = append(result.Issues, types.Issue{
				ID:       "python_venv_corrupt",
				Title:    "Python venv corrupted",
				Severity: "critical",
				Language: "python",
				Message:  issue,
				Fix:      "Run 'envfix repair' to rebuild the virtual environment",
			})
		}
	}

	// Check Node.js modules corruption
	if result.Node.Status == "healthy" && result.Node.ModulesFound {
		nodeIssues := d.CheckNodeModulesCorruption()
		for _, issue := range nodeIssues {
			result.Node.Issues = append(result.Node.Issues, issue)
			result.Issues = append(result.Issues, types.Issue{
				ID:       "node_modules_corrupt",
				Title:    "node_modules corrupted",
				Severity: "warning",
				Language: "nodejs",
				Message:  issue,
				Fix:      "Run 'envfix repair' to rebuild node_modules",
			})
		}
	}

	// Check for global package pollution
	pollutionIssues := d.CheckGlobalPackagePollution()
	result.Warnings = append(result.Warnings, pollutionIssues...)

	// Check cache corruption
	cacheIssues := d.CheckCacheCorruption()
	result.Warnings = append(result.Warnings, cacheIssues...)
}

// CheckIssues analyzes detected environments for problems
func (d *Detector) CheckIssues(result *types.ScanResult) {
	// Python issues
	if result.Python.Status == "missing" {
		for _, issue := range result.Python.Issues {
			result.Issues = append(result.Issues, types.Issue{
				ID:       "python_missing",
				Title:    "Python not found",
				Severity: "critical",
				Language: "python",
				Message:  issue,
			})
		}
	} else if len(result.Python.Issues) > 0 {
		for _, issue := range result.Python.Issues {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Python: %s", issue))
		}
	}

	// Node.js issues
	if result.Node.Status == "missing" {
		for _, issue := range result.Node.Issues {
			result.Issues = append(result.Issues, types.Issue{
				ID:       "nodejs_missing",
				Title:    "Node.js not found",
				Severity: "critical",
				Language: "nodejs",
				Message:  issue,
			})
		}
	} else if len(result.Node.Issues) > 0 {
		for _, issue := range result.Node.Issues {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Node.js: %s", issue))
		}
	}

	// Check for mixed package managers
	d.CheckPackageManagerConflicts(result)

	// Check PATH
	d.CheckPath(result)
}

func (d *Detector) CheckPackageManagerConflicts(result *types.ScanResult) {
	if len(result.Node.Issues) > 0 {
		for _, issue := range result.Node.Issues {
			if issue == "Multiple package manager lock files detected (npm, pnpm, yarn)" {
				result.Issues = append(result.Issues, types.Issue{
					ID:       "package_manager_conflict",
					Title:    "Package manager conflict",
					Severity: "warning",
					Language: "nodejs",
					Message:  issue,
				})
			}
		}
	}
}

func (d *Detector) CheckPath(result *types.ScanResult) {
	paths := utils.SplitPath()
	if len(paths) == 0 {
		result.Warnings = append(result.Warnings, "PATH is empty or not set")
		return
	}

	// Check for broken symlinks in PATH
	for _, p := range paths {
		if utils.IsSymlink(p) && !utils.PathExists(p) {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Broken symlink in PATH: %s", p))
		}
	}
}

// ScanAsYAML returns the scan result in YAML format
func (d *Detector) ScanAsYAML(result *types.ScanResult) (string, error) {
	// TODO: Implement YAML marshaling
	return fmt.Sprintf("%+v", result), nil
}
