package constants

import (
	"fmt"
	"strings"

	"github.com/e2b-dev/infra/packages/shared/pkg/consts"
)

func CheckRequired() error {
	var missing []string

	if consts.CloudProviderEnv == "" {
		missing = append(missing, "CLOUD_PROVIDER")
	}

	if consts.CloudProviderEnv == consts.AWS {
		if consts.AWSRegion == "" {
			missing = append(missing, "AWS_REGION")
		}
		if consts.AWSAccountID == "" {
			missing = append(missing, "AWS_ACCOUNT_ID")
		}
		if consts.AWSECRRegistry == "" {
			missing = append(missing, "AWS_ECR_REGISTRY")
		}

		if consts.AWSAccessKeyID == "" {
			missing = append(missing, "AWS_ACCESS_KEY_ID")
		}

		if consts.AWSSecretAccessKey == "" {
			missing = append(missing, "AWS_SECRET_ACCESS_KEY")
		}

	} else if consts.CloudProviderEnv == consts.GCP {
		if consts.GCPProject == "" {
			missing = append(missing, "GCP_PROJECT_ID")
		}

		if consts.DockerRegistry == "" {
			missing = append(missing, "GCP_DOCKER_REPOSITORY_NAME")
		}

		if consts.GoogleServiceAccountSecret == "" {
			missing = append(missing, "GOOGLE_SERVICE_ACCOUNT_BASE64")
		}

		if consts.GCPRegion == "" {
			missing = append(missing, "GCP_REGION")
		}

		if len(missing) > 0 {
			return fmt.Errorf("missing environment variables: %s", strings.Join(missing, ", "))
		}
	}

	return nil
}
