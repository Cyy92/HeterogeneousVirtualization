package config

const (
	FaaSCliVersion   = "1.0"
	GatewayEnvVarKey = "OPENFX_URL"

	DefaultProviderName = "faas"
	DefaultConfigFile   = "config.yaml"
	DefaultCloudConfig  = "cloudconfig.yaml"
	DefaultRegistry     = "10.0.1.150:5000/cyy"
	DefaultGatewayURL   = "localhost:31113"
	//FIXME
	DefaultRuntimeRepo = "https://github.com/keti-openfx/OpenFx-runtime.git"
	DefaultRuntimeDir  = "./runtime"
	DefaultCPU         = "50m"
	DefaultMemory      = "50Mi"
	DefaultGPU         = ""

	// Virt Vm Configuration
	DefaultCloudInitRepo = "https://github.com/Cyy92/Cloud-Init.git"
	DefaultCloudInitDir  = "./cloudinit"
	DefaultBinaryDir     = "/binaries"
	DefaultUser          = "root"
	DefaultPW            = "ketilinux"
	DefaultLockPW        = false
	DefaultSUDO          = "ALL=(ALL) NOPASSWD:ALL"
	DefaultVMMemory      = "2048M"
)

var (
	DefaultConstraints  = []string{"nodetype=cpunode"}
	DefaultOAuth2Server = "http://10.0.2.101:9096"
)

type VirtualMachines struct {
	Users []User   `yaml:"users,omitempty"`
	Cmds  []string `yaml:"runcmd,omitempty"`
}

type User struct {
	Name           string   `yaml:"name,omitempty"`
	Password       string   `yaml:"password,omitempty"`
	LockPW         bool     `yaml:"lock-passwd"`
	Sudo           string   `yaml:"sudo,omitempty"`
	AuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
}

type NewGitInfo struct {
	Git []NewGit `yaml:"git_info,omitempty"`
}

type NewGit struct {
	Username string `yaml:"git_user,omitempty"`
	Token    string `yaml:"token,omitempty"`
	Repo     string `yaml:"repo_name,omitempty"`
}

type ExistGitInfo struct {
	Git []ExistGit `yaml:"git_info,omitempty"`
}

type ExistGit struct {
	Username   string `yaml:"username,omitempty"`
	Token      string `yaml:"token,omitempty"`
	Repo       string `yaml:"repo_name,omitempty"`
	WorkingDir string `yaml:"exist_workingdir,omitempty"`
}

type Services struct {
	Functions map[string]Function `yaml:"functions,omitempty"`
	FaaS      FaaS                `yaml:"faas,omitempty"`
}

type FaaS struct {
	FxGatewayURL string `yaml:"gateway"`
}

type Handler struct {
	// Local directory to use for function
	Dir string `yaml:"dir",omitempty`
	// Local file to use for function
	File string `yaml:"file",omitempty`
	// function name to use for function
	Name string `yaml:"name"`
}

// Function as deployed or built on OpenFx
type Function struct {
	// Name of deployed function
	Name    string `yaml:"-"`
	Runtime string `yaml:"runtime"`

	Description string `yaml:"desc",omitempty`
	Maintainer  string `yaml:"maintainer",omitempty`

	// Handler to use for function
	Handler Handler `yaml:"handler"`

	// Doker private registry
	RegistryURL string `yaml:"docker_registry"`

	// Image Docker image name
	Image string `yaml:"image"`

	// Docker registry Authorization
	RegistryAuth string `yaml:"registry_auth,omitempty"`

	Environment map[string]string `yaml:"environment,omitempty"`

	// Secrets list of secrets to be made available to function
	Secrets []string `yaml:"secrets,omitempty"`

	//SkipBuild bool `yaml:"skip_build,omitempty"`

	Constraints *[]string `yaml:"constraints,omitempty"`

	// EnvironmentFile is a list of files to import and override environmental variables.
	// These are overriden in order.
	EnvironmentFile []string `yaml:"environment_file,omitempty"`

	Labels *map[string]string `yaml:"labels,omitempty"`

	// Limits for function
	Limits *FunctionResources `yaml:"limits,omitempty"`

	// Requests of resources requested by function
	Requests *FunctionResources `yaml:"requests,omitempty"`

	// BuildOptions to determine native packages
	BuildOptions []string `yaml:"build_options,omitempty"`

	// BuildOptions to determine native packages
	BuildArgs []string `yaml:"build_args,omitempty"`
}

// FunctionResources Memory and CPU, GPU
type FunctionResources struct {
	Memory string `yaml:"memory"`
	CPU    string `yaml:"cpu"`
	GPU    string `yaml:"gpu"`
}

// EnvironmentFile represents external file for environment data
type EnvironmentFile struct {
	Environment map[string]string `yaml:"environment"`
}
