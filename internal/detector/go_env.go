package detector

import (
	"envfix/internal/types"
	"envfix/internal/utils"
	"os"
	"strings"
)

// DetectGo detects Go environment
func (d *Detector) DetectGo() types.GoEnv {
	env := types.GoEnv{
		Status: "missing",
		Issues: []string{},
	}

	// Find go executable
	goExe := findGoExecutable()
	if goExe == "" {
		env.Issues = append(env.Issues, "Go executable not found")
		return env
	}

	env.Executable = goExe
	env.Status = "healthy"

	// Get Go version
	output, err := utils.RunCommand(goExe, "version")
	if err == nil {
		env.Version = strings.TrimSpace(output)
	}

	// Get GOROOT
	goRoot := os.Getenv("GOROOT")
	if goRoot == "" {
		// Try to auto-detect
		output, _ := utils.RunCommand(goExe, "env", "GOROOT")
		goRoot = strings.TrimSpace(output)
	}
	env.GoRoot = goRoot

	// Get GOPATH
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		output, _ := utils.RunCommand(goExe, "env", "GOPATH")
		goPath = strings.TrimSpace(output)
	}
	env.GoPath = goPath

	// Check for go.mod (indicates Go modules)
	env.GoModules = utils.PathExists("go.mod")

	return env
}

func findGoExecutable() string {
	names := []string{"go", "go.exe"}
	if !utils.IsWindows() {
		names = []string{"go"}
	}

	for _, name := range names {
		if exe, err := utils.FindExecutable(name); err == nil {
			return exe
		}
	}

	// Check common installation paths
	paths := []string{}
	if utils.IsWindows() {
		paths = []string{
			"C:\\Program Files\\Go\\bin\\go.exe",
			"C:\\Go\\bin\\go.exe",
		}
	} else {
		paths = []string{
			"/usr/local/go/bin/go",
			"/usr/bin/go",
		}
	}

	for _, p := range paths {
		if utils.PathExists(p) {
			return p
		}
	}

	return ""
}
