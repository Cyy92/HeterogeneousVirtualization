package virt

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/git"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/grpc"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/log"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"

	"github.com/spf13/cobra"
)

var (
	gateway string
	memory  string
	domain  string
)

func init() {
	deployVMCmd.Flags().StringVarP(&cloudconfig, "config", "f", "", "Path to cloud config file describing VM(s)")
	deployVMCmd.Flags().StringVarP(&memory, "mem", "", "", "Memory usage for VMs")
	deployVMCmd.Flags().StringVarP(&domain, "dom", "d", "", "Set domain name of VMs, such as ubuntu or centos")
	deployVMCmd.Flags().StringVarP(&gateway, "gateway", "g", "", "Set gateway URL")
	deployVMCmd.MarkFlagRequired("config")
	deployVMCmd.MarkFlagRequired("gateway")
	deployVMCmd.MarkFlagRequired("dom")
}

var deployVMCmd = &cobra.Command{
	Use:   `deploy -f <YAML_CLOUD_CONIFIG_FILE>`,
	Short: "Deploy FaaS VMs",
	Long: `
	Deploy FaaS VM containers via the supplied cloud config using the "-f" flag.
	`,
	Example: `  
	faas-cli virt deploy ubuntu -f cloudconfig.yaml --mem 4096Mi --dom ubuntu --gateway 10.0.0.180:31113
	faas-cli virt deploy ubuntu -f cloudconfig.yaml -g 10.0.0.180:31113 --dom ubuntu
        `,
	PreRunE: preRunDeployVM,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDeployVM(); err != nil {
			fmt.Println(err.Error())
		}

		return
	},
}

func preRunDeployVM(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args) > 1 {
		return fmt.Errorf("please provide a name for VM")
	}

	VM_Name = args[0]

	if cloudconfig == "" {
		ce := fmt.Sprintf("please provide a '-f' flag for vm creation\n")
		return errors.New(ce)
	} else {
		if err := parseCloudConfig(); err != nil {
			return err
		}
	}
	if gateway == "" {
		ge := fmt.Sprintf("please provide a '-g' flag for setting gateway address\n")
		return errors.New(ge)
	}

	if domain == "" {
		return fmt.Errorf("please provide a domain name for VM")
	}

	if memory == "" {
		memory = config.DefaultVMMemory
	}

	return nil
}

func deploy(gw string, vmName string, domain string, userdata string) error {
	resource := &config.FunctionResources{
		Memory: memory,
	}

	deployVMConfig := grpc.DeployVMConfig{
		FxGateway: gw,
		VMName:    vmName,
		Domain:    domain,
		UserData:  userdata,
		Requests:  resource,
	}

	if err := grpc.DeployVM(deployVMConfig); err != nil {
		return err
	}
	return nil
}

func runDeployVM() error {
	git.Push()

	f, _ := os.Open(cloudconfig)

	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	encoded := base64.StdEncoding.EncodeToString(content)

	log.Info("Deploying: %s ...\n", VM_Name)

	//DEPLOY
	if err := deploy(gateway, VM_Name, domain, encoded); err != nil {
		return err
	}

	return nil
}
