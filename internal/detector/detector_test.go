package detector

import (
	"testing"
)

func TestNewDetector(t *testing.T) {
	d := NewDetector()
	if d == nil {
		t.Errorf("NewDetector() = nil; want detector instance")
	}
}

func TestScan(t *testing.T) {
	d := NewDetector()
	result, err := d.Scan()

	if err != nil {
		t.Errorf("Scan() returned error: %v", err)
	}

	if result == nil {
		t.Errorf("Scan() returned nil result")
	}

	// Verify result structure
	if result.Platform.OS == "" {
		t.Errorf("Scan() result.Platform.OS is empty")
	}

	if result.ScanDuration < 0 {
		t.Errorf("Scan() result.ScanDuration is negative")
	}
}

func TestDetectPython(t *testing.T) {
	d := NewDetector()
	pythonEnv := d.DetectPython()

	// Should have a status
	if pythonEnv.Status == "" {
		t.Errorf("DetectPython() status is empty")
	}

	// Status should be one of: missing, broken, healthy
	validStatuses := map[string]bool{
		"missing": true,
		"broken":  true,
		"healthy": true,
	}
	if !validStatuses[pythonEnv.Status] {
		t.Errorf("DetectPython() status = %q; want missing, broken, or healthy", pythonEnv.Status)
	}
}

func TestDetectNode(t *testing.T) {
	d := NewDetector()
	nodeEnv := d.DetectNode()

	if nodeEnv.Status == "" {
		t.Errorf("DetectNode() status is empty")
	}

	validStatuses := map[string]bool{
		"missing": true,
		"broken":  true,
		"healthy": true,
	}
	if !validStatuses[nodeEnv.Status] {
		t.Errorf("DetectNode() status = %q; want missing, broken, or healthy", nodeEnv.Status)
	}
}

func TestDetectRust(t *testing.T) {
	d := NewDetector()
	rustEnv := d.DetectRust()

	if rustEnv.Status == "" {
		t.Errorf("DetectRust() status is empty")
	}
}

func TestDetectGo(t *testing.T) {
	d := NewDetector()
	goEnv := d.DetectGo()

	if goEnv.Status == "" {
		t.Errorf("DetectGo() status is empty")
	}
}
