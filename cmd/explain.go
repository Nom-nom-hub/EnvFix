package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var explainCmd = &cobra.Command{
	Use:   "explain <issue>",
	Short: "Explain an issue in human-readable terms",
	Long:  `Provides human-readable explanation for detected issues and why they matter.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issue := args[0]
		explanation := getExplanation(issue)

		fmt.Printf("\nIssue: %s\n", issue)
		fmt.Printf("Explanation:\n%s\n\n", explanation)
		return nil
	},
}

func getExplanation(issue string) string {
	explanations := map[string]string{
		"python_missing": `Python is not installed or not in your PATH. This means you can't run Python scripts directly.
Fix: Install Python from python.org or your package manager.
		
Windows: Visit python.org or use winget install Python.3.11
macOS: brew install python3
Linux: sudo apt-get install python3 (Ubuntu/Debian) or sudo yum install python3 (RHEL/Fedora)`,

		"node_missing": `Node.js is not installed or not in your PATH. This prevents you from running JavaScript projects.
Fix: Install Node.js from nodejs.org or your package manager.

Windows: Visit nodejs.org or use winget install OpenJS.NodeJS
macOS: brew install node
Linux: Use your package manager (apt, yum, pacman)`,

		"venv_corrupt": `Your Python virtual environment is corrupted or broken. This means dependencies may not load correctly.
Fix: Run 'envfix repair' to rebuild the virtual environment.`,

		"package_manager_conflict": `Multiple Node.js package managers detected (npm, yarn, pnpm). This causes dependency conflicts.
Fix: Run 'envfix repair' to resolve the conflict and use a single manager.`,

		"broken_symlink": `A broken symlink exists in your PATH. This can slow down command lookup.
Fix: Run 'envfix clean' to remove broken symlinks.`,
	}

	if exp, exists := explanations[issue]; exists {
		return exp
	}

	return fmt.Sprintf(`Unknown issue: %s

Run 'envfix scan' to see all detected issues, or 'envfix repair' to fix them automatically.`, issue)
}
