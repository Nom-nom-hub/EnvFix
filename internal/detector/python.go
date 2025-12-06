package detector

import (
	"envfix/internal/types"
	"envfix/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

// DetectPython detects Python environment
func (d *Detector) DetectPython() types.PythonEnv {
	env := types.PythonEnv{
		Status: "missing",
		Issues: []string{},
	}

	// Find Python executable
	pythonExe := findPythonExecutable()
	if pythonExe == "" {
		env.Issues = append(env.Issues, "Python executable not found in PATH")
		return env
	}

	env.Executable = pythonExe
	env.Status = "healthy"

	// Get Python version
	output, err := utils.RunCommand(pythonExe, "--version")
	if err == nil {
		env.Version = strings.TrimSpace(output)
	}

	// Check for virtual environment
	env.VenvFound = isVenvActive()
	if env.VenvFound {
		env.VenvPath = os.Getenv("VIRTUAL_ENV")
	}

	// Detect package manager
	if utils.PathExists("requirements.txt") {
		env.PackageManager = "pip"
	} else if utils.PathExists("pyproject.toml") {
		// Could be pip, poetry, or uv
		env.PackageManager = "uv/poetry"
	}

	// List installed packages (limited)
	if utils.PathExists("requirements.txt") {
		data, _ := os.ReadFile("requirements.txt")
		packages := strings.Split(string(data), "\n")
		if len(packages) > 0 {
			env.Packages = packages[:min(len(packages), 10)]
		}
	}

	return env
}

func findPythonExecutable() string {
	// Try common names
	names := []string{"python", "python3", "python3.11", "python3.10"}
	if utils.IsWindows() {
		names = []string{"python.exe", "python3.exe", "python3.11.exe", "python3.10.exe"}
	}

	for _, name := range names {
		if exe, err := utils.FindExecutable(name); err == nil {
			return exe
		}
	}

	// Check common installation paths
	paths := getPythonPaths()
	for _, p := range paths {
		if utils.PathExists(p) {
			return p
		}
	}

	return ""
}

func getPythonPaths() []string {
	home := utils.GetHomeDir()
	paths := []string{}

	if utils.IsWindows() {
		paths = append(paths,
			filepath.Join(home, "AppData", "Local", "Programs", "Python", "Python311", "python.exe"),
			filepath.Join(home, "AppData", "Local", "Programs", "Python", "Python310", "python.exe"),
			"C:\\Python311\\python.exe",
			"C:\\Python310\\python.exe",
		)
	} else if utils.IsMacOS() {
		paths = append(paths,
			"/opt/homebrew/bin/python3",
			"/usr/local/bin/python3",
			filepath.Join(home, ".pyenv", "shims", "python"),
		)
	} else {
		paths = append(paths,
			"/usr/bin/python3",
			"/usr/local/bin/python3",
			filepath.Join(home, ".pyenv", "shims", "python"),
		)
	}

	return paths
}

func isVenvActive() bool {
	venv := os.Getenv("VIRTUAL_ENV")
	return venv != ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
