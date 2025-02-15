package function

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/grpc"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/builder"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/log"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var (
	replace       bool
	update        bool
	deployVerbose bool
	registry      string
	token         string
	minreplicas   int32
	maxreplicas   int32
)

func init() {
	deployCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
	deployCmd.Flags().StringVarP(&registry, "registry", "", "", "Docker private registry url")
	deployCmd.Flags().StringVarP(&token, "token", "", "", "Access token for deploying function(s)")
	deployCmd.Flags().BoolVar(&replace, "replace", false, "Remove and re-create existing function(s)")
	deployCmd.Flags().BoolVar(&update, "update", true, "Perform rolling update on existing function(s)")
	deployCmd.Flags().BoolVarP(&deployVerbose, "deployverbose", "v", false, "Print function build log")
	deployCmd.Flags().Int32Var(&minreplicas, "min", 1, "Minimum Replicas for Function")
	deployCmd.Flags().Int32Var(&maxreplicas, "max", 1, "Maximum Replicas for Function")
	deployCmd.MarkFlagRequired("config")
	deployCmd.MarkFlagRequired("token")
}

var deployCmd = &cobra.Command{
	Use:   `deploy -f <YAML_CONIFIG_FILE>`,
	Short: "Deploy OpenFx functions",
	Long: `
	Push OpenFx function Image & Deploy OpenFx function containers via the supplied YAML config using the "-f" flag. Also write docker private registry using the "--registry" flag to push docker image into registry.
	`,
	Example: `  
	openfx-cli function deploy -f config.yml
  	openfx-cli function deploy -f ./config.yml --replace=false --update=true
	openfx-cli function deploy -f config.yml -v
	openfx-cli function deploy -f config.yml --registry 127.0.0.1:5000
	openfx-cli function deploy -f config.yml -g 10.0.0.180:31113
	openfx-cli function deploy -f config.yml --min 1 --max 5
        `,
	PreRunE: preRunDeploy,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDeploy(); err != nil {
			fmt.Println(err.Error())
		}

		return
	},
}

func preRunDeploy(cmd *cobra.Command, args []string) error {
	params := url.Values{}
	params.Add("access_token", token)

	apiUrl := fmt.Sprintf(config.DefaultOAuth2Server+"/verify?%s", params.Encode())
	resp, err := http.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("Token verify failed: %s", err)
	}
	defer resp.Body.Close()

	if update && replace {
		return errors.New(`one of "--update" flag or "--replace" flag must be false\n`)
	}

	var configURL string
	if configFile == "" {
		e := fmt.Sprintf("please provide a '-f' flag for function creation\n")
		return errors.New(e)
	} else {
		if err := parseConfigFile(); err != nil {
			return err
		}
		configURL = fxServices.FaaS.FxGatewayURL
	}
	gateway = config.GetFxGatewayURL(gateway, configURL)

	return nil
}

func deploy(gw string, function config.Function, update, replace bool, minreplicas int32, maxreplicas int32) error {

	//function.Secrets
	//sendRegistryAuth
	//EnvVar
	fileEnvironment, err := readFiles(function.EnvironmentFile)
	if err != nil {
		return err
	}
	allEnvironment := mergeMap(function.Environment, fileEnvironment)

	//Labels
	labelMap := map[string]string{}
	if function.Labels != nil {
		labelMap = *function.Labels
	}

	//Annotations
	AnnoMap := map[string]string{}
	if function.Maintainer != "" {
		AnnoMap["maintainer"] = function.Maintainer
	}
	if function.Description != "" {
		AnnoMap["desc"] = function.Description
	}

	// Get FxProcess to use from the ?
	deployConfig := grpc.DeployConfig{
		FxGateway:    gw,
		FunctionName: function.Name,
		Image:        function.Image,
		EnvVars:      allEnvironment,
		Labels:       labelMap,
		Annotations:  AnnoMap,
		Constraints:  function.Constraints,
		Secrets:      append(function.Secrets, "regcred"),
		Limits:       function.Limits,
		Requests:     function.Requests,

		MinReplicas: minreplicas,
		MaxReplicas: maxreplicas,
		Token:       token,
		Update:      update,
		Replace:     replace,
	}
	if err := grpc.Deploy(deployConfig); err != nil {
		return err
	}
	return nil
}

func runDeploy() error {
	if len(fxServices.Functions) <= 0 {
		return errors.New("")
	}

	for name, function := range fxServices.Functions {

		function.Name = name

		log.Info("Pushing: %s, Image: %s in Registry: %s ...\n", function.Name, function.Image, function.RegistryURL)
		if deployVerbose {
			err := builder.ExecCommandPipe("./", []string{"docker", "push", function.Image}, os.Stdout, os.Stderr)
			if err != nil {
				return err
			}
		} else {
			_, err := builder.ExecCommand("./", []string{"docker", "push", function.Image})
			if err != nil {
				return err
			}
		}

		log.Info("Deploying: %s ...\n", function.Name)

		//DEPLOY
		if err := deploy(gateway, function, update, replace, minreplicas, maxreplicas); err != nil {
			return err
		}

		log.Info("http trigger url: http://%s/function/%s \n", gateway, function.Name)
	}

	return nil
}

func readFiles(files []string) (map[string]string, error) {
	envs := make(map[string]string)

	for _, file := range files {
		bytesOut, readErr := ioutil.ReadFile(file)
		if readErr != nil {
			return nil, readErr
		}

		envFile := config.EnvironmentFile{}
		unmarshalErr := yaml.Unmarshal(bytesOut, &envFile)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		for k, v := range envFile.Environment {
			envs[k] = v
		}
	}
	return envs, nil
}

func mergeMap(i map[string]string, j map[string]string) map[string]string {
	merged := make(map[string]string)

	for k, v := range i {
		merged[k] = v
	}
	for k, v := range j {
		merged[k] = v
	}
	return merged
}
