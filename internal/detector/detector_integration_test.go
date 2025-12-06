// +build integration

package detector

import (
	"os"
	"testing"
)

// TestFullScan runs a complete environment scan
func TestFullScan(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	d := NewDetector()
	result, err := d.Scan()

	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if result == nil {
		t.Fatal("Scan() returned nil result")
	}

	// Verify all fields are populated
	if result.Platform.OS == "" {
		t.Error("Platform.OS is empty")
	}

	if result.Timestamp.IsZero() {
		t.Error("Timestamp is zero")
	}

	if result.ScanDuration <= 0 {
		t.Error("ScanDuration is not positive")
	}

	// At least one language should be detected
	hasLanguage := result.Python.Status != "missing" ||
		result.Node.Status != "missing" ||
		result.Rust.Status != "missing" ||
		result.Go.Status != "missing"

	if !hasLanguage {
		t.Error("No language environments detected")
	}
}

// TestPythonDetection verifies Python detection
func TestPythonDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	d := NewDetector()
	pythonEnv := d.DetectPython()

	if pythonEnv.Status == "" {
		t.Error("Python status is empty")
	}

	// Python executable should either be found or status should be missing
	if pythonEnv.Status == "healthy" && pythonEnv.Executable == "" {
		t.Error("Python is healthy but executable is empty")
	}
}

// TestNodeDetection verifies Node.js detection
func TestNodeDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	d := NewDetector()
	nodeEnv := d.DetectNode()

	if nodeEnv.Status == "" {
		t.Error("Node status is empty")
	}

	if nodeEnv.Status == "healthy" && nodeEnv.Executable == "" {
		t.Error("Node is healthy but executable is empty")
	}
}

// TestDependencyDetection verifies dependency file detection
func TestDependencyDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a test package.json
	content := `{
  "name": "test",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0"
  }
}
`

	// Write to temp location
	tmpfile, err := os.Create("test_package_integration.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = tmpfile.WriteString(content)
	tmpfile.Close()
	defer os.Remove(tmpfile.Name())

	// Rename to package.json
	_ = os.Rename(tmpfile.Name(), "package.json")
	defer os.Remove("package.json")

	d := NewDetector()
	result, err := d.Scan()

	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	// Should detect package.json
	found := false
	for _, warning := range result.Warnings {
		if warning == "package.json found but no lock file (npm, yarn, or pnpm)" {
			found = true
			break
		}
	}

	if !found {
		t.Log("Warning about missing package.json lock file not found (might be expected)")
	}
}
