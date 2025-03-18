package consts

import (
	"os"
)

var (
	AWSRegion          = os.Getenv("AWS_REGION")
	AWSAccountID       = os.Getenv("AWS_ACCOUNT_ID")
	AWSECRRegistry     = os.Getenv("AWS_ECR_REGISTRY")
	AWSAccessKeyID     = os.Getenv("AWS_ACCESS_KEY_ID")
	AWSSecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
)
