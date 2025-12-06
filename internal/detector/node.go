package detector

import (
	"envfix/internal/types"
	"envfix/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

// DetectNode detects Node.js environment
func (d *Detector) DetectNode() types.NodeEnv {
	env := types.NodeEnv{
		Status: "missing",
		Issues: []string{},
	}

	// Find Node executable
	nodeExe := findNodeExecutable()
	if nodeExe == "" {
		env.Issues = append(env.Issues, "Node.js executable not found in PATH")
		return env
	}

	env.Executable = nodeExe
	env.Status = "healthy"

	// Get Node version
	output, err := utils.RunCommand(nodeExe, "--version")
	if err == nil {
		env.Version = strings.TrimSpace(output)
	}

	// Get npm version
	if npmExe, err := utils.FindExecutable("npm"); err == nil {
		output, _ := utils.RunCommand(npmExe, "--version")
		env.NPMVersion = strings.TrimSpace(output)
	}

	// Check for pnpm
	if pnpmExe, err := utils.FindExecutable("pnpm"); err == nil {
		output, _ := utils.RunCommand(pnpmExe, "--version")
		env.PNPMVersion = strings.TrimSpace(output)
	}

	// Check for yarn
	if yarnExe, err := utils.FindExecutable("yarn"); err == nil {
		output, _ := utils.RunCommand(yarnExe, "--version")
		env.YarnVersion = strings.TrimSpace(output)
	}

	// Detect package manager in use
	if utils.PathExists("package-lock.json") {
		env.PackageManager = "npm"
	} else if utils.PathExists("pnpm-lock.yaml") {
		env.PackageManager = "pnpm"
	} else if utils.PathExists("yarn.lock") {
		env.PackageManager = "yarn"
	}

	// Check for node_modules
	env.ModulesFound = utils.PathExists("node_modules")
	if env.ModulesFound {
		absPath, _ := os.Getwd()
		env.ModulesPath = filepath.Join(absPath, "node_modules")
	}

	// Validate environment
	validateNodeEnv(&env)

	return env
}

func findNodeExecutable() string {
	names := []string{"node", "node.exe"}
	if !utils.IsWindows() {
		names = []string{"node"}
	}

	for _, name := range names {
		if exe, err := utils.FindExecutable(name); err == nil {
			return exe
		}
	}

	// Check common installation paths
	paths := getNodePaths()
	for _, p := range paths {
		if utils.PathExists(p) {
			return p
		}
	}

	return ""
}

func getNodePaths() []string {
	home := utils.GetHomeDir()
	var paths []string

	if utils.IsWindows() {
		paths = append(paths,
			filepath.Join(home, "AppData", "Local", "nodejs", "node.exe"),
			"C:\\Program Files\\nodejs\\node.exe",
			"C:\\Program Files (x86)\\nodejs\\node.exe",
		)
	} else if utils.IsMacOS() {
		paths = append(paths,
			"/opt/homebrew/bin/node",
			"/usr/local/bin/node",
			filepath.Join(home, ".nvm", "versions", "node", "v18.0.0", "bin", "node"),
		)
	} else {
		paths = append(paths,
			"/usr/bin/node",
			"/usr/local/bin/node",
			filepath.Join(home, ".nvm", "versions", "node", "v18.0.0", "bin", "node"),
		)
	}

	return paths
}

func validateNodeEnv(env *types.NodeEnv) {
	if env.ModulesFound && env.PackageManager == "" {
		env.Issues = append(env.Issues, "node_modules directory exists but no lock file found")
	}

	// Check for multiple lock files (conflict)
	lockFiles := 0
	if utils.PathExists("package-lock.json") {
		lockFiles++
	}
	if utils.PathExists("pnpm-lock.yaml") {
		lockFiles++
	}
	if utils.PathExists("yarn.lock") {
		lockFiles++
	}

	if lockFiles > 1 {
		env.Issues = append(env.Issues, "Multiple package manager lock files detected (npm, pnpm, yarn)")
	}
}
