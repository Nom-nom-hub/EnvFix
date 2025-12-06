package detector

import (
	"envfix/internal/types"
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name          string
		versionStr    string
		expectedMajor int
		expectedMinor int
		expectedPatch int
	}{
		{"Python 3.11.0", "3.11.0", 3, 11, 0},
		{"Node 18.12.0", "v18.12.0", 18, 12, 0},
		{"Node 16.13.2", "16.13.2", 16, 13, 2},
		{"Go 1.21.0", "go1.21.0", 1, 21, 0},
		{"Rust 1.73.0", "1.73.0", 1, 73, 0},
		{"No version", "invalid", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			major, minor, patch := ParseVersion(tt.versionStr)
			if major != tt.expectedMajor || minor != tt.expectedMinor || patch != tt.expectedPatch {
				t.Errorf("Expected %d.%d.%d, got %d.%d.%d",
					tt.expectedMajor, tt.expectedMinor, tt.expectedPatch,
					major, minor, patch)
			}
		})
	}
}

func TestCheckPythonMinimumVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		expect  bool
	}{
		{"Python 3.11", "3.11.0", true},
		{"Python 3.9", "3.9.0", true},
		{"Python 3.7", "3.7.0", true},
		{"Python 3.6", "3.6.0", false},
		{"Python 2.7", "2.7.18", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPythonMinimumVersion(tt.version)
			if result != tt.expect {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestCheckNodeMinimumVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		expect  bool
	}{
		{"Node 18", "18.12.0", true},
		{"Node 16", "16.13.2", true},
		{"Node 14", "14.21.0", true},
		{"Node 12", "12.22.0", false},
		{"Node 10", "10.24.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckNodeMinimumVersion(tt.version)
			if result != tt.expect {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestCheckRustMinimumVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		expect  bool
	}{
		{"Rust 1.73", "1.73.0", true},
		{"Rust 1.56", "1.56.0", true},
		{"Rust 1.50", "1.50.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckRustMinimumVersion(tt.version)
			if result != tt.expect {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestCheckGoMinimumVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		expect  bool
	}{
		{"Go 1.21", "1.21.0", true},
		{"Go 1.16", "1.16.0", true},
		{"Go 1.15", "1.15.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckGoMinimumVersion(tt.version)
			if result != tt.expect {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestIsMinimumVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		minMajor int
		minMinor int
		expect   bool
	}{
		{"Exact match", "3.7.0", 3, 7, true},
		{"Major higher", "4.0.0", 3, 7, true},
		{"Minor higher", "3.8.0", 3, 7, true},
		{"Below", "3.6.0", 3, 7, false},
		{"Major below", "2.7.0", 3, 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMinimumVersion(tt.version, tt.minMajor, tt.minMinor)
			if result != tt.expect {
				t.Errorf("Expected %v, got %v", tt.expect, result)
			}
		})
	}
}

func TestCheckVersionConflicts(t *testing.T) {
	d := &Detector{}
	result := &types.ScanResult{
		Python: types.PythonEnv{
			Status:     "healthy",
			Version:    "3.11.0",
			Executable: "python3",
		},
		Node: types.NodeEnv{
			Status:  "healthy",
			Version: "18.12.0",
		},
		Rust: types.RustEnv{
			Status:  "healthy",
			Version: "1.73.0",
		},
		Go: types.GoEnv{
			Status:  "healthy",
			Version: "1.21.0",
		},
		Issues:   []types.Issue{},
		Warnings: []string{},
	}

	d.CheckVersionConflicts(result)

	// Should have no issues for modern versions
	if len(result.Issues) > 0 {
		t.Errorf("Expected no issues for modern versions, got %d", len(result.Issues))
	}
}

func TestCheckVersionConflictsPython2(t *testing.T) {
	d := &Detector{}
	result := &types.ScanResult{
		Python: types.PythonEnv{
			Status:     "healthy",
			Version:    "Python 2.7.18",
			Executable: "python2",
		},
		Node: types.NodeEnv{
			Status: "missing",
		},
		Rust: types.RustEnv{
			Status: "missing",
		},
		Go: types.GoEnv{
			Status: "missing",
		},
		Issues:   []types.Issue{},
		Warnings: []string{},
	}

	d.CheckVersionConflicts(result)

	// Should detect Python 2 EOL
	foundPythonEOL := false
	for _, issue := range result.Issues {
		if issue.ID == "python_eol" {
			foundPythonEOL = true
			if issue.Severity != "critical" {
				t.Errorf("Python 2 EOL should be critical, got %s", issue.Severity)
			}
		}
	}

	if !foundPythonEOL {
		t.Errorf("Expected to detect Python 2 EOL")
	}
}

func TestCheckVersionConflictsNodeOld(t *testing.T) {
	d := &Detector{}
	result := &types.ScanResult{
		Python: types.PythonEnv{
			Status: "missing",
		},
		Node: types.NodeEnv{
			Status:  "healthy",
			Version: "12.22.0",
		},
		Rust: types.RustEnv{
			Status: "missing",
		},
		Go: types.GoEnv{
			Status: "missing",
		},
		Issues:   []types.Issue{},
		Warnings: []string{},
	}

	d.CheckVersionConflicts(result)

	// Should detect Node.js 12 EOL
	foundNodeEOL := false
	for _, issue := range result.Issues {
		if issue.ID == "node_eol" {
			foundNodeEOL = true
			if issue.Severity != "critical" {
				t.Errorf("Node 12 EOL should be critical, got %s", issue.Severity)
			}
		}
	}

	if !foundNodeEOL {
		t.Errorf("Expected to detect Node.js 12 EOL")
	}
}

func TestCheckVersionConflictsVersionLow(t *testing.T) {
	d := &Detector{}
	result := &types.ScanResult{
		Python: types.PythonEnv{
			Status:     "healthy",
			Version:    "3.6.0",
			Executable: "python3",
		},
		Node: types.NodeEnv{
			Status:  "healthy",
			Version: "14.21.0",
		},
		Rust: types.RustEnv{
			Status: "missing",
		},
		Go: types.GoEnv{
			Status: "missing",
		},
		Issues:   []types.Issue{},
		Warnings: []string{},
	}

	d.CheckVersionConflicts(result)

	// Should detect Python 3.6 is below minimum
	foundPythonLow := false
	for _, issue := range result.Issues {
		if issue.ID == "python_version_low" {
			foundPythonLow = true
		}
	}

	if !foundPythonLow {
		t.Errorf("Expected to detect Python version below minimum")
	}
}
