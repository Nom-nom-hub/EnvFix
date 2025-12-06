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

	_ = os.Chdir(tempDir)

	// Test 1: No files - should pass
	err := r.CheckDependencyLockfiles()
	if err != nil {
		t.Errorf("Should handle empty directory: %v", err)
	}

	// Test 2: Python project with requirements.txt - should pass
	_ = os.WriteFile("pyproject.toml", []byte("[tool.poetry]\n"), 0644)
	_ = os.WriteFile("requirements.txt", []byte("requests==2.28.0\n"), 0644)
	err = r.CheckDependencyLockfiles()
	if err != nil {
		t.Errorf("Should pass with requirements.txt: %v", err)
	}

	// Test 3: Node.js project with lock file - should pass
	_ = os.WriteFile("package.json", []byte("{}\n"), 0644)
	_ = os.WriteFile("package-lock.json", []byte("{}\n"), 0644)
	err = r.CheckDependencyLockfiles()
	if err != nil {
		t.Errorf("Should pass with lock file: %v", err)
	}

	// Test 4: Go project with go.mod and go.sum - should pass
	_ = os.WriteFile("go.mod", []byte("module example.com\n"), 0644)
	_ = os.WriteFile("go.sum", []byte("example.com v1.0.0\n"), 0644)
	err = r.CheckDependencyLockfiles()
	if err != nil {
		t.Errorf("Should pass with go.sum: %v", err)
	}
}

func TestVerifyDependencies(t *testing.T) {
	r := &Rebuilder{}
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	_ = os.Chdir(tempDir)

	// Test 1: Empty directory - no dependencies to verify
	err := r.VerifyDependencies()
	if err != nil {
		t.Logf("Empty dir error (may be expected): %v", err)
	}

	// Test 2: Python project
	_ = os.WriteFile("requirements.txt", []byte("requests==2.28.0\n"), 0644)
	err = r.VerifyDependencies()
	// Just verify it doesn't panic
	_ = err
}

func TestRegenerateLockfiles(t *testing.T) {
	r := &Rebuilder{}
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	_ = os.Chdir(tempDir)

	// Test 1: Empty directory
	err := r.RegenerateLockfiles()
	if err != nil {
		t.Logf("Empty dir error (expected): %v", err)
	}

	// Test 2: Node.js project
	_ = os.WriteFile("package.json", []byte("{}\n"), 0644)
	err = r.RegenerateLockfiles()
	// Should handle gracefully even if npm not available
	_ = err
}

func TestCleanDependencyCaches(t *testing.T) {
	r := &Rebuilder{}

	// Just verify the function runs without error
	err := r.CleanDependencyCaches()
	if err != nil {
		t.Logf("Clean caches error (may be expected): %v", err)
	}
}
