package function

import (
	"errors"
	"fmt"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/grpc"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/log"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
	deleteCmd.Flags().StringVarP(&gateway, "gateway", "g", "localhost:31113", "Gateway URL to store in YAML config file")
}

var deleteCmd = &cobra.Command{
	Use:     `delete -f <YAML_CONIFIG_FILE>`,
	Aliases: []string{"remove", "rm"},
	Short:   "Delete FaaS functions",
	Long: `
	Delete FaaS function via the supplied YAML config using
the "-f" flag or the function name(which may contain multiple function definitions)
`,
	Example: `  faas-cli function delete -f config.yml
		    faas-cli function delete echo
                  `,
	PreRunE: preRunDelete,
	Run: func(cmd *cobra.Command, args []string) {

		if err := runDelete(); err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func preRunDelete(cmd *cobra.Command, args []string) error {
	fxServices = config.NewServices()

	var configURL string
	if cmd.Flag("config").Value.String() != "" {
		if err := parseConfigFile(); err != nil {
			return err
		}
		configURL = fxServices.FaaS.FxGatewayURL
	}

	gateway = config.GetFxGatewayURL(gateway, configURL)

	if len(args) > 0 {
		fxServices.Functions = make(map[string]config.Function, 0)
		fxServices.Functions[args[0]] = config.Function{}
	}

	return nil
}

func runDelete() error {
	if len(fxServices.Functions) <= 0 {
		return errors.New("")
	}

	for name, function := range fxServices.Functions {
		function.Name = name
		if err := grpc.Delete(gateway, function.Name); err != nil {
			return err
		}

		log.Print("Deleted: %s\n", function.Name)
	}

	return nil
}
