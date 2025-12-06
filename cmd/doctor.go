package cmd

import (
	"envfix/internal/utils"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Scan and repair environment (convenience command)",
	Long:  `Combines scan + repair operations for one-shot environment diagnosis and fixing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.Info("Running environment doctor...")
		
		// Run scan first
		if err := scanCmd.RunE(cmd, args); err != nil {
			return err
		}

		// Then run repair
		if err := repairCmd.RunE(cmd, args); err != nil {
			return err
		}

		utils.Success("✓ Doctor finished")
		return nil
	},
}
