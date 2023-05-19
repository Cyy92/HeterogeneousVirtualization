package virt

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/api/git"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/runtime"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	runtimeName    string
	handlerDir     string
	handlerName    string
	gitrepo        string
	existdir       string
	user           string
	token          string
	cloneRepoCmd   string
	fxNewGitInfo   *config.NewGitInfo
	fxExistGitInfo *config.ExistGitInfo
)

func init() {
	initCmd.Flags().StringVarP(&cloudconfig, "config", "f", "", "Path to cloud config file describing VM(s)")
	initCmd.Flags().StringVarP(&runtimeName, "runtime", "r", "", "Runtime(Language) to use")
	initCmd.Flags().StringVarP(&gitrepo, "reponame", "", "", "User's git repo name for upload binaries")
	initCmd.Flags().StringVarP(&existdir, "exist-workingdir", "", "", "User's exist git working directory")
	initCmd.Flags().StringVarP(&user, "user", "", "", "Git user name")
	initCmd.Flags().StringVarP(&token, "token", "", "", "Git PW")
	//initCmd.Flags().StringVarP(&domain, "dom", "d", "", "Set domain of VM")
	initCmd.MarkFlagRequired("gitrepo")
	initCmd.MarkFlagRequired("runtime")
	initCmd.MarkFlagRequired("user")
	initCmd.MarkFlagRequired("token")
}

var initCmd = &cobra.Command{
	Use: `init <VM_NAME>
  faas-cli virt init <VM_NAME> [-f <APPEND_EXISTING_YAML_FILE>]`,
	Short: "Prepare a FaaS Virtual Machine",
	Long: `
	The init command creates a new VM template. When user execute init command, config file and directory with VM name are created. Also, in directory with VM name, there is handler file and user can modify this file later. 
`,
	Example: `  faas-cli virt init ubuntu --runtime go --repo faas
  faas-cli virt init centos -f ./cloudconfig.yaml -r python3 --repo faas
  `,
	PreRunE: preRunInit,
	RunE:    runInit,
}

func validateVMName(VM_Name string) error {
	// Regex for RFC-1123 validation:
	// 	k8s.io/kubernetes/pkg/util/validation/validation.go
	var validDNS = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	matched := validDNS.MatchString(VM_Name)
	if matched {
		return nil
	}
	return fmt.Errorf(`VM name can only contain a-z, 0-9 and dashes`)
}

func preRunInit(cmd *cobra.Command, args []string) error {

	if len(args) < 1 || len(args) > 1 {
		return fmt.Errorf("please provide a name for VM")
	}

	VM_Name = args[0]
	if err := validateVMName(VM_Name); err != nil {
		return err
	}

	if cloudconfig == "" {
		cloudconfig = config.DefaultCloudConfig
	}

	if gitrepo == "" {
		return fmt.Errorf("please provide a git repo name")
	}

	if _, err := os.Stat(cloudconfig); err != nil {
		if os.IsNotExist(err) {
			fxVMs = config.NewVMs()
		} else {
			return err
		}
	} else {
		if err := parseCloudConfig(); err != nil {
			return err
		}
	}

	if !runtime.ExistFileOrDir(config.DefaultCloudInitDir) {
		if err := runtime.DownloadCloudInits(config.DefaultCloudInitDir, config.DefaultCloudInitRepo); err != nil {
			return err
		}
	}

	return nil
}

