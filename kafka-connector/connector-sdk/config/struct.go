package config

const (
	FxCliVersion     = "0.5.0"
	GatewayEnvVarKey = "OPENFX_URL"

	DefaultProviderName = "openfx"
	DefaultConfigFile   = "config.yaml"
	DefaultRegistry     = "keti.asuscomm.com:5000"
	DefaultGatewayURL   = "10.0.0.202:31113"
	//FIXME
	DefaultRuntimeRepo = "https://github.com/keti-openfx/OpenFx-runtime.git"
	DefaultRuntimeDir  = "./runtime"
	DefaultCPU         = "50m"
	DefaultMemory      = "50Mi"
)

type Services struct {
	Functions map[string]Function `yaml:"functions,omitempty"`
	Openfx    Openfx              `yaml:"openfx,omitempty"`
}

type Openfx struct {
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
