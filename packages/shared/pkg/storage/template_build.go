package storage

import (
	"context"
	"fmt"
	
	"github.com/e2b-dev/infra/packages/shared/pkg/consts"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage/gcs"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage/header"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage/s3"
	"golang.org/x/sync/errgroup"
)

type TemplateBuild struct {
	files *TemplateFiles

	memfileHeader *header.Header
	rootfsHeader  *header.Header

	bucket *gcs.BucketHandle
	s3     *s3.BucketHandle
}

func NewTemplateBuild(
	memfileHeader *header.Header,
	rootfsHeader *header.Header,
	files *TemplateFiles,
) *TemplateBuild {
	if consts.CloudProviderEnv == consts.AWS {
		return &TemplateBuild{
			s3:    s3.GetTemplateBucket(),
			files: files,
		}
	}

	if consts.CloudProviderEnv == consts.GCP {
		return &TemplateBuild{
			bucket: gcs.GetTemplateBucket(),
			files:  files,
		}
	}

	panic(fmt.Sprintf("not implemented for cloud provider: %s", consts.CloudProviderEnv))
}

func (t *TemplateBuild) Remove(ctx context.Context) error {
	err := gcs.RemoveDir(ctx, t.bucket, t.files.StorageDir())
	if err != nil {
		return fmt.Errorf("error when removing template build '%s': %w", t.files.StorageDir(), err)
	}

	return nil
}

func (t *TemplateBuild) uploadMemfileHeader(ctx context.Context, h *header.Header) error {
	object := gcs.NewObject(ctx, t.bucket, t.files.StorageMemfileHeaderPath())

	serialized, err := header.Serialize(h.Metadata, h.Mapping)
	if err != nil {
		return fmt.Errorf("error when serializing memfile header: %w", err)
	}

	_, err = object.ReadFrom(serialized)
	if err != nil {
		return fmt.Errorf("error when uploading memfile header: %w", err)
	}

	return nil
}

func (t *TemplateBuild) uploadMemfile(ctx context.Context, memfilePath string) error {

	if consts.CloudProviderEnv == consts.AWS {
		fmt.Printf("Uploading memfile to s3://%s/%s\n", t.s3.Name(),t.files.StorageMemfilePath())
		object := s3.NewObject(t.s3.Client(), t.s3.Name(), t.files.StorageMemfilePath())
		object.UploadWithCli(ctx, memfilePath)
	}

	if consts.CloudProviderEnv == consts.GCP {
		object := gcs.NewObject(ctx, t.bucket, t.files.StorageMemfilePath())

		err := object.UploadWithCli(ctx, memfilePath)
		if err != nil {
			return fmt.Errorf("error when uploading memfile: %w", err)
		}
	}

	return nil
}

func (t *TemplateBuild) uploadRootfsHeader(ctx context.Context, h *header.Header) error {

	if consts.CloudProviderEnv == consts.AWS {
		fmt.Printf("Uploading rootfsheader to s3: %s\n", t.s3.Name(),t.files.StorageRootfsHeaderPath())
		object := s3.NewObject(t.s3.Client(), t.s3.Name(), t.files.StorageRootfsHeaderPath())
		object.UploadWithCli(ctx, t.files.StorageRootfsHeaderPath())
	}

	if consts.CloudProviderEnv == consts.GCP {
		object := gcs.NewObject(ctx, t.bucket, t.files.StorageRootfsHeaderPath())

		serialized, err := header.Serialize(h.Metadata, h.Mapping)
		if err != nil {
			return fmt.Errorf("error when serializing memfile header: %w", err)
		}

		_, err = object.ReadFrom(serialized)
		if err != nil {
			return fmt.Errorf("error when uploading memfile header: %w", err)
		}
	}

	return nil
}

func (t *TemplateBuild) uploadRootfs(ctx context.Context, rootfsPath string) error {

	if consts.CloudProviderEnv == consts.AWS {
		fmt.Printf("Uploading rootfs to s3: %s\n", t.s3.Name(),t.files.StorageRootfsPath())
		object := s3.NewObject(t.s3.Client(), t.s3.Name(), t.files.StorageRootfsPath())
		object.UploadWithCli(ctx, rootfsPath)
	}

	if consts.CloudProviderEnv == consts.GCP {
		object := gcs.NewObject(ctx, t.bucket, t.files.StorageRootfsPath())

		err := object.UploadWithCli(ctx, rootfsPath)
		if err != nil {
			return fmt.Errorf("error when uploading rootfs: %w", err)
		}
	}

	return nil
}

// Snapfile is small enough so we dont use composite upload.
func (t *TemplateBuild) uploadSnapfile(ctx context.Context, snapfilePath string) error {

	if consts.CloudProviderEnv == consts.AWS {
		fmt.Println("uploading snapfile to s3", t.files.StorageSnapfilePath())
		
		object := s3.NewObject(t.s3.Client(), t.s3.Name(), t.files.StorageSnapfilePath())
		object.UploadWithCli(ctx, snapfilePath)
	}

	if consts.CloudProviderEnv == consts.GCP {
		object := gcs.NewObject(ctx, t.bucket, t.files.StorageSnapfilePath())

		err := object.UploadWithCli(ctx, snapfilePath)
		if err != nil {
			return fmt.Errorf("error when uploading rootfs: %w", err)
		}
	}

	return nil
}

func (t *TemplateBuild) Upload(
	ctx context.Context,
	snapfilePath string,
	memfilePath *string,
	rootfsPath *string,
) chan error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if t.rootfsHeader == nil {
			return nil
		}

		err := t.uploadRootfsHeader(ctx, t.rootfsHeader)
		if err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		if rootfsPath == nil {
			return nil
		}

		err := t.uploadRootfs(ctx, *rootfsPath)
		if err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		if t.memfileHeader == nil {
			return nil
		}

		err := t.uploadMemfileHeader(ctx, t.memfileHeader)
		if err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		if memfilePath == nil {
			return nil
		}

		err := t.uploadMemfile(ctx, *memfilePath)
		if err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {

		if snapfilePath == "" {
			return nil
		}

		// snapfile, err := os.Open(snapfilePath)
		// if err != nil {
		// 	return err
		// }

		// defer snapfile.Close()

		err := t.uploadSnapfile(ctx, snapfilePath)
		if err != nil {
			return err
		}

		return nil
	})

	done := make(chan error)

	go func() {
		done <- eg.Wait()
	}()

	return done
}
