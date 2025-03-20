package s3

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// RemoveDir recursively removes all objects with the specified prefix
func RemoveDir(ctx context.Context, bucket *BucketHandle, dir string) error {
	// Ensure dir ends with a slash
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	// List objects with the prefix
	paginator := s3.NewListObjectsV2Paginator(bucket.Client(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket.Name()),
		Prefix: aws.String(dir),
	})

	// Collect all objects to delete
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("error when listing objects: %w", err)
		}

		if len(page.Contents) == 0 {
			continue
		}

		// Create a slice of objects to delete
		objectIds := make([]types.ObjectIdentifier, 0, len(page.Contents))
		for _, obj := range page.Contents {
			objectIds = append(objectIds, types.ObjectIdentifier{
				Key: obj.Key,
			})
		}

		// Delete objects in batches
		_, err = bucket.Client().DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(bucket.Name()),
			Delete: &types.Delete{
				Objects: objectIds,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return fmt.Errorf("error when deleting objects: %w", err)
		}
	}

	return nil
}
