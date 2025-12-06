package rebuilder

import (
	"envfix/internal/utils"
	"fmt"
	"os"
)

// RepairNodeModules rebuilds node_modules
func (r *Rebuilder) RepairNodeModules() error {
	utils.Info("Repairing Node.js modules...")

	// Detect package manager
	pm := detectNodePackageManager()
	if pm == "" {
		return fmt.Errorf("no package manager detected (npm, yarn, pnpm)")
	}

	// Backup existing node_modules
	if utils.PathExists("node_modules") {
		backupPath := "node_modules.backup"
		utils.Info("Backing up existing node_modules")
		if err := utils.RemoveDirectory(backupPath); err != nil {
			utils.Warning("Could not remove old backup: %v", err)
		}
		if err := os.Rename("node_modules", backupPath); err != nil {
			return fmt.Errorf("failed to backup node_modules: %w", err)
		}
	}

	// Install dependencies
	utils.Info("Running %s install...", pm)
	_, err := utils.RunCommand(pm, "install")
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	utils.Success("Node.js modules repaired")
	return nil
}

// ResolveNodePackageManagerConflict resolves conflicts between multiple package managers
func (r *Rebuilder) ResolveNodePackageManagerConflict() error {
	utils.Info("Resolving Node.js package manager conflicts...")

	hasNpm := utils.PathExists("package-lock.json")
	hasPnpm := utils.PathExists("pnpm-lock.yaml")
	hasYarn := utils.PathExists("yarn.lock")

	count := 0
	if hasNpm {
		count++
	}
	if hasPnpm {
		count++
	}
	if hasYarn {
		count++
	}

	if count <= 1 {
		utils.Info("No package manager conflicts detected")
		return nil
	}

	// Keep the most commonly used one (npm) and remove others
	utils.Info("Removing conflicting lock files...")
	if hasPnpm {
		if err := os.Remove("pnpm-lock.yaml"); err != nil {
			utils.Warning("Could not remove pnpm-lock.yaml: %v", err)
		}
	}
	if hasYarn {
		if err := os.Remove("yarn.lock"); err != nil {
			utils.Warning("Could not remove yarn.lock: %v", err)
		}
	}

	utils.Success("Package manager conflict resolved (using npm)")
	return nil
}

// CleanNodeCache clears Node.js caches
func (r *Rebuilder) CleanNodeCache() error {
	utils.Info("Cleaning Node.js cache...")

	pm := detectNodePackageManager()
	if pm == "" {
		pm = "npm"
	}

	_, err := utils.RunCommand(pm, "cache", "clean", "--force")
	if err != nil {
		utils.Warning("Failed to clean cache: %v", err)
		return nil // Non-fatal
	}

	utils.Success("Node.js cache cleaned")
	return nil
}

func detectNodePackageManager() string {
	if utils.PathExists("pnpm-lock.yaml") {
		if _, err := utils.FindExecutable("pnpm"); err == nil {
			return "pnpm"
		}
	}
	if utils.PathExists("yarn.lock") {
		if _, err := utils.FindExecutable("yarn"); err == nil {
			return "yarn"
		}
	}
	if _, err := utils.FindExecutable("npm"); err == nil {
		return "npm"
	}
	return ""
}
