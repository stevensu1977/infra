package filesystem

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"connectrpc.com/connect"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"

	"github.com/e2b-dev/infra/packages/envd/internal/logs"
	"github.com/e2b-dev/infra/packages/envd/internal/permissions"
	rpc "github.com/e2b-dev/infra/packages/envd/internal/services/spec/filesystem"
	"github.com/e2b-dev/infra/packages/shared/pkg/id"
)

type FileWatcher struct {
	watcher *fsnotify.Watcher
	Events  []*rpc.FilesystemEvent
	ctx     context.Context
	Error   error
}

func CreateFileWatcher(watchPath, operationID string, logger *zerolog.Logger) (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("error creating watcher: %w", err))
	}

	err = w.Add(watchPath)
	if err != nil {
		_ = w.Close()
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("error adding path %s to watcher: %w", watchPath, err))
	}
	fw := &FileWatcher{
		watcher: w,
		ctx:     context.Background(),
		Events:  []*rpc.FilesystemEvent{},
		Error:   nil,
	}

	go func() {
		for {
			select {
			case <-fw.ctx.Done():
				return
			case chErr, ok := <-w.Errors:
				if !ok {
					fw.Error = connect.NewError(connect.CodeInternal, fmt.Errorf("watcher error channel closed"))
					return
				}

				fw.Error = connect.NewError(connect.CodeInternal, fmt.Errorf("watcher error: %w", chErr))
				return
			case e, ok := <-w.Events:
				if !ok {
					fw.Error = connect.NewError(connect.CodeInternal, fmt.Errorf("watcher event channel closed"))
					return
				}

				// One event can have multiple operations.
				ops := []rpc.EventType{}

				if fsnotify.Create.Has(e.Op) {
					ops = append(ops, rpc.EventType_EVENT_TYPE_CREATE)
				}

				if fsnotify.Rename.Has(e.Op) {
					ops = append(ops, rpc.EventType_EVENT_TYPE_RENAME)
				}

				if fsnotify.Chmod.Has(e.Op) {
					ops = append(ops, rpc.EventType_EVENT_TYPE_CHMOD)
				}

				if fsnotify.Write.Has(e.Op) {
					ops = append(ops, rpc.EventType_EVENT_TYPE_WRITE)
				}

				if fsnotify.Remove.Has(e.Op) {
					ops = append(ops, rpc.EventType_EVENT_TYPE_REMOVE)
				}

				for _, op := range ops {
					name, nameErr := filepath.Rel(watchPath, e.Name)
					if nameErr != nil {
						fw.Error = connect.NewError(connect.CodeInternal, fmt.Errorf("error getting relative path: %w", nameErr))
						return
					}

					filesystemEvent := &rpc.WatchDirResponse_Filesystem{
						Filesystem: &rpc.FilesystemEvent{
							Name: name,
							Type: op,
						},
					}

					event := &rpc.WatchDirResponse{
						Event: filesystemEvent,
					}

					fw.Events = append(fw.Events, &rpc.FilesystemEvent{
						Name: name,
						Type: op,
					})

					logger.
						Debug().
						Str("event_type", "filesystem_event").
						Str(string(logs.OperationIDKey), operationID).
						Interface("filesystem_event", event).
						Msg("Streaming filesystem event")
				}
			}
		}
	}()

	return fw, nil
}

func (fw *FileWatcher) Close() {
	_ = fw.watcher.Close()
	fw.ctx.Done()
}

func (s Service) WatchDirStart(ctx context.Context, req *connect.Request[rpc.WatchDirRequest]) (*connect.Response[rpc.WatchDirStartResponse], error) {
	u, err := permissions.GetAuthUser(ctx)
	if err != nil {
		return nil, err
	}

	watchPath, err := permissions.ExpandAndResolve(req.Msg.GetPath(), u)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	info, err := os.Stat(watchPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("path %s not found: %w", watchPath, err))
		}

		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("error statting path %s: %w", watchPath, err))
	}

	if !info.IsDir() {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("path %s not a directory: %w", watchPath, err))
	}

	watcherId := "w" + id.Generate()

	w, err := CreateFileWatcher(watchPath, watcherId, s.logger)
	s.watchers.Store(watcherId, w)

	return connect.NewResponse(&rpc.WatchDirStartResponse{
		WatcherId: watcherId,
	}), nil
}

func (s Service) WatchDirPoll(_ context.Context, req *connect.Request[rpc.WatchDirPollRequest]) (*connect.Response[rpc.WatchDirPollResponse], error) {
	watcherId := req.Msg.GetWatcherId()

	w, ok := s.watchers.Load(watcherId)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("watcher with id %s not found", watcherId))
	}

	if w.Error != nil {
		return nil, w.Error
	}

	events := w.Events
	if int(req.Msg.Offset) >= len(w.Events) {
		events = []*rpc.FilesystemEvent{}
	} else {
		events = w.Events[req.Msg.Offset:]
	}

	return connect.NewResponse(&rpc.WatchDirPollResponse{
		Events: events,
	}), nil
}

func (s Service) WatchDirStop(_ context.Context, req *connect.Request[rpc.WatchDirStopRequest]) (*connect.Response[rpc.WatchDirStopResponse], error) {
	watcherId := req.Msg.GetWatcherId()

	w, ok := s.watchers.Load(watcherId)
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("watcher with id %s not found", watcherId))
	}

	w.Close()
	s.watchers.Delete(watcherId)

	return connect.NewResponse(&rpc.WatchDirStopResponse{}), nil
}
