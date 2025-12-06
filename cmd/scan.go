package cmd

import (
	"envfix/internal/detector"
	"envfix/internal/utils"
	"fmt"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan environment for issues",
	Long:  `Performs a deep diagnostic scan of the environment and current project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		det := detector.NewDetector()
		result, err := det.Scan()
		if err != nil {
			utils.Error("Scan failed: %v", err)
			return err
		}

		// Display results
		fmt.Println("\n" + "═══════════════════════════════════════════════════════")
		fmt.Println("               ENVIRONMENT SCAN REPORT")
		fmt.Println("═══════════════════════════════════════════════════════")

		// Platform info
		fmt.Printf("\nPlatform: %s %s\n", result.Platform.OS, result.Platform.Arch)

		// Language environments
		fmt.Println("\nLanguage Environments:")
		if result.Python.Status != "missing" {
			fmt.Printf("  ✓ Python:   %s\n", result.Python.Version)
		} else {
			fmt.Printf("  ✗ Python:   NOT FOUND\n")
		}

		if result.Node.Status != "missing" {
			fmt.Printf("  ✓ Node.js:  %s\n", result.Node.Version)
		} else {
			fmt.Printf("  ✗ Node.js:  NOT FOUND\n")
		}

		if result.Rust.Status != "missing" {
			fmt.Printf("  ✓ Rust:     %s\n", result.Rust.Version)
		} else {
			fmt.Printf("  ⚠ Rust:     Not installed (optional)\n")
		}

		if result.Go.Status != "missing" {
			fmt.Printf("  ✓ Go:       %s\n", result.Go.Version)
		} else {
			fmt.Printf("  ⚠ Go:       Not installed (optional)\n")
		}

		// Status
		fmt.Println("\nStatus:")
		if result.IsHealthy {
			utils.Success("  ✓ Environment is healthy")
		} else {
			fmt.Printf("  ⚠ Found %d issue(s)\n", len(result.Issues))
		}

		if len(result.Issues) > 0 {
			fmt.Println("\n" + "───────────────────────────────────────────────────────")
			fmt.Println("Critical Issues:")
			fmt.Println("───────────────────────────────────────────────────────")
			for i, issue := range result.Issues {
				fmt.Printf("\n%d. [%s] %s\n", i+1, issue.Language, issue.Title)
				fmt.Printf("   Severity: %s\n", issue.Severity)
				fmt.Printf("   Message:  %s\n", issue.Message)
				if issue.Fix != "" {
					fmt.Printf("   Fix:      %s\n", issue.Fix)
				}
			}
		}

		if len(result.Warnings) > 0 {
			fmt.Println("\n" + "───────────────────────────────────────────────────────")
			fmt.Println("Warnings:")
			fmt.Println("───────────────────────────────────────────────────────")
			for i, warning := range result.Warnings {
				fmt.Printf("%d. %s\n", i+1, warning)
			}
		}

		fmt.Println("\n" + "═══════════════════════════════════════════════════════")
		fmt.Printf("Scan completed in %.2f seconds\n", result.ScanDuration)
		fmt.Println("═══════════════════════════════════════════════════════")

		return nil
	},
}

func init() {
	// TODO: Implement JSON and YAML output formats
}
