package detector

import (
	"envfix/internal/types"
	"envfix/internal/utils"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// CheckVersionConflicts detects multiple incompatible versions
func (d *Detector) CheckVersionConflicts(result *types.ScanResult) {
	// Check for Python version mismatches
	d.checkPythonVersionConflicts(result)

	// Check for Node version managers with multiple versions
	d.checkNodeVersionConflicts(result)

	// Check minimum version requirements
	d.checkMinimumVersionRequirements(result)

	// Check for outdated versions
	d.checkOutdatedVersions(result)
}

func (d *Detector) checkPythonVersionConflicts(result *types.ScanResult) {
	if result.Python.Status != "healthy" {
		return
	}

	home := utils.GetHomeDir()
	pythonVersions := []string{}

	// Check pyenv
	pyenvShims := filepath.Join(home, ".pyenv", "shims", "python")
	if utils.PathExists(pyenvShims) {
		pythonVersions = append(pythonVersions, "pyenv")
	}

	// Check conda
	condaPython := filepath.Join(home, "anaconda3", "bin", "python")
	if utils.PathExists(condaPython) {
		pythonVersions = append(pythonVersions, "conda")
	}

	// Check system Python
	if _, err := utils.RunCommand(result.Python.Executable, "--version"); err == nil {
		pythonVersions = append(pythonVersions, "system")
	}

	if len(pythonVersions) > 1 {
		result.Warnings = append(result.Warnings,
			"Multiple Python version managers detected: "+strings.Join(pythonVersions, ", "))
	}
}

func (d *Detector) checkNodeVersionConflicts(result *types.ScanResult) {
	if result.Node.Status != "healthy" {
		return
	}

	home := utils.GetHomeDir()
	versionManagers := []string{}

	// Check nvm
	nvmPath := filepath.Join(home, ".nvm")
	if utils.PathExists(nvmPath) {
		versionManagers = append(versionManagers, "nvm")
	}

	// Check n (node version manager)
	nPath := filepath.Join(home, ".n")
	if utils.PathExists(nPath) {
		versionManagers = append(versionManagers, "n")
	}

	// Check nodeenv
	nodeenvPath := filepath.Join(home, ".nodeenv")
	if utils.PathExists(nodeenvPath) {
		versionManagers = append(versionManagers, "nodeenv")
	}

	if len(versionManagers) > 1 {
		result.Warnings = append(result.Warnings,
			"Multiple Node version managers detected: "+strings.Join(versionManagers, ", "))
	}
}

// ParseVersion extracts version numbers for comparison
func ParseVersion(versionStr string) (major, minor, patch int) {
	// Extract version pattern like "3.11.0" or "v18.0.0"
	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) >= 4 {
		if val, err := strconv.Atoi(matches[1]); err == nil {
			major = val
		}
		if val, err := strconv.Atoi(matches[2]); err == nil {
			minor = val
		}
		if val, err := strconv.Atoi(matches[3]); err == nil {
			patch = val
		}
	}

	return
}

// IsMinimumVersion checks if a version meets minimum requirements
func IsMinimumVersion(versionStr string, minMajor, minMinor int) bool {
	major, minor, _ := ParseVersion(versionStr)

	if major > minMajor {
		return true
	}
	if major == minMajor && minor >= minMinor {
		return true
	}

	return false
}

// DetectVersionMismatches checks for incompatible version combinations
func (d *Detector) DetectVersionMismatches(result *types.ScanResult) []types.Issue {
	issues := []types.Issue{}

	// Python 2 vs 3
	if result.Python.Status == "healthy" {
		if strings.Contains(result.Python.Version, "Python 2") {
			issues = append(issues, types.Issue{
				ID:       "python_v2_detected",
				Title:    "Python 2 detected (EOL)",
				Severity: "warning",
				Language: "python",
				Message:  "Python 2 is end-of-life. Consider upgrading to Python 3.",
				Fix:      "Install Python 3: https://www.python.org/downloads/",
			})
		}
	}

	// Node.js major version compatibility
	if result.Node.Status == "healthy" {
		major, _, _ := ParseVersion(result.Node.Version)
		if major < 14 {
			issues = append(issues, types.Issue{
				ID:       "node_version_old",
				Title:    "Node.js version is outdated",
				Severity: "warning",
				Language: "nodejs",
				Message:  "Node.js versions below 14 are not supported. Consider upgrading.",
				Fix:      "Install Node.js 18+ from https://nodejs.org/",
			})
		}
	}

	return issues
}

// CheckPythonMinimumVersion validates Python meets minimum version
func CheckPythonMinimumVersion(versionStr string) bool {
	return IsMinimumVersion(versionStr, 3, 7)
}

// CheckNodeMinimumVersion validates Node.js meets minimum version
func CheckNodeMinimumVersion(versionStr string) bool {
	return IsMinimumVersion(versionStr, 14, 0)
}

