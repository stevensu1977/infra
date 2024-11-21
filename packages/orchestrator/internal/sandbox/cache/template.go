package cache

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/e2b-dev/infra/packages/shared/pkg/storage"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage/block"
	"github.com/e2b-dev/infra/packages/shared/pkg/storage/gcs"
	"github.com/e2b-dev/infra/packages/shared/pkg/utils"
)

const (
	pageSize        = 2 << 11
	hugepageSize    = 2 << 20
	rootfsBlockSize = 2 << 11
)

type Template struct {
	files *storage.TemplateCacheFiles

	memfile  *utils.SetOnce[block.ReadonlyDevice]
	rootfs   *utils.SetOnce[block.ReadonlyDevice]
	snapfile *utils.SetOnce[*File]

	hugePages bool
}

func (t *Template) PageSize() int64 {
	if t.hugePages {
		return hugepageSize
	}

	return pageSize
}

func (t *TemplateCache) newTemplate(
	cacheIdentifier,
	templateId,
	buildId,
	kernelVersion,
	firecrackerVersion string,
	hugePages bool,
) *Template {
	files := storage.NewTemplateFiles(
		templateId,
		buildId,
		kernelVersion,
		firecrackerVersion,
	).NewTemplateCacheFiles(cacheIdentifier)

	return &Template{
		hugePages: hugePages,
		files:     files,
		memfile:   utils.NewSetOnce[block.ReadonlyDevice](),
		rootfs:    utils.NewSetOnce[block.ReadonlyDevice](),
		snapfile:  utils.NewSetOnce[*File](),
	}
}

func (t *Template) Fetch(ctx context.Context, bucket *gcs.BucketHandle) {
	err := os.MkdirAll(t.files.CacheDir(), os.ModePerm)
	if err != nil {
		errMsg := fmt.Errorf("failed to create directory %s: %w", t.files.CacheDir(), err)

		t.memfile.SetError(errMsg)
		t.rootfs.SetError(errMsg)
		t.snapfile.SetError(errMsg)

		return
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() error {
		defer wg.Done()

		snapfile, snapfileErr := NewFile(
			ctx,
			bucket,
			t.files.StorageSnapfilePath(),
			t.files.CacheSnapfilePath(),
		)
		if snapfileErr != nil {
			errMsg := fmt.Errorf("failed to fetch snapfile: %w", snapfileErr)

			return t.snapfile.SetError(errMsg)
		}

		return t.snapfile.SetValue(snapfile)
	}()

	wg.Add(1)
	go func() error {
		defer wg.Done()

		memfileStorage, memfileErr := block.NewStorage(
			ctx,
			bucket,
			t.files.StorageMemfilePath(),
			t.PageSize(),
			t.files.CacheMemfilePath(),
		)
		if memfileErr != nil {
			errMsg := fmt.Errorf("failed to create memfile storage: %w", memfileErr)

			return t.memfile.SetError(errMsg)
		}

		return t.memfile.SetValue(memfileStorage)
	}()

	wg.Add(1)
	go func() error {
		defer wg.Done()

		rootfsStorage, rootfsErr := block.NewStorage(
			ctx,
			bucket,
			t.files.StorageRootfsPath(),
			// TODO: This should ideally be the blockSize (4096), but we would need to implement more complex dirty block caching in cache there.
			ChunkSize,
			t.files.CacheRootfsPath(),
		)
		if rootfsErr != nil {
			errMsg := fmt.Errorf("failed to create rootfs storage: %w", rootfsErr)

			return t.rootfs.SetError(errMsg)
		}

		return t.rootfs.SetValue(rootfsStorage)
	}()

	wg.Wait()
}

func (t *Template) Close() error {
	var errs []error

	memfile, err := t.Memfile()
	if err == nil {
		errs = append(errs, memfile.Close())
	}

	rootfs, err := t.Rootfs()
	if err == nil {
		errs = append(errs, rootfs.Close())
	}

	snapfile, err := t.Snapfile()
	if err == nil {
		errs = append(errs, snapfile.Close())
	}

	return errors.Join(errs...)
}

func (t *Template) Files() *storage.TemplateCacheFiles {
	return t.files
}

func (t *Template) Memfile() (block.ReadonlyDevice, error) {
	return t.memfile.Wait()
}

func (t *Template) Rootfs() (block.ReadonlyDevice, error) {
	return t.rootfs.Wait()
}

func (t *Template) Snapfile() (*File, error) {
	return t.snapfile.Wait()
}
