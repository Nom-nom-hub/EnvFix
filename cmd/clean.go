package cmd

import (
	"envfix/internal/cleaner"
	"envfix/internal/utils"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean unused environment artifacts",
	Long:  `Removes leftover, unused, orphaned, or broken environment assets.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cl := cleaner.NewCleaner()
		err := cl.Clean()
		if err != nil {
			utils.Error("Clean failed: %v", err)
			return err
		}

		return nil
	},
}
