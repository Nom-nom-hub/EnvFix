package detector

import (
	"envfix/internal/types"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIsOrphanedVenv(t *testing.T) {
	d := &Detector{}

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	venvDir := filepath.Join(tempDir, "test_venv")
	_ = os.MkdirAll(venvDir, 0755)

	// Create pyvenv.cfg to mark it as a venv
	configFile := filepath.Join(venvDir, "pyvenv.cfg")
	if err := os.WriteFile(configFile, []byte("home = /usr/bin\n"), 0644); err != nil {
		t.Fatalf("Failed to create pyvenv.cfg: %v", err)
	}

	// Test 1: Orphaned venv (no project files)
	if !d.isOrphanedVenv(venvDir) {
		t.Errorf("Expected venv to be orphaned when no project files exist")
	}

	// Test 2: Non-orphaned venv (has requirements.txt)
	requirementsFile := filepath.Join(tempDir, "requirements.txt")
	if err := os.WriteFile(requirementsFile, []byte("requests==2.28.0\n"), 0644); err != nil {
		t.Fatalf("Failed to create requirements.txt: %v", err)
	}

	if d.isOrphanedVenv(venvDir) {
		t.Errorf("Expected venv to NOT be orphaned when requirements.txt exists")
	}

	// Test 3: Non-existent pyvenv.cfg
	notVenvDir := filepath.Join(tempDir, "not_venv")
	_ = os.MkdirAll(notVenvDir, 0755)

	if d.isOrphanedVenv(notVenvDir) {
		t.Errorf("Expected non-venv directory to return false")
	}
}

func TestGetDirectoryAge(t *testing.T) {
	d := &Detector{}
	tempDir := t.TempDir()

	tests := []struct {
		name      string
		age       time.Duration
		expectStr string
	}{
		{"less than hour", 30 * time.Minute, "less than 1 hour"},
		{"1 day old", 24 * time.Hour, "1 days"},
		{"7 days old", 7 * 24 * time.Hour, "1 weeks"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := filepath.Join(tempDir, tt.name)
			_ = os.MkdirAll(testDir, 0755)

			// Modify time
			past := time.Now().Add(-tt.age)
			_ = os.Chtimes(testDir, past, past)

			age := d.getDirectoryAge(testDir)
			if age == "unknown" {
				t.Errorf("Expected valid age, got 'unknown'")
			}
		})
	}
}

func TestCheckOrphanedEnvironments(t *testing.T) {
	d := &Detector{}
	result := &types.ScanResult{
		Issues:   []types.Issue{},
		Warnings: []string{},
	}

	// Run the check
	d.CheckOrphanedEnvironments(result)

	// Should not error
	if result == nil {
		t.Errorf("Result should not be nil")
	}

	// In a clean environment, should have no or minimal issues
	t.Logf("Found %d issues, %d warnings", len(result.Issues), len(result.Warnings))
}

func TestIsOrphanedVenvWithPyproject(t *testing.T) {
	d := &Detector{}

	tempDir := t.TempDir()
	venvDir := filepath.Join(tempDir, "venv")
	_ = os.MkdirAll(venvDir, 0755)

	// Create pyvenv.cfg
	configFile := filepath.Join(venvDir, "pyvenv.cfg")
	_ = os.WriteFile(configFile, []byte("home = /usr/bin\n"), 0644)

	// Test with pyproject.toml
	pyprojectFile := filepath.Join(tempDir, "pyproject.toml")
	_ = os.WriteFile(pyprojectFile, []byte("[tool.poetry]\n"), 0644)

	if d.isOrphanedVenv(venvDir) {
		t.Errorf("Expected venv to NOT be orphaned when pyproject.toml exists")
	}
}

func TestIsOrphanedVenvWithSetup(t *testing.T) {
	d := &Detector{}

	tempDir := t.TempDir()
	venvDir := filepath.Join(tempDir, "venv")
	_ = os.MkdirAll(venvDir, 0755)

	// Create pyvenv.cfg
	configFile := filepath.Join(venvDir, "pyvenv.cfg")
	_ = os.WriteFile(configFile, []byte("home = /usr/bin\n"), 0644)

	// Test with setup.py
	setupFile := filepath.Join(tempDir, "setup.py")
	_ = os.WriteFile(setupFile, []byte("from setuptools import setup\n"), 0644)

	if d.isOrphanedVenv(venvDir) {
		t.Errorf("Expected venv to NOT be orphaned when setup.py exists")
	}
}
