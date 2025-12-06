package rebuilder

import (
	"envfix/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// InstallPythonDependenciesUpgrade upgrades all Python dependencies
func (r *Rebuilder) InstallPythonDependenciesUpgrade() error {
	if !utils.PathExists("requirements.txt") {
		utils.Info("No requirements.txt found, skipping Python dependencies")
		return nil
	}

	utils.Info("Upgrading Python dependencies...")

	// Find pip
	pipExe := r.findPipExecutable()
	if pipExe == "" {
		return fmt.Errorf("pip not found in PATH or venv")
	}

	// Upgrade all packages
	_, err := utils.RunCommand(pipExe, "install", "--upgrade", "-r", "requirements.txt")
	if err != nil {
		utils.Warning("Failed to upgrade all packages, attempting basic install: %v", err)
		// Fallback to basic install
		_, err = utils.RunCommand(pipExe, "install", "-r", "requirements.txt")
		if err != nil {
			return fmt.Errorf("failed to install Python dependencies: %w", err)
		}
	}

	utils.Success("Python dependencies upgraded")
	return nil
}

// InstallNodeDependencies installs Node.js dependencies with npm ci (clean install)
func (r *Rebuilder) InstallNodeDependencies() error {
	if !utils.PathExists("package.json") {
		utils.Info("No package.json found, skipping Node.js dependencies")
		return nil
	}

	utils.Info("Installing Node.js dependencies with clean install (npm ci)...")

	// Detect package manager
	pm := detectNodePackageManager()
	if pm == "" {
		return fmt.Errorf("no package manager detected (npm, yarn, pnpm)")
	}

	// Use npm ci for clean, reproducible install
	if pm == "npm" {
		// Check if package-lock.json exists for npm ci
		if utils.PathExists("package-lock.json") {
			_, err := utils.RunCommand("npm", "ci")
			if err != nil {
				utils.Warning("npm ci failed, falling back to npm install: %v", err)
				_, err = utils.RunCommand("npm", "install")
				if err != nil {
					return fmt.Errorf("failed to install Node.js dependencies: %w", err)
				}
			}
		} else {
			// Use regular install if no lock file
			_, err := utils.RunCommand("npm", "install")
			if err != nil {
				return fmt.Errorf("failed to install Node.js dependencies: %w", err)
			}
		}
	} else if pm == "pnpm" {
		_, err := utils.RunCommand("pnpm", "install")
		if err != nil {
			return fmt.Errorf("failed to install with pnpm: %w", err)
		}
	} else if pm == "yarn" {
		_, err := utils.RunCommand("yarn", "install")
		if err != nil {
			return fmt.Errorf("failed to install with yarn: %w", err)
		}
	}

	utils.Success("Node.js dependencies installed")
	return nil
}

// InstallRustDependencies updates Rust dependencies via Cargo
func (r *Rebuilder) InstallRustDependencies() error {
	if !utils.PathExists("Cargo.toml") {
		utils.Info("No Cargo.toml found, skipping Rust dependencies")
		return nil
	}

	utils.Info("Updating Rust dependencies...")

	// Update Cargo.lock to latest compatible versions
	_, err := utils.RunCommand("cargo", "update")
	if err != nil {
		return fmt.Errorf("failed to update Rust dependencies: %w", err)
	}

	utils.Success("Rust dependencies updated")
	return nil
}

// TidyGoDependencies removes unused Go dependencies
func (r *Rebuilder) TidyGoDependencies() error {
	if !utils.PathExists("go.mod") {
		utils.Info("No go.mod found, skipping Go module tidying")
		return nil
	}

	utils.Info("Tidying Go module dependencies...")

	// Tidy up go.mod and go.sum
	_, err := utils.RunCommand("go", "mod", "tidy")
	if err != nil {
		return fmt.Errorf("failed to tidy Go modules: %w", err)
	}

	utils.Success("Go modules tidied")
	return nil
}

// VerifyDependencies checks if all dependencies are properly installed
func (r *Rebuilder) VerifyDependencies() error {
	utils.Info("Verifying dependencies...")

	verified := 0

	// Check Python
	if utils.PathExists("requirements.txt") {
		pipExe := r.findPipExecutable()
		if pipExe != "" {
			// Run pip check to verify installed packages
			_, err := utils.RunCommand(pipExe, "check")
			if err == nil {
				utils.Success("✓ Python dependencies verified")
				verified++
			} else {
				utils.Warning("Python dependency verification failed: %v", err)
			}
		}
	}

	// Check Node.js
	if utils.PathExists("package.json") {
		pm := detectNodePackageManager()
		if pm == "npm" {
			_, err := utils.RunCommand("npm", "list", "--depth=0")
			if err == nil {
				utils.Success("✓ Node.js dependencies verified")
				verified++
			} else {
				utils.Warning("Node.js dependency verification had issues: %v", err)
			}
		}
	}

	// Check Rust
	if utils.PathExists("Cargo.toml") {
		_, err := utils.RunCommand("cargo", "check")
		if err == nil {
			utils.Success("✓ Rust dependencies verified")
			verified++
		} else {
			utils.Warning("Rust dependency verification failed: %v", err)
		}
	}

	// Check Go
	if utils.PathExists("go.mod") {
		_, err := utils.RunCommand("go", "mod", "verify")
		if err == nil {
			utils.Success("✓ Go dependencies verified")
			verified++
		} else {
			utils.Warning("Go dependency verification failed: %v", err)
		}
	}

	if verified == 0 {
		utils.Info("No dependencies to verify")
	}

	return nil
}

// RegenerateLockfiles regenerates lock files for reproducible builds
func (r *Rebuilder) RegenerateLockfiles() error {
	utils.Info("Regenerating lock files...")

	// Python (create requirements.txt from virtual environment)
	if utils.PathExists(".venv") {
		pipExe := r.findPipExecutable()
		if pipExe != "" {
			utils.Info("Regenerating requirements.txt...")
			output, err := utils.RunCommand(pipExe, "freeze")
			if err == nil && output != "" {
				err = os.WriteFile("requirements.txt", []byte(output), 0644)
				if err == nil {
					utils.Success("✓ requirements.txt regenerated")
				}
			}
		}
	}

	// Node.js (npm regenerates package-lock.json automatically)
	if utils.PathExists("package.json") && utils.PathExists("node_modules") {
		pm := detectNodePackageManager()
		if pm == "npm" {
			utils.Info("Regenerating package-lock.json...")
			_, err := utils.RunCommand("npm", "install")
			if err == nil {
				utils.Success("✓ package-lock.json regenerated")
			}
		}
	}

	// Rust (cargo.lock is regenerated on build)
	if utils.PathExists("Cargo.toml") {
		utils.Info("Regenerating Cargo.lock...")
		_, err := utils.RunCommand("cargo", "update")
		if err == nil {
			utils.Success("✓ Cargo.lock regenerated")
		}
	}

	// Go (go.sum is regenerated on go mod tidy)
	if utils.PathExists("go.mod") {
		utils.Info("Regenerating go.sum...")
		_, err := utils.RunCommand("go", "mod", "tidy")
		if err == nil {
			utils.Success("✓ go.sum regenerated")
		}
	}

	return nil
}

// CleanDependencyCaches clears all dependency caches
func (r *Rebuilder) CleanDependencyCaches() error {
	utils.Info("Cleaning dependency caches...")

	home := utils.GetHomeDir()

	// Python pip cache
	pipCachePath := filepath.Join(home, ".cache", "pip")
	if utils.IsWindows() {
		pipCachePath = filepath.Join(home, "AppData", "Local", "pip", "Cache")
	}

	if utils.PathExists(pipCachePath) {
		_, _ = utils.RunCommand("pip", "cache", "purge")
		utils.Success("✓ Python pip cache cleared")
	}

	// Node.js npm cache
	_, _ = utils.RunCommand("npm", "cache", "clean", "--force")
	utils.Success("✓ Node.js npm cache cleared")

	// Cargo cache (partial)
	cargoHome := os.Getenv("CARGO_HOME")
	if cargoHome == "" {
		cargoHome = filepath.Join(home, ".cargo")
	}
	cargoCache := filepath.Join(cargoHome, "registry", "cache")
	if utils.PathExists(cargoCache) {
		utils.RemoveDirectory(cargoCache)
		utils.Success("✓ Rust cargo cache cleared")
	}

	// Go module cache
	_, _ = utils.RunCommand("go", "clean", "-modcache")
	utils.Success("✓ Go module cache cleared")

	return nil
}

// findPipExecutable finds the pip executable in venv or system PATH
func (r *Rebuilder) findPipExecutable() string {
	venvPath := ".venv"

	// Check venv pip first
	if utils.IsWindows() {
		pipPath := filepath.Join(venvPath, "Scripts", "pip.exe")
		if utils.PathExists(pipPath) {
			return pipPath
		}
	} else {
		pipPath := filepath.Join(venvPath, "bin", "pip")
		if utils.PathExists(pipPath) {
			return pipPath
		}
	}

	// Fall back to system pip
	pipExe, err := utils.FindExecutable("pip")
	if err == nil {
		return pipExe
	}

	pipExe, err = utils.FindExecutable("pip3")
	if err == nil {
		return pipExe
	}

	return ""
}

// CheckDependencyLockfiles validates lock file consistency
func (r *Rebuilder) CheckDependencyLockfiles() error {
	utils.Info("Checking dependency lock files...")

	issues := []string{}

	// Python: Check requirements.txt exists
	if !utils.PathExists("requirements.txt") && (utils.PathExists("setup.py") || utils.PathExists("pyproject.toml")) {
		issues = append(issues, "Python project missing requirements.txt")
	}

	// Node.js: Check lock file consistency
	if utils.PathExists("package.json") {
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

		if count > 1 {
			issues = append(issues, "Multiple Node.js lock files detected (should use one)")
		} else if count == 0 {
			issues = append(issues, "Node.js package.json exists but no lock file found")
		}
	}

	// Rust: Check Cargo.lock exists for binary projects
	if utils.PathExists("Cargo.toml") {
		content, err := os.ReadFile("Cargo.toml")
		if err == nil {
			if strings.Contains(string(content), "[package]") && !utils.PathExists("Cargo.lock") {
				// Binary project should have Cargo.lock
				if strings.Contains(string(content), "[[bin]]") {
					issues = append(issues, "Rust binary project missing Cargo.lock")
				}
			}
		}
	}

	// Go: Check go.sum consistency
	if utils.PathExists("go.mod") {
		if !utils.PathExists("go.sum") {
			issues = append(issues, "Go project has go.mod but missing go.sum")
		}
	}

	if len(issues) == 0 {
		utils.Success("✓ All dependency lock files are consistent")
		return nil
	}

	for _, issue := range issues {
		utils.Warning("%s", issue)
	}

	return nil
}
