package detector

import (
	"envfix/internal/types"
	"envfix/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CheckOrphanedEnvironments detects unused venvs, node_modules, and toolchains
func (d *Detector) CheckOrphanedEnvironments(result *types.ScanResult) {
	// Check for orphaned Python venvs
	d.checkOrphanedPythonVenvs(result)

	// Check for orphaned Node.js modules
	d.checkOrphanedNodeModules(result)

	// Check for orphaned toolchains
	d.checkOrphanedToolchains(result)
}

// checkOrphanedPythonVenvs scans for stale or disconnected virtual environments
func (d *Detector) checkOrphanedPythonVenvs(result *types.ScanResult) {
	home := utils.GetHomeDir()

	// Common venv locations
	venvLocations := []string{
		filepath.Join(home, ".venv"),
		filepath.Join(home, "venv"),
		filepath.Join(home, ".virtualenvs"),
		filepath.Join(home, "virtualenvs"),
		filepath.Join(home, ".local", "share", "venvs"),
		filepath.Join(home, ".pyenv", "versions"),
	}

	for _, venvPath := range venvLocations {
		if !utils.PathExists(venvPath) {
			continue
		}

		entries, err := utils.ListDirectory(venvPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			fullPath := filepath.Join(venvPath, entry.Name())
			if d.isOrphanedVenv(fullPath) {
				age := d.getDirectoryAge(fullPath)
				result.Issues = append(result.Issues, types.Issue{
					ID:       "orphaned_venv",
					Title:    "Orphaned Python virtual environment detected",
					Severity: "info",
					Language: "python",
					Message:  "Found orphaned venv: " + fullPath + " (unused for " + age + ")",
					Fix:      "Run 'envfix clean' to remove orphaned environments",
				})
				result.Warnings = append(result.Warnings,
					"Orphaned venv: "+fullPath+" ("+age+" old)")
			}
		}
	}
}

// checkOrphanedNodeModules scans for old node_modules directories
func (d *Detector) checkOrphanedNodeModules(result *types.ScanResult) {
	home := utils.GetHomeDir()

	// Check for orphaned node_modules in common locations
	locations := []string{
		home,
		filepath.Join(home, "Projects"),
		filepath.Join(home, "work"),
		filepath.Join(home, "Developer"),
	}

	for _, location := range locations {
		if !utils.PathExists(location) {
			continue
		}

		d.walkForOrphanedModules(location, result)
	}
}

// walkForOrphanedModules recursively searches for orphaned node_modules
func (d *Detector) walkForOrphanedModules(dir string, result *types.ScanResult) {
	entries, err := utils.ListDirectory(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())

		// Skip certain directories
		if strings.HasPrefix(entry.Name(), ".") ||
			entry.Name() == "node_modules" {
			continue
		}

		// Check if this directory contains an orphaned node_modules
		nodeModulesPath := filepath.Join(fullPath, "node_modules")
		if utils.PathExists(nodeModulesPath) {
			// Check if parent has package.json
			packageJsonPath := filepath.Join(fullPath, "package.json")
			if !utils.PathExists(packageJsonPath) {
				age := d.getDirectoryAge(nodeModulesPath)
				result.Warnings = append(result.Warnings,
					"Orphaned node_modules: "+nodeModulesPath+" ("+age+" old)")
			}
		}
	}
}

// checkOrphanedToolchains detects unused language toolchains
func (d *Detector) checkOrphanedToolchains(result *types.ScanResult) {
	home := utils.GetHomeDir()

	// Check for multiple Go versions via gimme or similar
	goVersionDirs := []string{
		filepath.Join(home, ".gimme", "versions"),
		filepath.Join(home, ".go", "versions"),
	}

	for _, dir := range goVersionDirs {
		if !utils.PathExists(dir) {
			continue
		}

		entries, err := utils.ListDirectory(dir)
		if err != nil {
			continue
		}

		if len(entries) > 1 {
			age := d.getDirectoryAge(dir)
			result.Warnings = append(result.Warnings,
				"Multiple Go versions detected in "+dir+" (oldest: "+age+")")
		}
	}

	// Check for Rust toolchains
	rustupHome := os.Getenv("RUSTUP_HOME")
	if rustupHome == "" {
		rustupHome = filepath.Join(home, ".rustup")
	}

	toolchainsDir := filepath.Join(rustupHome, "toolchains")
	if utils.PathExists(toolchainsDir) {
		entries, err := utils.ListDirectory(toolchainsDir)
		if err == nil && len(entries) > 2 {
			age := d.getDirectoryAge(toolchainsDir)
			result.Warnings = append(result.Warnings,
				"Multiple Rust toolchains detected (oldest: "+age+")")
		}
	}
}

// isOrphanedVenv checks if a venv is orphaned (no parent project)
func (d *Detector) isOrphanedVenv(venvPath string) bool {
	// Check if it's actually a venv
	if !utils.PathExists(filepath.Join(venvPath, "pyvenv.cfg")) {
		return false
	}

	// Get parent directory
	parentDir := filepath.Dir(venvPath)

	// Check if parent has Python project files
	projectFiles := []string{
		filepath.Join(parentDir, "requirements.txt"),
		filepath.Join(parentDir, "pyproject.toml"),
		filepath.Join(parentDir, "setup.py"),
		filepath.Join(parentDir, "poetry.lock"),
		filepath.Join(parentDir, "Pipfile"),
	}

	for _, file := range projectFiles {
		if utils.PathExists(file) {
			return false // Has project files, not orphaned
		}
	}

	// If venv is very new (less than 24 hours), consider it non-orphaned
	// to avoid false positives for recently created venvs
	age := d.getDirectoryAge(venvPath)
	return !strings.Contains(age, "hours") && !strings.Contains(age, "days")
}

// getDirectoryAge returns human-readable age of a directory
func (d *Detector) getDirectoryAge(dirPath string) string {
	info, err := os.Stat(dirPath)
	if err != nil {
		return "unknown"
	}

	age := time.Since(info.ModTime())

	if age < time.Hour {
		return "less than 1 hour"
	} else if age < 24*time.Hour {
		return "less than 1 day"
	} else if age < 7*24*time.Hour {
		days := int(age.Hours() / 24)
		return fmt.Sprintf("%d days", days)
	} else if age < 30*24*time.Hour {
		weeks := int(age.Hours() / 24 / 7)
		return fmt.Sprintf("%d weeks", weeks)
	} else if age < 365*24*time.Hour {
		months := int(age.Hours() / 24 / 30)
		return fmt.Sprintf("%d months", months)
	}

	years := int(age.Hours() / 24 / 365)
	return fmt.Sprintf("%d years", years)
}
