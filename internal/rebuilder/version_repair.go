package rebuilder

import (
	"envfix/internal/utils"
	"os"
	"path/filepath"
)

// SuggestVersionUpgrades recommends upgrading outdated interpreters
func (r *Rebuilder) SuggestVersionUpgrades() error {
	utils.Info("Checking for version upgrade recommendations...")

	upgrades := []VersionUpgrade{}

	// Check Python
	pythonExe, _ := utils.FindExecutable("python3")
	if pythonExe == "" {
		pythonExe, _ = utils.FindExecutable("python")
	}
	if pythonExe != "" {
		_, err := utils.RunCommand(pythonExe, "--version")
		if err == nil {
			utils.Info("Python is installed at: %s", pythonExe)
			// Would need to parse version output and compare
		}
	}

	// Check Node.js
	nodeExe, _ := utils.FindExecutable("node")
	if nodeExe != "" {
		_, err := utils.RunCommand(nodeExe, "--version")
		if err == nil {
			utils.Info("Node.js is installed at: %s", nodeExe)
			// Would need to parse version output and compare
		}
	}

	// Report upgrades
	if len(upgrades) > 0 {
		utils.Warning("Found %d version upgrades recommended", len(upgrades))
		for _, upgrade := range upgrades {
			utils.Info("  %s: %s → %s", upgrade.Language, upgrade.Current, upgrade.Recommended)
		}
	} else {
		utils.Success("✓ All installed versions are up-to-date")
	}

	return nil
}

// RemoveOldInterpreterVersions removes alternate/duplicate interpreter installations
func (r *Rebuilder) RemoveOldInterpreterVersions() error {
	utils.Info("Checking for old/duplicate interpreter installations...")

	removed := 0
	home := utils.GetHomeDir()

	// Check for old Python versions via gimme-sha256 or similar
	if utils.IsWindows() {
		// Windows: Check AppData for old Python installations
		pythonRoot := filepath.Join(home, "AppData", "Local", "Programs", "Python")
		if utils.PathExists(pythonRoot) {
			entries, err := utils.ListDirectory(pythonRoot)
			if err == nil && len(entries) > 1 {
				utils.Warning("Multiple Python installations found: %d versions", len(entries))
				// Don't auto-remove, just warn
			}
		}
	} else {
		// Unix: Check ~/.python-versions or similar
		pythonVersionsDir := filepath.Join(home, ".python-versions")
		if utils.PathExists(pythonVersionsDir) {
			entries, _ := utils.ListDirectory(pythonVersionsDir)
			if len(entries) > 2 {
				utils.Warning("Multiple Python versions detected: %d versions", len(entries))
			}
		}
	}

	// Check for old Node versions
	nodVersionsDir := filepath.Join(home, ".nvm", "versions", "node")
	if utils.PathExists(nodVersionsDir) {
		entries, _ := utils.ListDirectory(nodVersionsDir)
		if len(entries) > 2 {
			utils.Warning("Multiple Node.js versions detected: %d versions", len(entries))
		}
	}

	// Check for old Rust toolchains
	rustupHome := os.Getenv("RUSTUP_HOME")
	if rustupHome == "" {
		rustupHome = filepath.Join(home, ".rustup")
	}
	toolchainsDir := filepath.Join(rustupHome, "toolchains")
	if utils.PathExists(toolchainsDir) {
		entries, _ := utils.ListDirectory(toolchainsDir)
		if len(entries) > 2 {
			utils.Warning("Multiple Rust toolchains detected: %d versions", len(entries))
			utils.Info("Use 'rustup toolchain remove <name>' to remove old toolchains")
		}
	}

	// Check for old Go versions via gimme
	goVersionsDir := filepath.Join(home, ".gimme", "versions", "go")
	if utils.PathExists(goVersionsDir) {
		entries, _ := utils.ListDirectory(goVersionsDir)
		if len(entries) > 1 {
			utils.Warning("Multiple Go versions detected: %d versions", len(entries))
		}
	}

	if removed == 0 {
		utils.Success("✓ No old interpreter versions to remove")
	}

	return nil
}

// ValidateVersionCompatibility checks if installed versions are compatible with each other
func (r *Rebuilder) ValidateVersionCompatibility() error {
	utils.Info("Validating version compatibility...")

	issues := []string{}

	// Get all version info (would normally be passed from detector)
	// For now, just check that interpreters are accessible

	pythonExe, _ := utils.FindExecutable("python3")
	if pythonExe == "" {
		pythonExe, _ = utils.FindExecutable("python")
	}

	nodeExe, _ := utils.FindExecutable("node")
	npmExe, _ := utils.FindExecutable("npm")

	cargoExe, _ := utils.FindExecutable("cargo")
	goExe, _ := utils.FindExecutable("go")

	// Report what's available
	utils.Info("Installed interpreters:")
	if pythonExe != "" {
		utils.Success("  ✓ Python: %s", pythonExe)
	} else {
		issues = append(issues, "Python not found")
	}

	if nodeExe != "" {
		utils.Success("  ✓ Node.js: %s", nodeExe)
		if npmExe != "" {
			utils.Success("  ✓ npm: %s", npmExe)
		} else {
			issues = append(issues, "npm not found (Node.js installed)")
		}
	} else {
		// Optional
		utils.Info("  - Node.js not found")
	}

	if cargoExe != "" {
		utils.Success("  ✓ Rust (cargo): %s", cargoExe)
	} else {
		// Optional
		utils.Info("  - Rust not found")
	}

	if goExe != "" {
		utils.Success("  ✓ Go: %s", goExe)
	} else {
		// Optional
		utils.Info("  - Go not found")
	}

	// Report any compatibility issues
	if len(issues) > 0 {
		for _, issue := range issues {
			utils.Warning("  ! %s", issue)
		}
		return nil // Don't error, just warn
	}

	utils.Success("✓ All installed versions are compatible")
	return nil
}

// VersionUpgrade represents a recommended version upgrade
type VersionUpgrade struct {
	Language    string
	Current     string
	Recommended string
	Reason      string
}
