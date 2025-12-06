package cleaner

import (
	"envfix/internal/utils"
	"os"
	"path/filepath"
)

type Cleaner struct{}

func NewCleaner() *Cleaner {
	return &Cleaner{}
}

// Clean removes orphaned and unused environment artifacts
func (c *Cleaner) Clean() error {
	utils.Info("Starting environment cleanup...")

	cleanedSize := int64(0)
	var err error

	// Clean Python
	size, err := c.CleanPython()
	if err != nil {
		utils.Warning("Python cleanup failed: %v", err)
	} else {
		cleanedSize += size
	}

	// Clean Node.js
	size, err = c.CleanNode()
	if err != nil {
		utils.Warning("Node.js cleanup failed: %v", err)
	} else {
		cleanedSize += size
	}

	// Clean broken symlinks in PATH
	err = c.CleanBrokenSymlinks()
	if err != nil {
		utils.Warning("Symlink cleanup failed: %v", err)
	}

	// Clean caches
	size, err = c.CleanCaches()
	if err != nil {
		utils.Warning("Cache cleanup failed: %v", err)
	} else {
		cleanedSize += size
	}

	// Clean orphaned environments
	size, err = c.CleanOrphanedEnvironments()
	if err != nil {
		utils.Warning("Orphaned environment cleanup failed: %v", err)
	} else {
		cleanedSize += size
	}

	utils.Success("✓ Cleanup complete (freed ~%.2f MB)", float64(cleanedSize)/1024/1024)
	return nil
}

// CleanPython removes old venv directories and caches
func (c *Cleaner) CleanPython() (int64, error) {
	var totalSize int64

	utils.Info("Cleaning Python environments...")

	// Clean pip cache
	home := utils.GetHomeDir()
	pipCache := filepath.Join(home, ".cache", "pip")
	if utils.IsWindows() {
		pipCache = filepath.Join(home, "AppData", "Local", "pip", "Cache")
	}

	if utils.PathExists(pipCache) {
		size, _ := utils.DirectorySize(pipCache)
		if err := utils.RemoveDirectory(pipCache); err == nil {
			utils.Info("Cleaned pip cache: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	// Clean uv cache
	uvCache := filepath.Join(home, ".cache", "uv")
	if utils.PathExists(uvCache) {
		size, _ := utils.DirectorySize(uvCache)
		if err := utils.RemoveDirectory(uvCache); err == nil {
			utils.Info("Cleaned uv cache: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	return totalSize, nil
}

// CleanNode removes old node_modules and caches
func (c *Cleaner) CleanNode() (int64, error) {
	var totalSize int64

	utils.Info("Cleaning Node.js environments...")

	home := utils.GetHomeDir()

	// Clean npm cache
	npmCache := filepath.Join(home, ".npm")
	if utils.PathExists(npmCache) {
		size, _ := utils.DirectorySize(npmCache)
		if err := utils.RemoveDirectory(npmCache); err == nil {
			utils.Info("Cleaned npm cache: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	// Clean pnpm cache
	pnpmCache := filepath.Join(home, ".pnpm-store")
	if utils.PathExists(pnpmCache) {
		size, _ := utils.DirectorySize(pnpmCache)
		if err := utils.RemoveDirectory(pnpmCache); err == nil {
			utils.Info("Cleaned pnpm cache: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	// Clean yarn cache
	yarnCache := filepath.Join(home, ".yarn")
	if utils.PathExists(yarnCache) {
		// Don't remove .yarn, just cache
		yarnCacheDir := filepath.Join(yarnCache, "cache")
		if utils.PathExists(yarnCacheDir) {
			size, _ := utils.DirectorySize(yarnCacheDir)
			if err := utils.RemoveDirectory(yarnCacheDir); err == nil {
				utils.Info("Cleaned yarn cache: %.2f MB", float64(size)/1024/1024)
				totalSize += size
			}
		}
	}

	// Clean old node_modules backup
	if utils.PathExists("node_modules.backup") {
		size, _ := utils.DirectorySize("node_modules.backup")
		if err := utils.RemoveDirectory("node_modules.backup"); err == nil {
			utils.Info("Removed node_modules backup: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	return totalSize, nil
}

// CleanBrokenSymlinks removes broken symlinks from PATH
func (c *Cleaner) CleanBrokenSymlinks() error {
	utils.Info("Checking for broken symlinks...")

	paths := utils.SplitPath()
	brokenCount := 0

	for _, p := range paths {
		entries, err := utils.ListDirectory(p)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			fullPath := filepath.Join(p, entry.Name())
			if utils.IsSymlink(fullPath) && !utils.PathExists(fullPath) {
				utils.Info("Found broken symlink: %s", fullPath)
				if err := os.Remove(fullPath); err == nil {
					brokenCount++
				}
			}
		}
	}

	if brokenCount > 0 {
		utils.Success("Removed %d broken symlink(s)", brokenCount)
	}

	return nil
}

// CleanCaches removes various package manager caches
func (c *Cleaner) CleanCaches() (int64, error) {
	var totalSize int64

	utils.Info("Cleaning caches...")

	home := utils.GetHomeDir()

	// Rust cache
	cargoHome := os.Getenv("CARGO_HOME")
	if cargoHome == "" {
		cargoHome = filepath.Join(home, ".cargo")
	}
	registryCache := filepath.Join(cargoHome, "registry", "cache")
	if utils.PathExists(registryCache) {
		size, _ := utils.DirectorySize(registryCache)
		if err := utils.RemoveDirectory(registryCache); err == nil {
			utils.Info("Cleaned cargo cache: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	return totalSize, nil
}

// RemoveOldVenvs removes venv backups
func (c *Cleaner) RemoveOldVenvs() (int64, error) {
	var totalSize int64

	if utils.PathExists(".venv.old") {
		size, _ := utils.DirectorySize(".venv.old")
		if err := utils.RemoveDirectory(".venv.old"); err == nil {
			utils.Info("Removed old venv: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	return totalSize, nil
}

// CleanOrphanedEnvironments removes orphaned venvs and old node_modules
func (c *Cleaner) CleanOrphanedEnvironments() (int64, error) {
	var totalSize int64

	utils.Info("Cleaning orphaned environments...")

	home := utils.GetHomeDir()

	// Clean orphaned venvs in common locations
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
			if c.isOrphanedVenv(fullPath) {
				size, _ := utils.DirectorySize(fullPath)
				if err := utils.RemoveDirectory(fullPath); err == nil {
					utils.Info("Removed orphaned venv: %s (%.2f MB)", entry.Name(), float64(size)/1024/1024)
					totalSize += size
				}
			}
		}
	}

	// Clean old node_modules backups
	if utils.PathExists("node_modules.backup") {
		size, _ := utils.DirectorySize("node_modules.backup")
		if err := utils.RemoveDirectory("node_modules.backup"); err == nil {
			utils.Info("Removed old node_modules backup: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	if utils.PathExists("node_modules.old") {
		size, _ := utils.DirectorySize("node_modules.old")
		if err := utils.RemoveDirectory("node_modules.old"); err == nil {
			utils.Info("Removed old node_modules: %.2f MB", float64(size)/1024/1024)
			totalSize += size
		}
	}

	return totalSize, nil
}

// isOrphanedVenv checks if a venv directory is orphaned
func (c *Cleaner) isOrphanedVenv(venvPath string) bool {
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
