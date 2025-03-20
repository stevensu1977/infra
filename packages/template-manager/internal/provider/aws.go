package provider

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AWSProvider struct {
	ECRClient *ecr.Client
	S3Client  *s3.Client
	config    aws.Config
	accountID string
}

func NewProvider(region string) *AWSProvider {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region), // Replace with your AWS region
		config.WithSharedConfigProfile(""),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS configuration: %v", err)
		return nil
	}

	// Create ECR client
	ecrClient := ecr.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)

	return &AWSProvider{
		ECRClient: ecrClient,
		S3Client:  s3Client,
		config:    cfg,
	}

}

func (p *AWSProvider) GetAWSAccountID() (string, error) {
	stsClient := sts.NewFromConfig(p.config)
	result, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return *result.Account, nil
}

func (p *AWSProvider) GetECRPassword() (string, error) {
	output, err := p.ECRClient.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		log.Fatalf("Failed to get ECR authorization token: %v", err)
	}

	// Check if authorization data was received
	if len(output.AuthorizationData) == 0 {
		log.Fatal("No ECR authorization data received")
	}

	authToken := *output.AuthorizationData[0].AuthorizationToken
	decodedToken, err := base64.StdEncoding.DecodeString(authToken)
	if err != nil {
		return "", fmt.Errorf("failed to decode auth token: %v", err)
	}

	parts := strings.SplitN(string(decodedToken), ":", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid auth token format")
	}

	return parts[1], nil // 返回密码部分
}
