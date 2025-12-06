package rebuilder

import (
	"envfix/internal/utils"
	"fmt"
	"os"
	"path/filepath"
)

// RepairPythonVenv rebuilds a Python virtual environment
func (r *Rebuilder) RepairPythonVenv() error {
	utils.Info("Repairing Python virtual environment...")

	pythonExe, err := utils.FindExecutable("python3")
	if err != nil {
		pythonExe, _ = utils.FindExecutable("python")
	}

	if pythonExe == "" {
		return fmt.Errorf("python executable not found")
	}

	venvPath := ".venv"

	// Backup existing venv if it exists
	if utils.PathExists(venvPath) {
		backupPath := venvPath + ".old"
		utils.Info("Backing up existing venv to %s", backupPath)
		if err := utils.RemoveDirectory(backupPath); err != nil {
			utils.Warning("Could not remove old backup: %v", err)
		}
		if err := os.Rename(venvPath, backupPath); err != nil {
			return fmt.Errorf("failed to backup venv: %w", err)
		}
	}

	// Create new venv
	utils.Info("Creating new virtual environment at %s", venvPath)
	_, err = utils.RunCommand(pythonExe, "-m", "venv", venvPath)
	if err != nil {
		return fmt.Errorf("failed to create venv: %w", err)
	}

	utils.Success("Python virtual environment repaired")
	return nil
}

// InstallPythonDependencies installs dependencies from requirements.txt
func (r *Rebuilder) InstallPythonDependencies() error {
	if !utils.PathExists("requirements.txt") {
		utils.Info("No requirements.txt found")
		return nil
	}

	utils.Info("Installing Python dependencies...")

	// Find pip
	venvPath := ".venv"
	var pipExe string

	if utils.PathExists(filepath.Join(venvPath, "Scripts", "pip.exe")) {
		pipExe = filepath.Join(venvPath, "Scripts", "pip.exe")
	} else if utils.PathExists(filepath.Join(venvPath, "bin", "pip")) {
		pipExe = filepath.Join(venvPath, "bin", "pip")
	} else {
		var err error
		pipExe, err = utils.FindExecutable("pip")
		if err != nil {
			return fmt.Errorf("pip not found")
		}
	}

	// Install dependencies
	_, err := utils.RunCommand(pipExe, "install", "-r", "requirements.txt")
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	utils.Success("Python dependencies installed")
	return nil
}

// CheckPythonPath validates Python executable in PATH
func (r *Rebuilder) CheckPythonPath() error {
	utils.Info("Checking Python PATH...")

	pythonExe, err := utils.FindExecutable("python3")
	if err != nil {
		pythonExe, err = utils.FindExecutable("python")
		if err != nil {
			return fmt.Errorf("python not in PATH")
		}
	}

	utils.Success("Python found at: %s", pythonExe)
	return nil
}
