package utils

import (
	"os"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	home := GetHomeDir()
	if home == "" || home == "unknown" {
		t.Errorf("GetHomeDir() = %q; want non-empty home directory", home)
	}
}

func TestIsWindows(t *testing.T) {
	// This will vary by platform, just test it doesn't panic
	_ = IsWindows()
}

func TestIsMacOS(t *testing.T) {
	_ = IsMacOS()
}

func TestIsLinux(t *testing.T) {
	_ = IsLinux()
}

func TestPathExists(t *testing.T) {
	// Test with a file that should exist
	home := GetHomeDir()
	if !PathExists(home) {
		t.Errorf("PathExists(%q) = false; want true for home directory", home)
	}

	// Test with a nonexistent path
	if PathExists("/nonexistent/path/to/file") {
		t.Errorf("PathExists() = true; want false for nonexistent path")
	}
}

func TestExpandPath(t *testing.T) {
	// Test ~ expansion
	result := ExpandPath("~/test")
	if result == "~/test" {
		t.Errorf("ExpandPath(~/test) did not expand ~")
	}

	// Test absolute path passthrough
	result = ExpandPath("/tmp/test")
	if result != "/tmp/test" {
		t.Errorf("ExpandPath(/tmp/test) = %q; want /tmp/test", result)
	}
}

func TestGetVersionFromOutput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Python 3.11.0\n", "Python 3.11.0"},
		{"v18.0.0\nmore info", "v18.0.0"},
		{"rustc 1.70.0", "rustc 1.70.0"},
		{"", ""},
	}

	for _, tt := range tests {
		result := GetVersionFromOutput(tt.input)
		if result != tt.expected {
			t.Errorf("GetVersionFromOutput(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetPathSeparator(t *testing.T) {
	sep := GetPathSeparator()
	if sep != ";" && sep != ":" {
		t.Errorf("GetPathSeparator() = %q; want : or ;", sep)
	}
}

func TestSplitPath(t *testing.T) {
	paths := SplitPath()
	if len(paths) == 0 {
		t.Errorf("SplitPath() returned empty slice")
	}
}

func TestReadEnvVariable(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "test_value")
	defer os.Unsetenv("TEST_ENV_VAR")

	result := ReadEnvVariable("TEST_ENV_VAR")
	if result != "test_value" {
		t.Errorf("ReadEnvVariable() = %q; want test_value", result)
	}
}

func TestGetSystemInfo(t *testing.T) {
	osName, arch, home := GetSystemInfo()

	if osName == "" {
		t.Errorf("GetSystemInfo() osName = empty string")
	}

	if arch == "" {
		t.Errorf("GetSystemInfo() arch = empty string")
	}

	if home == "" {
		t.Errorf("GetSystemInfo() home = empty string")
	}
}
