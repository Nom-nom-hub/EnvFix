package spec

import (
	"envfix/internal/detector"
	"envfix/internal/types"
	"envfix/internal/utils"
	"time"

	"gopkg.in/yaml.v3"
)

type Generator struct {
	detector *detector.Detector
}

func NewGenerator() *Generator {
	return &Generator{
		detector: detector.NewDetector(),
	}
}

// Generate creates an environment manifest
func (g *Generator) Generate() (*types.EnvironmentManifest, error) {
	// Scan environment first
	scanResult, err := g.detector.Scan()
	if err != nil {
		return nil, err
	}

	osName, arch, _ := utils.GetSystemInfo()

	manifest := &types.EnvironmentManifest{
		Version:   "1.0",
		Timestamp: time.Now(),
		Platform: types.PlatformInfo{
			OS:   osName,
			Arch: arch,
		},
		Dependencies: make(map[string]string),
		Checksums:    make(map[string]string),
		Environment:  make(map[string]string),
	}

	// Add Python spec if detected
	if scanResult.Python.Status != "missing" {
		manifest.Python = types.PythonSpec{
			Version:        scanResult.Python.Version,
			VenvPath:       scanResult.Python.VenvPath,
			PackageManager: scanResult.Python.PackageManager,
		}
	}

	// Add Node spec if detected
	if scanResult.Node.Status != "missing" {
		manifest.Node = types.NodeSpec{
			Version:        scanResult.Node.Version,
			PackageManager: scanResult.Node.PackageManager,
		}
	}

	// Add Rust spec if detected
	if scanResult.Rust.Status != "missing" {
		manifest.Rust = types.RustSpec{
			Version: scanResult.Rust.Version,
		}
	}

	// Add Go spec if detected
	if scanResult.Go.Status != "missing" {
		manifest.Go = types.GoSpec{
			Version: scanResult.Go.Version,
		}
	}

	return manifest, nil
}

// GenerateYAML converts manifest to YAML
func (g *Generator) GenerateYAML() ([]byte, error) {
	manifest, err := g.Generate()
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(manifest)
}

// GenerateString converts manifest to YAML string
func (g *Generator) GenerateString() (string, error) {
	manifest, err := g.Generate()
	if err != nil {
		return "", err
	}

	data, _ := yaml.Marshal(manifest)
	return string(data), nil
}
