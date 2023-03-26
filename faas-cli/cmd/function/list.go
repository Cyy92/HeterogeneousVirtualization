package function

import (
	"fmt"
	"strings"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/grpc"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"
	"github.com/spf13/cobra"
)

func init() {
	listCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
}

var listCmd = &cobra.Command{
	Use:   `list -f <YAML_CONIFIG_FILE>`,
	Short: "Lists FaaS functions",
	Long: `
	Lists FaaS function
`,
	Example: `  faas-cli function list -f config.yml
	faas-cli funtion list -g localhost:31113
                  `,
	PreRunE: preRunList,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runList(); err != nil {
			fmt.Println(err.Error())
		}
		return
	},
}

func preRunList(cmd *cobra.Command, args []string) error {
	var configURL string
	if cmd.Flag("config").Value.String() != "" {
		if err := parseConfigFile(); err != nil {
			return err
		}
		configURL = fxServices.FaaS.FxGatewayURL
	}
	gateway = config.GetFxGatewayURL(gateway, configURL)

	return nil
}

func runList() error {
	fnList, err := grpc.List(gateway)
	if err != nil {
		return err
	}

	fmt.Printf("%-15s\t%-20s\t%-15s\t%-10s\t%-10s\t%-10s\t%-40s\n", "Function", "Image", "Maintainer", "Invocations", "Replicas", "Status", "Description")
	for _, fn := range fnList.Functions {

		if fxServices != nil {
			if _, ok := fxServices.Functions[fn.Name]; !ok {
				continue
			}
		}
		var fnImage string
		if fn.Image != config.DefaultRegistry+"/"+fn.Name {
			fnImage = strings.Replace(fn.Image, strings.Split(fn.Image, "/")[0], "$(repo)", 1)
		} else {
			fnImage = strings.Replace(fn.Image, config.DefaultRegistry, "$(repo)", 1)
		}

		if len(fnImage) > 30 {
			fnImage = fnImage[0:28] + ".."
		}

		var fnMaintainer, fnDesc, fnStatus string
		if v, ok := fn.Annotations["maintainer"]; ok {
			fnMaintainer = v
		}

		if v, ok := fn.Annotations["desc"]; ok {
			fnDesc = v
			if len(fnDesc) > 40 {
				fnDesc = fnDesc[0:38] + ".."
			}
		}

		if fn.AvailableReplicas == 0 {
			fnStatus = "Not Ready"
		} else {
			fnStatus = "Ready"
		}

		fmt.Printf("%-15s\t%-20s\t%-15s\t%-10d\t%-10d\t%-10s\t%-40s\n", fn.Name, fnImage, fnMaintainer, fn.InvocationCount, fn.Replicas, fnStatus, fnDesc)

	}

	return nil
}
