package consts

import "os"

const NodeIDLength = 8

var OrchestratorPort = os.Getenv("ORCHESTRATOR_PORT")
var CloudProviderEnv = os.Getenv("CLOUD_PROVIDER")

const (
	AWS string = "aws"
	GCP string = "gcp"
)