// FindAlternativeInterpreterVersions finds other installed interpreter versions
func FindAlternativeInterpreterVersions(interpreter string) []string {
	var versions []string
	home := utils.GetHomeDir()

	if interpreter == "python" {
		// Check common Python locations
		pythonPaths := map[string]bool{
			"/usr/bin/python3.11":                     true,
			"/usr/bin/python3.10":                     true,
			"/usr/bin/python3.9":                      true,
			"C:\\Python311\\python.exe":               true,
			"C:\\Python310\\python.exe":               true,
			"C:\\Python39\\python.exe":                true,
			filepath.Join(home, ".pyenv", "versions"): true,
		}

		for path := range pythonPaths {
			if utils.PathExists(path) {
				versions = append(versions, path)
			}
		}
	}

	return versions
}

// checkMinimumVersionRequirements validates that all interpreters meet minimum versions
func (d *Detector) checkMinimumVersionRequirements(result *types.ScanResult) {
	// Python minimum: 3.7
	if result.Python.Status == "healthy" && result.Python.Version != "" {
		if !CheckPythonMinimumVersion(result.Python.Version) {
			result.Issues = append(result.Issues, types.Issue{
				ID:       "python_version_low",
				Title:    "Python version below minimum requirement",
				Severity: "warning",
				Language: "python",
				Message:  "Python 3.7+ is required. Current: " + result.Python.Version,
				Fix:      "Install Python 3.8 or later from https://www.python.org/downloads/",
			})
		}
	}

	// Node.js minimum: 14
	if result.Node.Status == "healthy" && result.Node.Version != "" {
		if !CheckNodeMinimumVersion(result.Node.Version) {
			result.Issues = append(result.Issues, types.Issue{
				ID:       "node_version_low",
				Title:    "Node.js version below minimum requirement",
				Severity: "warning",
				Language: "nodejs",
				Message:  "Node.js 14+ is required. Current: " + result.Node.Version,
				Fix:      "Install Node.js 16+ from https://nodejs.org/",
			})
		}
	}
}

// checkOutdatedVersions detects versions that are outdated or deprecated
func (d *Detector) checkOutdatedVersions(result *types.ScanResult) {
	// Check Python 2 (EOL)
	if result.Python.Status == "healthy" && result.Python.Version != "" {
		if strings.Contains(result.Python.Version, "Python 2") {
			result.Issues = append(result.Issues, types.Issue{
				ID:       "python_eol",
				Title:    "Python 2 is end-of-life",
				Severity: "critical",
				Language: "python",
				Message:  "Python 2 reached end-of-life on January 1, 2020. Upgrade to Python 3 immediately.",
				Fix:      "Install Python 3.8+ from https://www.python.org/downloads/",
			})
		}
	}

	// Check Node.js 12 (EOL)
	if result.Node.Status == "healthy" && result.Node.Version != "" {
		major, _, _ := ParseVersion(result.Node.Version)
		if major < 14 {
			if major == 12 {
				result.Issues = append(result.Issues, types.Issue{
					ID:       "node_eol",
					Title:    "Node.js 12 is end-of-life",
					Severity: "critical",
					Language: "nodejs",
					Message:  "Node.js 12 reached end-of-life on April 30, 2022. Upgrade to 16+ immediately.",
					Fix:      "Install Node.js 18+ from https://nodejs.org/",
				})
			} else {
				result.Issues = append(result.Issues, types.Issue{
					ID:       "node_version_old",
					Title:    "Node.js version is very outdated",
					Severity: "warning",
					Language: "nodejs",
					Message:  "Node.js " + result.Node.Version + " is outdated. Recommend 18+.",
					Fix:      "Install Node.js 18+ from https://nodejs.org/",
				})
			}
		}
	}

	// Check Rust (1.56 is reasonable minimum)
	if result.Rust.Status == "healthy" && result.Rust.Version != "" {
		if !CheckRustMinimumVersion(result.Rust.Version) {
			result.Warnings = append(result.Warnings,
				"Rust version is outdated: "+result.Rust.Version+" (recommend 1.56+)")
		}
	}

	// Check Go (1.16 is reasonable minimum)
	if result.Go.Status == "healthy" && result.Go.Version != "" {
		if !CheckGoMinimumVersion(result.Go.Version) {
			result.Warnings = append(result.Warnings,
				"Go version is outdated: "+result.Go.Version+" (recommend 1.16+)")
		}
	}
}

// CheckRustMinimumVersion validates Rust meets minimum version
func CheckRustMinimumVersion(versionStr string) bool {
	return IsMinimumVersion(versionStr, 1, 56)
}

// CheckGoMinimumVersion validates Go meets minimum version
func CheckGoMinimumVersion(versionStr string) bool {
	return IsMinimumVersion(versionStr, 1, 16)
}
