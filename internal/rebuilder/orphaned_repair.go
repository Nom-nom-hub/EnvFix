package rebuilder

import (
	"envfix/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

// RemoveOrphanedEnvironments removes unused venvs and toolchains
func (r *Rebuilder) RemoveOrphanedEnvironments() error {
	utils.Info("Removing orphaned environments...")

	home := utils.GetHomeDir()
	totalRemoved := 0

	// Remove orphaned venvs
	venvLocations := []string{
		filepath.Join(home, ".venv"),
		filepath.Join(home, "venv"),
		filepath.Join(home, ".virtualenvs"),
		filepath.Join(home, ".local", "share", "venvs"),
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
			if r.isOrphanedVenv(fullPath) {
				if err := utils.RemoveDirectory(fullPath); err == nil {
					utils.Success("✓ Removed orphaned venv: %s", entry.Name())
					totalRemoved++
				} else {
					utils.Warning("Failed to remove %s: %v", fullPath, err)
				}
			}
		}
	}

	// Remove old node_modules backups
	backups := []string{
		"node_modules.backup",
		"node_modules.old",
		".node_modules.backup",
	}

	for _, backup := range backups {
		if utils.PathExists(backup) {
			if err := utils.RemoveDirectory(backup); err == nil {
				utils.Success("✓ Removed old backup: %s", backup)
				totalRemoved++
			}
		}
	}

	if totalRemoved == 0 {
		utils.Info("No orphaned environments found")
	} else {
		utils.Success("✓ Removed %d orphaned environment(s)", totalRemoved)
	}

	return nil
}

// isOrphanedVenv checks if a venv directory is orphaned
func (r *Rebuilder) isOrphanedVenv(venvPath string) bool {
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

	return true
}

// ValidateEnvironmentToolchains checks for multiple toolchain versions
func (r *Rebuilder) ValidateEnvironmentToolchains() error {
	utils.Info("Validating environment toolchains...")

	home := utils.GetHomeDir()
	issues := []string{}

	// Check for multiple Go versions
	goVersionDirs := []string{
		filepath.Join(home, ".gimme", "versions"),
		filepath.Join(home, ".go", "versions"),
	}

	for _, dir := range goVersionDirs {
		if utils.PathExists(dir) {
			entries, err := utils.ListDirectory(dir)
			if err == nil && len(entries) > 1 {
				issues = append(issues, fmt.Sprintf("Multiple Go versions in %s", dir))
			}
		}
	}

	// Check for multiple Rust toolchains
	rustupHome := os.Getenv("RUSTUP_HOME")
	if rustupHome == "" {
		rustupHome = filepath.Join(home, ".rustup")
	}

	toolchainsDir := filepath.Join(rustupHome, "toolchains")
	if utils.PathExists(toolchainsDir) {
		entries, err := utils.ListDirectory(toolchainsDir)
		if err == nil && len(entries) > 2 {
			issues = append(issues, fmt.Sprintf("Multiple Rust toolchains (%d found)", len(entries)))
		}
	}

	if len(issues) == 0 {
		utils.Success("✓ All toolchains are valid")
		return nil
	}

	for _, issue := range issues {
		utils.Warning("%s", issue)
	}

	return nil
}
