package template

import (
	"context"
	"fmt"
	"log"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"

	"github.com/e2b-dev/infra/packages/shared/pkg/consts"
	"github.com/e2b-dev/infra/packages/shared/pkg/telemetry"
	"github.com/e2b-dev/infra/packages/template-manager/internal/provider"
	"github.com/gogo/status"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
)


var awsProvider = provider.NewProvider(consts.AWSRegion)


func GetDockerImageURL(templateID string) string {
	// DockerImagesURL is the URL to the docker images in the artifact registry
	
	
	if consts.CloudProviderEnv == consts.GCP {
		return fmt.Sprintf("projects/%s/locations/%s/repositories/%s/packages/%s", consts.GCPProject, consts.GCPRegion, consts.DockerRegistry, templateID)
	}
	if consts.CloudProviderEnv == consts.AWS {
		accountID := consts.AWSAccountID
		if accountID == "" {
			_accountID, err := awsProvider.GetAWSAccountID()
			if err != nil {
				errMsg := fmt.Errorf("error getting AWS account ID: %w", err)
				panic(errMsg)
			}
			accountID = _accountID
		}
		return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:latest", accountID, consts.AWSRegion, templateID)
	}
	fmt.Errorf("unsupported cloud provider: %s", consts.CloudProviderEnv)
	return ""
}

func Delete(
	ctx context.Context,
	tracer trace.Tracer,
	artifactRegistry *artifactregistry.Client,
	ecrClient *ecr.Client,
	templateStorage *Storage,
	buildId string,
) error {
	childCtx, childSpan := tracer.Start(ctx, "delete-template")
	defer childSpan.End()

	err := templateStorage.Remove(ctx, buildId)
	if err != nil {
		return fmt.Errorf("error when deleting template objects: %w", err)
	}

	switch consts.CloudProviderEnv {
	case consts.GCP:
		return deleteFromGCP(childCtx, artifactRegistry, buildId)
	case consts.AWS:
		return deleteFromAWS(childCtx, buildId)
	default:
		return fmt.Errorf("unsupported cloud provider: %s", consts.CloudProviderEnv)
	}
}

func deleteFromGCP(ctx context.Context, artifactRegistry *artifactregistry.Client, buildId string) error {
	op, artifactRegistryDeleteErr := artifactRegistry.DeletePackage(ctx, &artifactregistrypb.DeletePackageRequest{
		Name: GetDockerImageURL(buildId),
	})

	if artifactRegistryDeleteErr != nil {
		if status.Code(artifactRegistryDeleteErr) == codes.NotFound {
			log.Printf("template image not found in GCP registry, skipping deletion: %v", artifactRegistryDeleteErr)
			telemetry.ReportEvent(ctx, fmt.Sprintf("template image not found in GCP registry, skipping deletion: %v", artifactRegistryDeleteErr))
			return nil
		}
		errMsg := fmt.Errorf("error when deleting template image from GCP registry: %w", artifactRegistryDeleteErr)
		telemetry.ReportCriticalError(ctx, errMsg)
		return errMsg
	}

	telemetry.ReportEvent(ctx, "started deleting template image from GCP registry")

	if waitErr := op.Wait(ctx); waitErr != nil {
		errMsg := fmt.Errorf("error when waiting for template image deletion from GCP registry: %w", waitErr)
		telemetry.ReportCriticalError(ctx, errMsg)
		return errMsg
	}

	telemetry.ReportEvent(ctx, "deleted template image from GCP registry")
	return nil
}

func deleteFromAWS(ctx context.Context, buildId string) error {
	// TODO: Implement AWS ECR deletion logic
	// This would typically involve:
	// 1. Creating an AWS ECR client
	// 2. Using BatchDeleteImage or DeleteRepository API
	fmt.Println(GetDockerImageURL(buildId))
	telemetry.ReportEvent(ctx, "AWS deletion not implemented yet")
	return fmt.Errorf("AWS deletion not implemented yet")
}
