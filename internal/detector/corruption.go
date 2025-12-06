package detector

import (
	"envfix/internal/utils"
	"path/filepath"
	"strings"
)

// CheckVenvCorruption detects if venv is broken or incomplete
func (d *Detector) CheckVenvCorruption(venvPath string) []string {
	issues := []string{}

	if !utils.PathExists(venvPath) {
		return issues
	}

	// Check for required venv files
	requiredPaths := []string{
		filepath.Join(venvPath, "pyvenv.cfg"),
	}

	if utils.IsWindows() {
		requiredPaths = append(requiredPaths,
			filepath.Join(venvPath, "Scripts", "activate.bat"),
			filepath.Join(venvPath, "Scripts", "python.exe"),
		)
	} else {
		requiredPaths = append(requiredPaths,
			filepath.Join(venvPath, "bin", "activate"),
			filepath.Join(venvPath, "bin", "python"),
		)
	}

	missingFiles := []string{}
	for _, p := range requiredPaths {
		if !utils.PathExists(p) {
			missingFiles = append(missingFiles, filepath.Base(p))
		}
	}

	if len(missingFiles) > 0 {
		issues = append(issues, "Virtual environment missing files: "+strings.Join(missingFiles, ", "))
	}

	return issues
}

// CheckNodeModulesCorruption detects broken node_modules
func (d *Detector) CheckNodeModulesCorruption() []string {
	issues := []string{}

	if !utils.PathExists("node_modules") {
		return issues
	}

	// Check if node_modules directory is readable
	entries, err := utils.ListDirectory("node_modules")
	if err != nil {
		issues = append(issues, "node_modules directory exists but is unreadable")
		return issues
	}

	// Check for empty node_modules
	if len(entries) == 0 {
		issues = append(issues, "node_modules directory is empty")
	}

	// Check for typical corruption patterns
	hasBrokenDeps := false

	for _, entry := range entries {
		if entry.IsDir() {
			// Check for broken package directories (no package.json)
			pkgJsonPath := filepath.Join("node_modules", entry.Name(), "package.json")
			if !utils.PathExists(pkgJsonPath) && !strings.HasPrefix(entry.Name(), ".") {
				hasBrokenDeps = true
				break
			}
		}
	}

	if hasBrokenDeps {
		issues = append(issues, "node_modules contains broken or incomplete packages")
	}

	return issues
}

// CheckGlobalPackagePollution detects globally installed packages that should be local
func (d *Detector) CheckGlobalPackagePollution() []string {
	issues := []string{}

	home := utils.GetHomeDir()

	// Check npm global packages
	if utils.IsWindows() {
		npmGlobal := filepath.Join(home, "AppData", "Roaming", "npm", "node_modules")
		if utils.PathExists(npmGlobal) {
			entries, _ := utils.ListDirectory(npmGlobal)
			if len(entries) > 5 {
				issues = append(issues, "High number of globally installed npm packages may cause conflicts")
			}
		}
	} else {
		npmGlobal := filepath.Join(home, ".npm-global", "lib", "node_modules")
		if utils.PathExists(npmGlobal) {
			entries, _ := utils.ListDirectory(npmGlobal)
			if len(entries) > 5 {
				issues = append(issues, "High number of globally installed npm packages may cause conflicts")
			}
		}
	}

	// Check pip global packages
	pipGlobal := filepath.Join(home, ".local", "lib")
	if utils.PathExists(pipGlobal) {
		entries, _ := utils.ListDirectory(pipGlobal)
		if len(entries) > 10 {
			issues = append(issues, "High number of globally installed pip packages detected")
		}
	}

	return issues
}

// CheckOrphanedVenvs detects venv directories that are no longer connected to projects
func (d *Detector) CheckOrphanedVenvs(venvName string) bool {
	if !utils.PathExists(venvName) {
		return false
	}

	// Venv is considered orphaned if:
	// 1. It exists but parent project has no requirements.txt or pyproject.toml
	// 2. It's stale (older than project files)

	hasProjectFiles := utils.PathExists("requirements.txt") ||
		utils.PathExists("pyproject.toml") ||
		utils.PathExists("setup.py")

	return !hasProjectFiles
}

// CheckCacheCorruption detects corrupted cache directories
func (d *Detector) CheckCacheCorruption() []string {
	issues := []string{}

	home := utils.GetHomeDir()

	// Check pip cache
	pipCachePath := filepath.Join(home, ".cache", "pip")
	if utils.IsWindows() {
		pipCachePath = filepath.Join(home, "AppData", "Local", "pip", "Cache")
	}

	if utils.PathExists(pipCachePath) {
		// Try to list cache - if fails, it's corrupted
		if entries, err := utils.ListDirectory(pipCachePath); err != nil || len(entries) == 0 {
			issues = append(issues, "pip cache appears corrupted or empty")
		}
	}

	// Check npm cache
	npmCachePath := filepath.Join(home, ".npm")
	if utils.PathExists(npmCachePath) {
		if entries, err := utils.ListDirectory(npmCachePath); err != nil || len(entries) == 0 {
			issues = append(issues, "npm cache appears corrupted or empty")
		}
	}

	return issues
}
