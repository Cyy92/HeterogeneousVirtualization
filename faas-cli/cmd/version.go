package cmd

import (
	"fmt"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the FaaS CLI version information",
	Long: `
	display the FaaS CLI version information
	`,
	Example: `faas-cli version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", config.FaaSCliVersion)
	},
}