func runInit(cmd *cobra.Command, args []string) error {

	if _, err := os.Stat(VM_Name); err == nil {
		return fmt.Errorf("folder: %s already exists", VM_Name)
	}
	/*
		input := bufio.NewReader(os.Stdin)
		fmt.Print("GitHub Username: ")
		username, _ := input.ReadString('\n')

		fmt.Print("GitHub Password: ")
		bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
		password := string(bytePassword)
		fmt.Print("\n")

		client := git.NewClient(username, password)
	*/
	client := git.NewClient(user, token)
	if existdir == "" {
		git.CreateNewRepo(client, gitrepo)
	}

	r, err := runtime.GetRuntime(VM_Name, runtimeName, config.DefaultCloudInitDir)
	if err == nil {
		fmt.Printf("Folder: %s created with %s.\n", VM_Name, r.Handler.Name)
	} else {
		return fmt.Errorf("folder: could not create %s\n %s", VM_Name, err)
	}

	pubkey, err := ioutil.ReadFile("/root/.ssh/id_rsa.pub")
	if err != nil {
		return fmt.Errorf("Public key doesn't exist. please create public key by ssh-keygen -t rsa command\n %s", err)
	}

	fxVMs.Users = []config.User{
		{
			Name:           config.DefaultUser,
			Password:       config.DefaultPW,
			LockPW:         config.DefaultLockPW,
			Sudo:           config.DefaultSUDO,
			AuthorizedKeys: []string{string(pubkey)},
		},
	}

	cloneRepoCmd = "git clone https://github.com/" + strings.TrimSpace(user) + "/" + gitrepo + ".git" + " " + config.DefaultBinaryDir
	fxVMs.Cmds = []string{"apt install -y git", cloneRepoCmd}

	entries, err := os.ReadDir("./" + VM_Name + "/apps")
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() == true {
			executeBin := "cd /binaries/apps/" + e.Name() + " && ./executor &"
			fxVMs.Cmds = append(fxVMs.Cmds, executeBin)
		}
	}

	confYaml, err := yaml.Marshal(&fxVMs)
	if err != nil {
		return err
	}

	fmt.Printf("VM handler created in folder: %s\n", VM_Name+"/apps/src")
	fmt.Printf("Rewrite the VM handler code in %s folder\n", VM_Name+"/apps/src")

	if err := os.Mkdir("./"+VM_Name+"/config", os.ModePerm); err != nil {
		return err
	}

	initialConfErr := ioutil.WriteFile("./"+VM_Name+"/config/"+cloudconfig, []byte("#cloud-config\n\npackage_update: true\n\n"), 0600)
	if initialConfErr != nil {
		return fmt.Errorf("error writing config file %s\n", initialConfErr)
	}

	conf, err := os.OpenFile("./"+VM_Name+"/config/"+cloudconfig, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file with %s\n", err)
	}

	if _, err := conf.Write(confYaml); err != nil {
		return fmt.Errorf("cannot write file with %s\n", err)
	}

	if err := conf.Close(); err != nil {
		return fmt.Errorf("%s\n", err)
	}

	if existdir == "" {
		fxNewGitInfo = config.NewUser()
		fxNewGitInfo.Git = []config.NewGit{
			{
				Username: user,
				Token:    token,
				Repo:     gitrepo,
			},
		}

		git.Clone(VM_Name, strings.TrimSpace(user), gitrepo)

		userinfoYaml, err := yaml.Marshal(&fxNewGitInfo)
		if err != nil {
			return err
		}

		userinfoWriteErr := ioutil.WriteFile("./"+VM_Name+"/config/userinfo.yaml", userinfoYaml, 0600)
		if userinfoWriteErr != nil {
			return fmt.Errorf("error writing config file %s", userinfoWriteErr)
		}

		fmt.Printf("Config files written \n")
	} else {
		fxExistGitInfo = config.ExistUser()
		fxExistGitInfo.Git = []config.ExistGit{
			{
				Username:   user,
				Token:      token,
				Repo:       gitrepo,
				WorkingDir: existdir,
			},
		}

		userinfoYaml, err := yaml.Marshal(&fxExistGitInfo)
		if err != nil {
			return err
		}

		userinfoWriteErr := ioutil.WriteFile("./"+VM_Name+"/config/userinfo.yaml", userinfoYaml, 0600)
		if userinfoWriteErr != nil {
			return fmt.Errorf("error writing config file %s", userinfoWriteErr)
		}

		fmt.Printf("Config files written \n")
	}

	return nil
}
