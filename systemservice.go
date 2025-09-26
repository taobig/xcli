package xcli

import (
	"context"

	"github.com/kardianos/service"
	"github.com/urfave/cli/v3"
)

type SystemService struct {
	ctx context.Context
	cmd *cli.Command

	startCallback func(ctx context.Context, cmd *cli.Command)
	stopCallback  func() error
}

func (ss *SystemService) Start(s service.Service) error {
	go ss.run()
	return nil
}

func (ss *SystemService) run() {
	ss.startCallback(ss.ctx, ss.cmd)
}

func (ss *SystemService) Stop(s service.Service) error {
	return ss.stopCallback()
}
