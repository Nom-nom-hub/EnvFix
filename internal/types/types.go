package types

import "time"

// ScanResult contains all detected issues in the environment
type ScanResult struct {
	Timestamp    time.Time    `yaml:"timestamp"`
	Platform     PlatformInfo `yaml:"platform"`
	Python       PythonEnv    `yaml:"python"`
	Node         NodeEnv      `yaml:"nodejs"`
	Rust         RustEnv      `yaml:"rust"`
	Go           GoEnv        `yaml:"go"`
	Issues       []Issue      `yaml:"issues"`
	Warnings     []string     `yaml:"warnings"`
	IsHealthy    bool         `yaml:"is_healthy"`
	ScanDuration float64      `yaml:"scan_duration_seconds"`
}

// PlatformInfo describes the system
type PlatformInfo struct {
	OS       string `yaml:"os"`
	Arch     string `yaml:"arch"`
	Version  string `yaml:"version"`
	HomeDir  string `yaml:"home_dir"`
	PathSep  string `yaml:"path_sep"`
	LineEnds string `yaml:"line_endings"`
}

// PythonEnv describes Python environment state
type PythonEnv struct {
	Version        string   `yaml:"version"`
	Executable     string   `yaml:"executable"`
	VenvFound      bool     `yaml:"venv_found"`
	VenvPath       string   `yaml:"venv_path"`
	PackageManager string   `yaml:"package_manager"`
	Packages       []string `yaml:"packages"`
	Issues         []string `yaml:"issues"`
	Status         string   `yaml:"status"` // healthy, broken, missing
}

// NodeEnv describes Node.js environment state
type NodeEnv struct {
	Version        string   `yaml:"version"`
	Executable     string   `yaml:"executable"`
	PackageManager string   `yaml:"package_manager"`
	NPMVersion     string   `yaml:"npm_version"`
	PNPMVersion    string   `yaml:"pnpm_version"`
	YarnVersion    string   `yaml:"yarn_version"`
	ModulesFound   bool     `yaml:"node_modules_found"`
	ModulesPath    string   `yaml:"node_modules_path"`
	Issues         []string `yaml:"issues"`
	Status         string   `yaml:"status"`
}

// RustEnv describes Rust environment state
type RustEnv struct {
	Version    string   `yaml:"version"`
	Executable string   `yaml:"executable"`
	CargoHome  string   `yaml:"cargo_home"`
	RustupHome string   `yaml:"rustup_home"`
	Issues     []string `yaml:"issues"`
	Status     string   `yaml:"status"`
}

// GoEnv describes Go environment state
type GoEnv struct {
	Version    string   `yaml:"version"`
	Executable string   `yaml:"executable"`
	GoRoot     string   `yaml:"go_root"`
	GoPath     string   `yaml:"go_path"`
	GoModules  bool     `yaml:"go_modules"`
	Issues     []string `yaml:"issues"`
	Status     string   `yaml:"status"`
}

// Issue represents a detected problem
type Issue struct {
	ID       string `yaml:"id"`
	Title    string `yaml:"title"`
	Severity string `yaml:"severity"` // critical, warning, info
	Language string `yaml:"language"`
	Message  string `yaml:"message"`
	Path     string `yaml:"path,omitempty"`
	Fix      string `yaml:"fix,omitempty"`
}

// EnvironmentManifest is the env.yaml structure
type EnvironmentManifest struct {
	Version      string            `yaml:"version"`
	Timestamp    time.Time         `yaml:"timestamp"`
	Platform     PlatformInfo      `yaml:"platform"`
	Python       PythonSpec        `yaml:"python,omitempty"`
	Node         NodeSpec          `yaml:"nodejs,omitempty"`
	Rust         RustSpec          `yaml:"rust,omitempty"`
	Go           GoSpec            `yaml:"go,omitempty"`
	Dependencies map[string]string `yaml:"dependencies,omitempty"`
	Checksums    map[string]string `yaml:"checksums,omitempty"`
	Environment  map[string]string `yaml:"environment,omitempty"`
}

type PythonSpec struct {
	Version        string `yaml:"version"`
	VenvPath       string `yaml:"venv_path"`
	PackageManager string `yaml:"package_manager"`
}

type NodeSpec struct {
	Version        string `yaml:"version"`
	PackageManager string `yaml:"package_manager"`
}

type RustSpec struct {
	Version string `yaml:"version"`
}

type GoSpec struct {
	Version string `yaml:"version"`
}

// RepairAction describes an action to fix an issue
type RepairAction struct {
	ID          string
	Description string
	Action      func() error
	Reversible  bool
}
