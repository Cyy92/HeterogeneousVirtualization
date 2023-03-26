package virt

import (
	"fmt"
	"os"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/builder"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/log"
	"github.com/spf13/cobra"
)

var (
	buildVerbose bool
)

func init() {
	buildCmd.Flags().BoolVarP(&buildVerbose, "buildverbose", "v", false, "Print function build log")
}

var buildCmd = &cobra.Command{
	Use:   `build`,
	Short: "Build FaaS function",
	Long: `
	Build FaaS function for executing in VM by binaries via the supplied YAML config 
	`,
	Example: `
	faas-cli virt build 
	faas-cli virt build -v
	`,
	PreRunE: preRunBuild,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runBuild(); err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func preRunBuild(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat("./apps/src"); err != nil {
		return fmt.Errorf("Unable to build with %s\n", err)
	}

	return nil
}

func buildExecutor(verbose bool) error {
	result, err := builder.BuildGoExecutor(verbose)
	if err != nil {
		log.Print(result)
		return err
	}

	return nil
}

func buildHandler(verbose bool) error {
	result, err := builder.BuildGoHandler(verbose)
	if err != nil {
		log.Print(result)
		return err
	}

	return nil
}

func runBuild() error {
	log.Info("Building executor ...\n")
	if err := buildExecutor(buildVerbose); err != nil {
		return err
	}

	log.Info("Building handler ...\n")
	if err := buildHandler(buildVerbose); err != nil {
		return err
	}

	return nil
}
