package template

import (
	"context"
	"fmt"

	"github.com/e2b-dev/infra/packages/shared/pkg/consts"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage/gcs"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage/s3"
)

type Storage struct {
	bucket   *gcs.BucketHandle
	s3bucket *s3.BucketHandle
}

func NewStorage(ctx context.Context) *Storage {
	if consts.CloudProviderEnv == consts.AWS {
		return &Storage{
			bucket:   nil,
			s3bucket: s3.GetTemplateBucket(),
		}
	}

	if consts.CloudProviderEnv == consts.GCP {
		return &Storage{
			bucket:   gcs.GetTemplateBucket(),
			s3bucket: nil,
		}
	}

	panic(fmt.Sprintf("invalid cloud provider: %s", consts.CloudProviderEnv))
}

func (t *Storage) Remove(ctx context.Context, buildId string) error {

	if consts.CloudProviderEnv == consts.AWS {
		err := s3.RemoveDir(ctx, t.s3bucket, buildId)
		if err != nil {
			return fmt.Errorf("error when removing template '%s': %w", buildId, err)
		}
	}

	if consts.CloudProviderEnv == consts.GCP {
		err := gcs.RemoveDir(ctx, t.bucket, buildId)
		if err != nil {
			return fmt.Errorf("error when removing template '%s': %w", buildId, err)
		}
	}

	panic(fmt.Sprintf("invalid cloud provider: %s", consts.CloudProviderEnv))

	return nil
}

func (t *Storage) NewBuild(files *storage.TemplateFiles) *storage.TemplateBuild {
	return storage.NewTemplateBuild(nil, nil, files)
}
