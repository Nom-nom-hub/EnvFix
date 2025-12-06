package cmd

import (
	"envfix/internal/spec"
	"envfix/internal/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	lockOutput string
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Generate environment manifest",
	Long:  `Generates an environment manifest (env.yaml) capturing language versions, toolchain versions, dependencies, and more.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.Info("Generating environment manifest...")

		gen := spec.NewGenerator()
		yamlBytes, err := gen.GenerateYAML()
		if err != nil {
			utils.Error("Manifest generation failed: %v", err)
			return err
		}

		// Write to file
		if err := os.WriteFile(lockOutput, yamlBytes, 0644); err != nil {
			utils.Error("Failed to write manifest: %v", err)
			return err
		}

		utils.Success("✓ Manifest written to %s", lockOutput)

		// Also print to stdout
		fmt.Println("\n" + string(yamlBytes))
		return nil
	},
}

func init() {
	lockCmd.Flags().StringVar(&lockOutput, "output", "env.yaml", "Output file for manifest")
}
