package detector

import (
	"envfix/internal/parser"
	"envfix/internal/types"
	"envfix/internal/utils"
)

// DetectDependencyIssues detects problems with dependencies
func (d *Detector) DetectDependencyIssues(result *types.ScanResult) {
	// Check for outdated lockfiles
	d.CheckOutdatedLockfiles(result)

	// Check for missing lockfiles
	d.CheckMissingLockfiles(result)

	// Check for dependency conflicts
	d.CheckDependencyConflicts(result)
}

// CheckOutdatedLockfiles detects old or missing lockfiles
func (d *Detector) CheckOutdatedLockfiles(result *types.ScanResult) {
	// Node.js: Check for lock file inconsistency with package.json
	if utils.PathExists("package.json") {
		hasLock := utils.PathExists("package-lock.json") ||
			utils.PathExists("yarn.lock") ||
			utils.PathExists("pnpm-lock.yaml")

		if !hasLock {
			result.Warnings = append(result.Warnings, "package.json found but no lock file (npm, yarn, or pnpm)")
		}
	}

	// Go: Check for missing go.sum
	if utils.PathExists("go.mod") && !utils.PathExists("go.sum") {
		result.Warnings = append(result.Warnings, "go.mod found but go.sum is missing")
	}
}

// CheckMissingLockfiles detects lock files without package files
func (d *Detector) CheckMissingLockfiles(result *types.ScanResult) {
	// package-lock.json without package.json
	if utils.PathExists("package-lock.json") && !utils.PathExists("package.json") {
		result.Warnings = append(result.Warnings, "package-lock.json found but package.json is missing")
	}

	// go.sum without go.mod
	if utils.PathExists("go.sum") && !utils.PathExists("go.mod") {
		result.Warnings = append(result.Warnings, "go.sum found but go.mod is missing")
	}
}

// CheckDependencyConflicts detects potential dependency problems
func (d *Detector) CheckDependencyConflicts(result *types.ScanResult) {
	deps := parser.GetAllDependencies()

	// Count dependencies by manager
	depsByManager := make(map[string]int)
	for _, dep := range deps {
		depsByManager[dep.Manager]++
	}

	// Add info to manifest
	if result.Python.Status == "healthy" && depsByManager["pip"] > 0 {
		result.Python.Packages = make([]string, 0)
		for _, dep := range deps {
			if dep.Manager == "pip" {
				result.Python.Packages = append(result.Python.Packages, dep.Name)
			}
		}
	}
}

// GetDependencySummary returns a summary of detected dependencies
func (d *Detector) GetDependencySummary() map[string]int {
	summary := make(map[string]int)
	deps := parser.GetAllDependencies()

	for _, dep := range deps {
		summary[dep.Manager]++
	}

	return summary
}
