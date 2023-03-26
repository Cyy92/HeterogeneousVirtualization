package virt

import (
	"os"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"

	"github.com/spf13/cobra"
)

var (
	cloudconfig string
	configFile  string
	fxVMs       *config.VirtualMachines
	VM_Name     string
)

var VirtCmd = &cobra.Command{
	Use: "virt SUBCOMMAND",
	//Aliases: []string{"fn"},
	Short: "virt specific operations",
	Long: `
	virt command allows user to init, list, deploy, delete VMs running on FaaS
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	VirtCmd.AddCommand(initCmd)
	VirtCmd.AddCommand(buildCmd)
	VirtCmd.AddCommand(deployVMCmd)
	//FunctionCmd.AddCommand(deleteCmd)
	//FunctionCmd.AddCommand(listCmd)
	//FunctionCmd.AddCommand(callCmd)
	//FunctionCmd.AddCommand(infoCmd)
	//FunctionCmd.AddCommand(logCmd)

	//VirtCmd.PersistentFlags().StringVarP(&gitrepo, "repo", "", "", "User's git repo for upload binaries")
	//VirtCmd.PersistentFlags().StringVarP(&domain, "dom", "d", "", "Set domain of VM")
}

func parseCloudConfig() error {

	var err error
	if _, err := os.Stat(cloudconfig); err != nil {
		if os.IsNotExist(err) {
		}
		return err
	}
	if fxVMs, err = config.ParseCloudConfig(cloudconfig); err != nil {
		return err
	}

	return nil
}
