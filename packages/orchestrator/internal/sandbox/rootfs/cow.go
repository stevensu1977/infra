package rootfs

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/e2b-dev/infra/packages/orchestrator/internal/sandbox/block"
	"github.com/e2b-dev/infra/packages/orchestrator/internal/sandbox/nbd"
	"github.com/e2b-dev/infra/packages/shared/pkg/utils"
)

type CowDevice struct {
	overlay block.Device
	mnt     *nbd.DirectPathMount

	ready *utils.SetOnce[string]
}

func NewCowDevice(rootfs block.ReadonlyDevice, cachePath string, blockSize int64) (*CowDevice, error) {
	size, err := rootfs.Size()
	if err != nil {
		return nil, fmt.Errorf("error getting device size: %w", err)
	}

	cache, err := block.NewCache(size, blockSize, cachePath)
	if err != nil {
		return nil, fmt.Errorf("error creating cache: %w", err)
	}

	overlay := block.NewOverlay(rootfs, cache, blockSize)

	mnt := nbd.NewDirectPathMount(overlay)

	return &CowDevice{
		mnt:     mnt,
		overlay: overlay,
		ready:   utils.NewSetOnce[string](),
	}, nil
}

func (o *CowDevice) Start(ctx context.Context) error {
	deviceIndex, err := o.mnt.Open(ctx)
	if err != nil {
		return o.ready.SetError(fmt.Errorf("error opening overlay file: %w", err))
	}

	return o.ready.SetValue(nbd.GetDevicePath(deviceIndex))
}

func (o *CowDevice) Export(ctx context.Context, path string) error {
	devicePath, err := o.ready.Wait()
	if err != nil {
		return fmt.Errorf("error getting overlay path: %w", err)
	}

	_, err = exec.CommandContext(ctx, "cp", devicePath, path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error exporting overlay: %w", err)
	}

	return nil
}

func (o *CowDevice) Close() error {
	var errs []error

	err := o.mnt.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing overlay mount: %w", err))
	}

	err = o.overlay.Close()
	if err != nil {
		errs = append(errs, fmt.Errorf("error closing overlay cache: %w", err))
	}

	devicePath, err := o.ready.Wait()
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting overlay path: %w", err))

		return errors.Join(errs...)
	}

	slot, err := nbd.GetDeviceSlot(devicePath)
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting overlay slot: %w", err))

		return errors.Join(errs...)
	}

	counter := 0
	for {
		counter++
		err := nbd.Pool.ReleaseDevice(slot)
		if errors.Is(err, nbd.ErrDeviceInUse{}) {
			if counter%100 == 0 {
				log.Printf("[%dth try] error releasing overlay device: %v\n", counter, err)
			}

			continue
		}

		if err != nil {
			return fmt.Errorf("error releasing overlay device: %w", err)
		}

		break
	}

	return nil
}

func (o *CowDevice) Path() (string, error) {
	return o.ready.Wait()
}
