package function

import (
	"fmt"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/grpc"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"
	"github.com/spf13/cobra"
)

func init() {
}

var logCmd = &cobra.Command{
	Use:   `log <FUNCTION_NAME>`,
	Short: "Display FaaS function logs",
	Long: `
	Display FaaS function logs
`,
	Example: `  faas-cli function log resizeImg
`,
	PreRunE: preRunLog,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runLog(); err != nil {
			return err
		}
		return nil
	},
}

func preRunLog(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a name for the function")
	}

	functionName = args[0]

	gateway = config.GetFxGatewayURL(gateway, "")
	return nil
}

func runLog() error {

	fnLog, err := grpc.GetLog(functionName, gateway)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", string(fnLog))

	return nil
}
