package s3

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/e2b-dev/infra/packages/shared/pkg/utils"
)

// BucketHandle is a wrapper for S3 bucket operations
type BucketHandle struct {
	client *s3.Client
	name   string
}

var getClient = sync.OnceValue(func() *s3.Client {
	return utils.Must(NewClient(context.Background()))
})

func newBucket(bucket string) *BucketHandle {
	return &BucketHandle{
		client: getClient(),
		name:   bucket,
	}
}

func getTemplateBucketName() string {
	// Check if the environment variable is set
	bucketName := os.Getenv("TEMPLATE_BUCKET_NAME")
	if bucketName != "" {
		return bucketName
	}

	// For mock-sandbox or development environments, provide a fallback value
	if os.Getenv("MOCK_SANDBOX") == "true" || os.Getenv("DEV_ENV") == "true" {
		return "mock-template-bucket"
	}

	// In production, still require the environment variable
	return utils.RequiredEnv("TEMPLATE_BUCKET_NAME", "bucket for storing template files")
}

func GetTemplateBucket() *BucketHandle {
	return newBucket(getTemplateBucketName())
}

// Name returns the name of the bucket
func (b *BucketHandle) Name() string {
	return b.name
}

// Client returns the S3 client
func (b *BucketHandle) Client() *s3.Client {
	return b.client
}

// Object returns a handle to interact with a specific object in this bucket
func (b *BucketHandle) Object(name string) *Object {
	return &Object{
		client: b.client,
		bucket: b.name,
		name:   name,
		retry:  true,
	}
}
