package rebuilder

import (
	"envfix/internal/detector"
	"envfix/internal/types"
	"envfix/internal/utils"
)

type Rebuilder struct {
	detector *detector.Detector
	dry_run  bool
	actions  []types.RepairAction
}

func NewRebuilder() *Rebuilder {
	return &Rebuilder{
		detector: detector.NewDetector(),
		dry_run:  false,
		actions:  []types.RepairAction{},
	}
}

// Repair fixes detected environment issues
func (r *Rebuilder) Repair() error {
	utils.Info("Starting environment repair...")

	// First, scan for issues
	result, err := r.detector.Scan()
	if err != nil {
		return err
	}

	if result.IsHealthy {
		utils.Success("✓ Environment is healthy, no repairs needed")
		return nil
	}

	// Build repair plan
	r.planRepairs(result)

	// Execute repairs
	if err := r.executeRepairs(); err != nil {
		return err
	}

	utils.Success("✓ Environment repair complete")
	return nil
}

// planRepairs creates a list of repair actions based on detected issues
func (r *Rebuilder) planRepairs(result *types.ScanResult) {
	r.actions = []types.RepairAction{}

	// Python repairs
	if result.Python.Status == "missing" {
		utils.Warning("Python not found in PATH - manual installation required")
	} else if len(result.Python.Issues) > 0 {
		r.actions = append(r.actions, types.RepairAction{
			ID:          "python_venv",
			Description: "Rebuild Python virtual environment",
			Action:      r.RepairPythonVenv,
			Reversible:  true,
		})

		if utils.PathExists("requirements.txt") {
			r.actions = append(r.actions, types.RepairAction{
				ID:          "python_deps",
				Description: "Install Python dependencies",
				Action:      r.InstallPythonDependencies,
				Reversible:  false,
			})
		}
	}

	// Node.js repairs
	if result.Node.Status == "healthy" && len(result.Node.Issues) > 0 {
		// Check for package manager conflicts
		hasConflict := false
		for _, issue := range result.Node.Issues {
			if issue == "Multiple package manager lock files detected (npm, pnpm, yarn)" {
				hasConflict = true
				break
			}
		}

		if hasConflict {
			r.actions = append(r.actions, types.RepairAction{
				ID:          "node_pm_conflict",
				Description: "Resolve Node.js package manager conflicts",
				Action:      r.ResolveNodePackageManagerConflict,
				Reversible:  true,
			})
		}

		r.actions = append(r.actions, types.RepairAction{
			ID:          "node_modules",
			Description: "Rebuild Node.js modules",
			Action:      r.RepairNodeModules,
			Reversible:  true,
		})

		r.actions = append(r.actions, types.RepairAction{
			ID:          "node_cache",
			Description: "Clean Node.js cache",
			Action:      r.CleanNodeCache,
			Reversible:  false,
		})
	}

	// Check for orphaned environments (warnings only, not critical)
	hasOrphanedIssues := false
	for _, issue := range result.Issues {
		if issue.ID == "orphaned_venv" {
			hasOrphanedIssues = true
			break
		}
	}

	if hasOrphanedIssues {
		r.actions = append(r.actions, types.RepairAction{
			ID:          "orphaned_envs",
			Description: "Remove orphaned virtual environments",
			Action:      r.RemoveOrphanedEnvironments,
			Reversible:  false,
		})
	}

	// Add dependency management actions
	r.actions = append(r.actions, types.RepairAction{
		ID:          "verify_deps",
		Description: "Verify all dependencies are properly installed",
		Action:      r.VerifyDependencies,
		Reversible:  false,
	})

	r.actions = append(r.actions, types.RepairAction{
		ID:          "check_lockfiles",
		Description: "Check dependency lock file consistency",
		Action:      r.CheckDependencyLockfiles,
		Reversible:  false,
	})

	// Add version validation actions
	r.actions = append(r.actions, types.RepairAction{
		ID:          "validate_versions",
		Description: "Validate interpreter version compatibility",
		Action:      r.ValidateVersionCompatibility,
		Reversible:  false,
	})
}

// executeRepairs runs all planned repair actions
func (r *Rebuilder) executeRepairs() error {
	if len(r.actions) == 0 {
		utils.Info("No repairs needed")
		return nil
	}

	utils.Info("\nExecuting %d repair action(s)...\n", len(r.actions))

	for i, action := range r.actions {
		utils.Info("[%d/%d] %s", i+1, len(r.actions), action.Description)

		if err := action.Action(); err != nil {
			utils.Error("Failed: %s - %v", action.ID, err)
			if action.Reversible {
				utils.Info("This operation has a backup (*.old or *.backup)")
			}
			return err
		}

		utils.Success("✓ Completed: %s", action.Description)
	}

	utils.Success("\n✓ All repairs completed successfully")
	return nil
}

// SetDryRun enables dry-run mode
func (r *Rebuilder) SetDryRun(dryRun bool) {
	r.dry_run = dryRun
}

// GetActions returns the list of planned repair actions
func (r *Rebuilder) GetActions() []types.RepairAction {
	return r.actions
}
