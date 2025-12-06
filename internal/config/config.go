package config

import (
	"envfix/internal/utils"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents envfix configuration
type Config struct {
	Version         string        `yaml:"version"`
	Excludes        []string      `yaml:"excludes"`
	Include         IncludeConfig `yaml:"include"`
	Repair          RepairConfig  `yaml:"repair"`
	Clean           CleanConfig   `yaml:"clean"`
	Environment     EnvVars       `yaml:"environment"`
	NotificationURL string        `yaml:"notification_url"`
}

// IncludeConfig specifies what to include in scans
type IncludeConfig struct {
	Python bool `yaml:"python"`
	Node   bool `yaml:"node"`
	Rust   bool `yaml:"rust"`
	Go     bool `yaml:"go"`
}

// RepairConfig specifies repair behavior
type RepairConfig struct {
	CreateBackups bool   `yaml:"create_backups"`
	DryRun        bool   `yaml:"dry_run"`
	Strategy      string `yaml:"strategy"` // aggressive, conservative, interactive
}

// CleanConfig specifies cleanup behavior
type CleanConfig struct {
	RemoveCaches     bool `yaml:"remove_caches"`
	RemoveBackups    bool `yaml:"remove_backups"`
	MaxCacheAgeDays  int  `yaml:"max_cache_age_days"`
	DryRun           bool `yaml:"dry_run"`
	SkipConfirmation bool `yaml:"skip_confirmation"`
}

// EnvVars defines environment variables
type EnvVars map[string]string

var defaultConfig = Config{
	Version: "1.0",
	Excludes: []string{
		".git",
		"node_modules",
		".venv",
		"venv",
		"dist",
		"build",
	},
	Include: IncludeConfig{
		Python: true,
		Node:   true,
		Rust:   true,
		Go:     true,
	},
	Repair: RepairConfig{
		CreateBackups: true,
		DryRun:        false,
		Strategy:      "conservative",
	},
	Clean: CleanConfig{
		RemoveCaches:     true,
		RemoveBackups:    false,
		MaxCacheAgeDays:  30,
		DryRun:           false,
		SkipConfirmation: false,
	},
	Environment: EnvVars{},
}

// LoadConfig loads configuration from .envfix.yaml
func LoadConfig() (*Config, error) {
	cfg := defaultConfig

	// Look for .envfix.yaml in current directory
	if !utils.PathExists(".envfix.yaml") {
		return &cfg, nil
	}

	data, err := os.ReadFile(".envfix.yaml")
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves configuration to .envfix.yaml
func SaveConfig(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(".envfix.yaml", data, 0644)
}

// LoadGlobalConfig loads configuration from user home directory
func LoadGlobalConfig() (*Config, error) {
	home := utils.GetHomeDir()
	globalConfigPath := filepath.Join(home, ".envfix", "config.yaml")

	cfg := defaultConfig

	if !utils.PathExists(globalConfigPath) {
		return &cfg, nil
	}

	data, err := os.ReadFile(globalConfigPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveGlobalConfig saves configuration to user home directory
func SaveGlobalConfig(cfg *Config) error {
	home := utils.GetHomeDir()
	configDir := filepath.Join(home, ".envfix")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	return os.WriteFile(configPath, data, 0644)
}

// IsExcluded checks if a path should be excluded
func (c *Config) IsExcluded(path string) bool {
	for _, exclude := range c.Excludes {
		if matched, _ := filepath.Match(exclude, filepath.Base(path)); matched {
			return true
		}
	}
	return false
}

// GetDefault returns the default configuration
func GetDefault() *Config {
	return &defaultConfig
}
