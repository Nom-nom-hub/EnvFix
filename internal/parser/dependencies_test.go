package parser

import (
	"os"
	"testing"
)

func TestParseRequirementsTxt(t *testing.T) {
	// Create a temporary requirements.txt
	content := `requests==2.28.1
flask>=2.0.0
django~=4.0
# comment line
pytest

numpy==1.23.0
`

	tmpfile, err := os.Create("test_requirements.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Rename it to requirements.txt temporarily
	os.Rename(tmpfile.Name(), "requirements.txt")
	defer os.Remove("requirements.txt")

	deps, err := ParseRequirementsTxt()
	if err != nil {
		t.Errorf("ParseRequirementsTxt() error = %v", err)
	}

	expectedPackages := []string{"requests", "flask", "django", "pytest", "numpy"}
	for _, pkg := range expectedPackages {
		if _, exists := deps[pkg]; !exists {
			t.Errorf("ParseRequirementsTxt() missing package: %s", pkg)
		}
	}
}

func TestParsePackageJSON(t *testing.T) {
	content := `{
  "name": "test-project",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0",
    "axios": "^1.0.0"
  },
  "devDependencies": {
    "jest": "^29.0.0",
    "eslint": "^8.0.0"
  }
}
`

	tmpfile, err := os.Create("test_package.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	os.Rename(tmpfile.Name(), "package.json")
	defer os.Remove("package.json")

	deps, err := ParsePackageJSON()
	if err != nil {
		t.Errorf("ParsePackageJSON() error = %v", err)
	}

	expectedPackages := []string{"express", "axios", "jest", "eslint"}
	for _, pkg := range expectedPackages {
		if _, exists := deps[pkg]; !exists {
			t.Errorf("ParsePackageJSON() missing package: %s", pkg)
		}
	}
}

func TestGetAllDependencies(t *testing.T) {
	deps := GetAllDependencies()
	
	// Should return a slice (may be empty if no package files exist in test env)
	// Just verify it's not nil (Go initializes empty slices)
	if deps == nil {
		t.Errorf("GetAllDependencies() returned nil; want []DependencyInfo")
	}
	// It's ok if it's empty - we may not have package files in test environment
	_ = deps
}
