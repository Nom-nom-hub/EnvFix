package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// GetSystemInfo returns platform information
func GetSystemInfo() (os string, arch string, homeDir string) {
	return runtime.GOOS, runtime.GOARCH, GetHomeDir()
}

// GetHomeDir returns the user's home directory
func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "unknown"
	}
	return home
}

// FindExecutable locates a command in PATH
func FindExecutable(name string) (string, error) {
	return exec.LookPath(name)
}

// RunCommand executes a command and returns output
func RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// IsWindows returns true on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsMacOS returns true on macOS
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux returns true on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// GetPathSeparator returns the OS-specific PATH separator
func GetPathSeparator() string {
	if IsWindows() {
		return ";"
	}
	return ":"
}

// SplitPath splits PATH environment variable
func SplitPath() []string {
	pathVar := os.Getenv("PATH")
	return strings.Split(pathVar, GetPathSeparator())
}

// PathExists checks if a path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsSymlink checks if a path is a symlink
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// FindInPaths searches for a file in a list of directories
func FindInPaths(filename string, paths []string) string {
	for _, dir := range paths {
		fullPath := filepath.Join(dir, filename)
		if PathExists(fullPath) {
			return fullPath
		}
		// On Windows, also try with .exe and .cmd extensions
		if IsWindows() {
			if PathExists(fullPath + ".exe") {
				return fullPath + ".exe"
			}
			if PathExists(fullPath + ".cmd") {
				return fullPath + ".cmd"
			}
		}
	}
	return ""
}

// GetVersionFromOutput extracts version from command output
func GetVersionFromOutput(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

// ReadEnvVariable reads an environment variable
func ReadEnvVariable(key string) string {
	return os.Getenv(key)
}

// ExpandPath expands ~ and environment variables in a path
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home := GetHomeDir()
		return filepath.Join(home, path[1:])
	}
	return os.ExpandEnv(path)
}

// DirectorySize returns the size of a directory in bytes
func DirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// ListDirectory lists entries in a directory
func ListDirectory(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

// CopyFile copies a file
func CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// RemoveDirectory safely removes a directory
func RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}

// CreateBackup creates a backup of a file
func CreateBackup(path string) (string, error) {
	backupPath := path + ".backup"
	if err := CopyFile(path, backupPath); err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}
	return backupPath, nil
}
