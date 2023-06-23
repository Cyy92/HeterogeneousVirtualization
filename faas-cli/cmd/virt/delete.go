package virt

import (
	"errors"
	"fmt"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/grpc"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/log"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringVarP(&gateway, "gateway", "g", "", "Set gateway URL")
	deleteCmd.MarkFlagRequired("gateway")
}

var deleteCmd = &cobra.Command{
	Use:     `delete <VM_NAME>`,
	Aliases: []string{"remove", "rm"},
	Short:   "Delete FaaS VMs",
	Long: `
	Delete FaaS vm using
the function name(which may contain multiple function definitions)
`,
	Example: `faas-cli virt delete ubuntu --gateway 10.0.2.101:31113`,
	PreRunE: preRunDelete,
	Run: func(cmd *cobra.Command, args []string) {

		if err := runDelete(); err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func preRunDelete(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args) > 1 {
		return fmt.Errorf("please provide a name for VM")
	}

	VM_Name = args[0]

	if gateway == "" {
		ge := fmt.Sprintf("please provide a '-g' flag for setting gateway address\n")
		return errors.New(ge)
	}

	return nil
}

func runDelete() error {
	log.Info("Deleting: %s ...\n", VM_Name)
	if err := grpc.DeleteVM(gateway, VM_Name); err != nil {
		return err
	}

	log.Print("Deleted: %s\n", VM_Name)

	return nil
}
