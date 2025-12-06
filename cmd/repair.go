package cmd

import (
	"envfix/internal/rebuilder"
	"envfix/internal/utils"

	"github.com/spf13/cobra"
)

var (
	repairYes  bool
	repairDry bool
)

var repairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Repair detected environment issues",
	Long:  `Automatically fixes detected issues using safe, atomic operations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rb := rebuilder.NewRebuilder()
		rb.SetDryRun(repairDry)
		
		err := rb.Repair()
		if err != nil {
			utils.Error("Repair failed: %v", err)
			return err
		}

		return nil
	},
}

func init() {
	repairCmd.Flags().BoolVar(&repairYes, "yes", false, "Skip confirmation prompts")
	repairCmd.Flags().BoolVar(&repairDry, "dry-run", false, "Preview changes without applying them")
}
