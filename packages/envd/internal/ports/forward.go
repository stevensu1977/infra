package ports

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"time"

	"github.com/rs/zerolog"
	psnet "github.com/shirou/gopsutil/v4/net"
)

type PortState string

const (
	PortStateForward PortState = "FORWARD"
	PortStateDelete  PortState = "DELETE"

	scanPeriod = 1 * time.Second
)

var (
	forwardedIP = "127.0.0.1"
	gatewayIP   = net.IPv4(169, 254, 0, 21)
)

type forwarding struct {
	cmd *exec.Cmd
}

func (f *forwarding) Stop() error {
	return f.cmd.Process.Kill()
}

type Forwarder struct {
	logger *zerolog.Logger
	ctx    context.Context

	ports map[uint32]*forwarding
}

func NewForwarder(
	ctx context.Context,
	logger *zerolog.Logger,
) *Forwarder {
	return &Forwarder{
		logger: logger,
		ctx:    ctx,
		ports:  make(map[uint32]*forwarding),
	}
}

func (f *Forwarder) Start() {
	ticker := time.NewTicker(scanPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-f.ctx.Done():
			return
		case <-ticker.C:
			cs, err := psnet.Connections("tcp")
			if err != nil {
				f.logger.Err(err).Msg("failed to get connections")

				return
			}

			newForwarding := make(map[uint32]*forwarding)

			for _, conn := range cs {
				if conn.Laddr.IP != forwardedIP {
					continue
				}

				forwarding, ok := f.ports[conn.Laddr.Port]
				if ok {
					newForwarding[conn.Laddr.Port] = forwarding

					delete(f.ports, conn.Laddr.Port)

					continue
				}

				forwarding, forwardErr := f.forwardPort(conn.Laddr.Port)
				if forwardErr != nil {
					f.logger.Err(forwardErr).Msg("failed to forward port")

					continue
				}

				newForwarding[conn.Laddr.Port] = forwarding
			}

			for _, forwarding := range newForwarding {
				forwarding.Stop()
			}

			f.ports = newForwarding
		}
	}
}

func (f *Forwarder) forwardPort(port uint32) (*forwarding, error) {
	// https://unix.stackexchange.com/questions/311492/redirect-application-listening-on-localhost-to-listening-on-external-interface
	// socat -d -d TCP4-LISTEN:4000,bind=169.254.0.21,fork TCP4:localhost:4000
	socatCmd := fmt.Sprintf(
		"socat -d -d -d TCP4-LISTEN:%v,bind=%s,fork TCP4:localhost:%v",
		port,
		gatewayIP.String(),
		port,
	)

	cmd := exec.CommandContext(f.ctx, "/bin/bash", "-c", socatCmd)

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start port forwarding - failed to start socat: %w", err)
	}

	go cmd.Wait()

	return &forwarding{
		cmd: cmd,
	}, nil
}
