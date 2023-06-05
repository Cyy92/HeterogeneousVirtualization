package virt

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/grpc"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/log"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"
	"github.com/spf13/cobra"
)

func init() {
}

var callCmd = &cobra.Command{
	Use:     `call <VM_NAME>`,
	Aliases: []string{"invoke"},
	Short:   "Call FaaS VMs",
	Long: `
	Call FaaS vm and reads from STDIN for handler(user defined function)'s input(bytes)
	`,
	Example: `  faas-cli virt call ubuntu
	cat "sample.png" | faas-cli virt call -g localhost:31113 ubuntu
        `,
	PreRunE: preRunCall,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("please provide a name for the function")
		}

		VM_Name = args[0]

		if err := runCall(); err != nil {
			return err
		}
		return nil
	},
}

func preRunCall(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		log.Fatal("Invalid function name. please describe name of function correctly\n")

	}

	VM_Name = args[0]

	gateway = config.GetFxGatewayURL(gateway, "")
	return nil
}

func runCall() error {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintf(os.Stderr, "Reading from STDIN - hit (Control + D) to stop.\n")
	}

	functionInput, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("unable to read standard input: %s", err.Error())
	}

	resp, err := grpc.InvokeVM(VM_Name, gateway, functionInput)
	if err != nil {
		return err
	}

	if resp != "" {
		os.Stdout.WriteString(resp)
	}

	return nil
}
