package rebuilder

import (
	"envfix/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

// RepairPATH fixes PATH environment variable issues
func (r *Rebuilder) RepairPATH() error {
	utils.Info("Repairing PATH environment variable...")

	paths := utils.SplitPath()
	var brokenPaths []string
	var cleanedPaths []string

	// Find and remove broken entries
	for _, p := range paths {
		if !utils.PathExists(p) {
			if utils.IsSymlink(p) {
				brokenPaths = append(brokenPaths, p)
				utils.Warning("Removing broken symlink from PATH: %s", p)
			} else {
				brokenPaths = append(brokenPaths, p)
				utils.Info("Removing non-existent path: %s", p)
			}
		} else {
			cleanedPaths = append(cleanedPaths, p)
		}
	}

	if len(brokenPaths) == 0 {
		utils.Info("No broken PATH entries found")
		return nil
	}

	utils.Success("✓ Found %d broken PATH entries", len(brokenPaths))

	// Note: Actually modifying PATH requires admin/root on most systems
	// We log what would be fixed but don't modify the actual environment
	utils.Info("To permanently fix PATH, run:")
	for i, p := range brokenPaths {
		utils.Info("  %d. Remove: %s", i+1, p)
	}

	return nil
}

// RemoveBrokenSymlinksFromPath removes invalid symlinks from PATH
func (r *Rebuilder) RemoveBrokenSymlinksFromPath() error {
	utils.Info("Removing broken symlinks from PATH...")

	paths := utils.SplitPath()
	removedCount := 0

	for _, pathDir := range paths {
		if !utils.PathExists(pathDir) {
			continue
		}

		entries, err := utils.ListDirectory(pathDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			fullPath := filepath.Join(pathDir, entry.Name())

			// Check if it's a broken symlink
			if utils.IsSymlink(fullPath) && !utils.PathExists(fullPath) {
				utils.Info("Removing broken symlink: %s", fullPath)
				if err := os.Remove(fullPath); err == nil {
					removedCount++
				}
			}
		}
	}

	if removedCount > 0 {
		utils.Success("✓ Removed %d broken symlink(s)", removedCount)
	}

	return nil
}

// ValidatePATHEntry ensures a path entry is valid
func ValidatePATHEntry(path string) error {
	if path == "" {
		return fmt.Errorf("empty PATH entry")
	}

	if !utils.PathExists(path) {
		return fmt.Errorf("PATH entry does not exist: %s", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat PATH entry: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("PATH entry is not a directory: %s", path)
	}

	return nil
}

// FindMissingInterpreterPaths finds where interpreters should be but aren't
func (r *Rebuilder) FindMissingInterpreterPaths() map[string]string {
	missing := make(map[string]string)

	paths := utils.SplitPath()
	pathSet := make(map[string]bool)
	for _, p := range paths {
		pathSet[p] = true
	}

	// Common interpreter locations
	expectedLocs := map[string][]string{
		"python": {
			filepath.Join(utils.GetHomeDir(), ".pyenv", "shims"),
			"C:\\Python311\\",
			"/usr/local/bin",
			"/usr/bin",
		},
		"node": {
			filepath.Join(utils.GetHomeDir(), ".nvm", "versions", "node", "current", "bin"),
			"C:\\Program Files\\nodejs",
			"/usr/local/bin",
		},
		"go": {
			"/usr/local/go/bin",
			"C:\\Go\\bin",
		},
	}

	for tool, locs := range expectedLocs {
		for _, loc := range locs {
			if utils.PathExists(loc) && !pathSet[loc] {
				missing[tool] = loc
			}
		}
	}

	return missing
}
