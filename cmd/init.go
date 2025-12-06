package cmd

import (
	"envfix/internal/config"
	"envfix/internal/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	initGlobal bool
	initForce  bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize envfix configuration",
	Long:  `Creates a .envfix.yaml configuration file in the current directory or global config.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetDefault()

		configPath := ".envfix.yaml"
		if initGlobal {
			utils.Info("Creating global configuration...")
			if err := config.SaveGlobalConfig(cfg); err != nil {
				utils.Error("Failed to save global config: %v", err)
				return err
			}
			utils.Success("✓ Global configuration created")
			return nil
		}

		// Check if config already exists
		if utils.PathExists(configPath) && !initForce {
			utils.Warning("Configuration file already exists at %s", configPath)
			fmt.Print("Overwrite? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Save local config
		if err := config.SaveConfig(cfg); err != nil {
			utils.Error("Failed to save config: %v", err)
			return err
		}

		utils.Success("✓ Configuration created at %s", configPath)

		// Show example
		fmt.Println("\nYou can now customize .envfix.yaml:")
		fmt.Println("  - Set repair strategy (conservative, aggressive, interactive)")
		fmt.Println("  - Exclude directories from scans")
		fmt.Println("  - Configure cache cleanup behavior")
		fmt.Println("  - Set environment variables")

		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&initGlobal, "global", false, "Create global configuration (~/.envfix/config.yaml)")
	initCmd.Flags().BoolVar(&initForce, "force", false, "Overwrite existing configuration")
	rootCmd.AddCommand(initCmd)
}
