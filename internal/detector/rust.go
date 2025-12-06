package detector

import (
	"envfix/internal/types"
	"envfix/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

// DetectRust detects Rust environment
func (d *Detector) DetectRust() types.RustEnv {
	env := types.RustEnv{
		Status: "missing",
		Issues: []string{},
	}

	// Find rustc executable
	rustcExe := findRustcExecutable()
	if rustcExe == "" {
		env.Issues = append(env.Issues, "Rust compiler (rustc) not found")
		return env
	}

	env.Executable = rustcExe
	env.Status = "healthy"

	// Get Rust version
	output, err := utils.RunCommand(rustcExe, "--version")
	if err == nil {
		env.Version = strings.TrimSpace(output)
	}

	// Get CARGO_HOME
	cargoHome := os.Getenv("CARGO_HOME")
	if cargoHome == "" {
		cargoHome = filepath.Join(utils.GetHomeDir(), ".cargo")
	}
	env.CargoHome = cargoHome

	// Get RUSTUP_HOME
	rustupHome := os.Getenv("RUSTUP_HOME")
	if rustupHome == "" {
		rustupHome = filepath.Join(utils.GetHomeDir(), ".rustup")
	}
	env.RustupHome = rustupHome

	return env
}

func findRustcExecutable() string {
	names := []string{"rustc", "rustc.exe"}
	if !utils.IsWindows() {
		names = []string{"rustc"}
	}

	for _, name := range names {
		if exe, err := utils.FindExecutable(name); err == nil {
			return exe
		}
	}

	// Check common installation paths
	home := utils.GetHomeDir()
	cargoPath := filepath.Join(home, ".cargo", "bin", "rustc")
	if utils.IsWindows() {
		cargoPath = filepath.Join(home, ".cargo", "bin", "rustc.exe")
	}

	if utils.PathExists(cargoPath) {
		return cargoPath
	}

	return ""
}
