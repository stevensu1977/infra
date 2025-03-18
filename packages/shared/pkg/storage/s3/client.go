package s3

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewClient(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.RetryMaxAttempts = 10
		o.RetryMode = aws.RetryModeStandard
	})

	return client, nil
}

func NewClientWithTimeout(ctx context.Context, timeout time.Duration) (*s3.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return NewClient(ctx)
}