package template

import (
	"context"
	"fmt"
	"log"
	"errors"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	"cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"

	"github.com/e2b-dev/infra/packages/shared/pkg/consts"
	"github.com/e2b-dev/infra/packages/shared/pkg/telemetry"
	"github.com/e2b-dev/infra/packages/template-manager/internal/provider"
	"github.com/gogo/status"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)


var awsProvider = provider.NewProvider(consts.AWSRegion)


func GetDockerImageURL(templateID string, buildId string) string {
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
		return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com/%s:%s", accountID, consts.AWSRegion, templateID, buildId)
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
	templateID string,
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
		return deleteFromGCP(childCtx, artifactRegistry, templateID, buildId)
	case consts.AWS:
		return deleteFromAWS(childCtx,ecrClient, templateID, buildId)
	default:
		return fmt.Errorf("unsupported cloud provider: %s", consts.CloudProviderEnv)
	}
}

func deleteFromGCP(ctx context.Context, artifactRegistry *artifactregistry.Client, templateID string, buildId string) error {
	op, artifactRegistryDeleteErr := artifactRegistry.DeletePackage(ctx, &artifactregistrypb.DeletePackageRequest{
		Name: GetDockerImageURL(templateID, buildId),
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

func deleteFromAWS(ctx context.Context, ecrClient *ecr.Client, templateID string, buildId string) error {
	fmt.Println("Deleting template image from AWS ECR registry",templateID, buildId)
	imageIdentifier := &ecr.BatchDeleteImageInput{
		RepositoryName: &templateID,
		ImageIds: []types.ImageIdentifier{
			{
				ImageTag: &buildId,
			},
		},
	}

	_, err := ecrClient.BatchDeleteImage(ctx, imageIdentifier)
	if err != nil {
		var notFoundErr *types.RepositoryNotFoundException
		if errors.As(err, &notFoundErr) {
			log.Printf("template image not found in AWS ECR registry, skipping deletion: %v", err)
			telemetry.ReportEvent(ctx, fmt.Sprintf("template image not found in AWS ECR registry, skipping deletion: %v", err))
			return nil
		}
		errMsg := fmt.Errorf("error when deleting template image from AWS ECR registry: %w", err)
		telemetry.ReportCriticalError(ctx, errMsg)
		return errMsg
	}

	telemetry.ReportEvent(ctx, "deleted template image from AWS ECR registry")
	return nil
}
