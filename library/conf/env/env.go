package env

import "os"

// deploy env.
const (
	DeployEnvDev    = "dev"
	DeployEnvFat1   = "fat1"
	DeployEnvUat    = "uat"
	DeployEnvOffice = "office"
	DeployEnvAvalon = "avalon" // deprecated, no avalon forever
	DeployEnvPre    = "pre"
	DeployEnvProd   = "prod"
)

// env default value.
const (
	// env
	_region    = "sh"
	_zone      = "sh001"
	_deployEnv = "dev"
)

// env configuration.
var (
	// Region avaliable region where app at.
	Region string
	// Zone avaliable zone where app at.
	Zone string
	// Hostname machine hostname.
	Hostname string
	// DeployEnv deploy env where app at.
	DeployEnv string
	// Color is the identification of different experimental group in one caster cluster.
	Color string
	// AppID is global unique application id, register by service tree.
	// such as main.arch.disocvery.
	AppID string
	// DiscoveryAppID global unique application id for disocvery, register by service tree.
	// such as main.arch.disocvery.
	DiscoveryAppID string
	// DiscoveryZone is discovery zone.
	DiscoveryZone string
	// DiscoveryHost is discovery host.
	DiscoveryHost string
	// IP FIXME(haoguanwei) #240
	IP = os.Getenv("POD_IP")
)

// app default value.
const (
	_httpPort  = "8000"
	_gorpcPort = "8099"
	_grpcPort  = "9000"
)

// app configraution.
var (
	// HTTPPort app listen http port.
	HTTPPort string
	// GORPCPort app listen gorpc port.
	GORPCPort string
	// GRPCPort app listen grpc port.
	GRPCPort string
)
