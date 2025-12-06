package cmd

import (
	"envfix/internal/utils"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	logFile string
)

var rootCmd = &cobra.Command{
	Use:   "envfix",
	Short: "Environment and dependency doctor",
	Long: `envfix is a cross-language environment and dependency doctor that diagnoses,
repairs, standardizes, and cleans developer environments across Python,
Node.js, Rust, Go, and system-level toolchains.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logger if log file specified
		if logFile != "" {
			if err := utils.InitLogger(logFile); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		utils.CloseLogger()
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Log file path")

	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(repairCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(explainCmd)
}
