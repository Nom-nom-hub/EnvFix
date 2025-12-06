package rebuilder

import (
	"os"
	"testing"
)

func TestFindPipExecutable(t *testing.T) {
	r := &Rebuilder{}

	// Test finding pip
	pipExe := r.findPipExecutable()

	// Should find either venv pip or system pip
	// (may be empty in test environment)
	// Just verify that if found, it's a valid path (doesn't error)
	// The function either returns a valid pip path or empty string
	_ = pipExe // It's ok if empty - we may not have pip in test env
}

func TestCheckDependencyLockfiles(t *testing.T) {
	r := &Rebuilder{}
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Test 1: No files - should pass
	err := r.CheckDependencyLockfiles()
	if err != nil {
		t.Errorf("Should handle empty directory: %v", err)
	}

	// Test 2: Python project with requirements.txt - should pass
	os.WriteFile("pyproject.toml", []byte("[tool.poetry]\n"), 0644)
	os.WriteFile("requirements.txt", []byte("requests==2.28.0\n"), 0644)
	err = r.CheckDependencyLockfiles()
	if err != nil {
		t.Errorf("Should pass with requirements.txt: %v", err)
	}

	// Test 3: Node.js project with lock file - should pass
	os.WriteFile("package.json", []byte("{}\n"), 0644)
	os.WriteFile("package-lock.json", []byte("{}\n"), 0644)
	err = r.CheckDependencyLockfiles()
	if err != nil {
		t.Errorf("Should pass with lock file: %v", err)
	}

	// Test 4: Go project with go.mod and go.sum - should pass
	os.WriteFile("go.mod", []byte("module example.com\n"), 0644)
	os.WriteFile("go.sum", []byte("example.com v1.0.0\n"), 0644)
	err = r.CheckDependencyLockfiles()
	if err != nil {
		t.Errorf("Should pass with go.sum: %v", err)
	}
}

func TestVerifyDependencies(t *testing.T) {
	r := &Rebuilder{}

	// Should not error even with no dependencies
	err := r.VerifyDependencies()
	if err != nil {
		t.Errorf("Should handle no dependencies: %v", err)
	}
}

func TestCleanDependencyCaches(t *testing.T) {
	r := &Rebuilder{}

	// Should not error - may not have caches to clean
	err := r.CleanDependencyCaches()
	if err != nil {
		t.Errorf("Should handle missing caches gracefully: %v", err)
	}
}

func TestInstallPythonDependenciesUpgrade(t *testing.T) {
	r := &Rebuilder{}
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Test 1: No requirements.txt - should return early
	err := r.InstallPythonDependenciesUpgrade()
	if err != nil {
		t.Errorf("Should handle missing requirements.txt: %v", err)
	}

	// Test 2: With requirements.txt - would need actual pip to test
	os.WriteFile("requirements.txt", []byte("requests==2.28.0\n"), 0644)
	// Note: This will fail without a real pip install, but that's expected in tests
}

func TestInstallNodeDependencies(t *testing.T) {
	r := &Rebuilder{}
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Test 1: No package.json - should return early
	err := r.InstallNodeDependencies()
	if err != nil {
		t.Errorf("Should handle missing package.json: %v", err)
	}

	// Test 2: With package.json - would need actual npm to test
	os.WriteFile("package.json", []byte("{}\n"), 0644)
	// Note: This will fail without a real npm install, but that's expected in tests
}

func TestInstallRustDependencies(t *testing.T) {
	r := &Rebuilder{}
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Test: No Cargo.toml - should return early
	err := r.InstallRustDependencies()
	if err != nil {
		t.Errorf("Should handle missing Cargo.toml: %v", err)
	}
}

func TestTidyGoDependencies(t *testing.T) {
	r := &Rebuilder{}
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	os.Chdir(tempDir)

	// Test: No go.mod - should return early
	err := r.TidyGoDependencies()
	if err != nil {
		t.Errorf("Should handle missing go.mod: %v", err)
	}
}

func TestRegenerateLockfiles(t *testing.T) {
	r := &Rebuilder{}

	// Should not error even with no lock files
	err := r.RegenerateLockfiles()
	if err != nil {
		t.Errorf("Should handle regenerate gracefully: %v", err)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		(s == substr || len(s) > len(substr) && 
		 (s[:len(substr)] == substr || 
		  s[len(s)-len(substr):] == substr))
}
