package parser

import (
	"bufio"
	"encoding/json"
	"envfix/internal/utils"
	"os"
	"strings"
)

// ParseRequirementsTxt parses Python requirements.txt
func ParseRequirementsTxt() (map[string]string, error) {
	deps := make(map[string]string)

	if !utils.PathExists("requirements.txt") {
		return deps, nil
	}

	file, err := os.Open("requirements.txt")
	if err != nil {
		return deps, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle version specifiers
		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == '=' || r == '<' || r == '>' || r == '!' || r == '~'
		})

		if len(parts) > 0 {
			pkg := strings.TrimSpace(parts[0])
			version := ""

			// Extract version
			if idx := strings.IndexAny(line, "=<>!~"); idx != -1 {
				version = strings.TrimSpace(line[idx:])
			}

			deps[pkg] = version
		}
	}

	return deps, scanner.Err()
}

// ParsePackageJSON parses Node.js package.json
func ParsePackageJSON() (map[string]string, error) {
	deps := make(map[string]string)

	if !utils.PathExists("package.json") {
		return deps, nil
	}

	data, err := os.ReadFile("package.json")
	if err != nil {
		return deps, err
	}

	var packageJSON struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}

	if err := json.Unmarshal(data, &packageJSON); err != nil {
		return deps, err
	}

	// Merge both dependencies and devDependencies
	for pkg, version := range packageJSON.Dependencies {
		deps[pkg] = version
	}
	for pkg, version := range packageJSON.DevDependencies {
		deps[pkg] = version
	}

	return deps, nil
}

// ParseGoMod parses Go module dependencies
func ParseGoMod() (map[string]string, error) {
	deps := make(map[string]string)

	if !utils.PathExists("go.mod") {
		return deps, nil
	}

	file, err := os.Open("go.mod")
	if err != nil {
		return deps, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}

		if inRequire && line == ")" {
			inRequire = false
			continue
		}

		if inRequire && line != "" && !strings.HasPrefix(line, "//") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps[parts[0]] = parts[1]
			}
		} else if !inRequire && strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				deps[parts[1]] = parts[2]
			}
		}
	}

	return deps, scanner.Err()
}

// ParseCargoToml parses Rust Cargo.toml dependencies
func ParseCargoToml() (map[string]string, error) {
	deps := make(map[string]string)

	if !utils.PathExists("Cargo.toml") {
		return deps, nil
	}

	file, err := os.Open("Cargo.toml")
	if err != nil {
		return deps, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	inDependencies := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "[dependencies]") {
			inDependencies = true
			continue
		}

		if inDependencies && strings.HasPrefix(line, "[") {
			inDependencies = false
			continue
		}

		if inDependencies && line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.Split(line, "=")
			if len(parts) >= 2 {
				pkg := strings.TrimSpace(parts[0])
				version := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				deps[pkg] = version
			}
		}
	}

	return deps, scanner.Err()
}

// DependencyInfo holds information about a dependency
type DependencyInfo struct {
	Name    string
	Version string
	Manager string // pip, npm, cargo, go
	Type    string // dependency, dev-dependency
}

// GetAllDependencies aggregates dependencies from all package managers
func GetAllDependencies() []DependencyInfo {
	all := []DependencyInfo{}

	// Python dependencies
	if pyDeps, err := ParseRequirementsTxt(); err == nil {
		for pkg, version := range pyDeps {
			all = append(all, DependencyInfo{
				Name:    pkg,
				Version: version,
				Manager: "pip",
				Type:    "dependency",
			})
		}
	}

	// Node.js dependencies
	if nodeDeps, err := ParsePackageJSON(); err == nil {
		for pkg, version := range nodeDeps {
			all = append(all, DependencyInfo{
				Name:    pkg,
				Version: version,
				Manager: "npm",
				Type:    "dependency",
			})
		}
	}

	// Go dependencies
	if goDeps, err := ParseGoMod(); err == nil {
		for pkg, version := range goDeps {
			all = append(all, DependencyInfo{
				Name:    pkg,
				Version: version,
				Manager: "go",
				Type:    "dependency",
			})
		}
	}

	// Rust dependencies
	if rustDeps, err := ParseCargoToml(); err == nil {
		for pkg, version := range rustDeps {
			all = append(all, DependencyInfo{
				Name:    pkg,
				Version: version,
				Manager: "cargo",
				Type:    "dependency",
			})
		}
	}

	return all
}
