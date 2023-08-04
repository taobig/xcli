package xcli

import (
	"github.com/kardianos/service"
	"github.com/urfave/cli/v2"
)

type SystemService struct {
	cCtx *cli.Context

	startCallback func(cCtx *cli.Context)
	stopCallback  func() error
}

func (ss *SystemService) Start(s service.Service) error {
	go ss.run()
	return nil
}

func (ss *SystemService) run() {
	ss.startCallback(ss.cCtx)
}

func (ss *SystemService) Stop(s service.Service) error {
	return ss.stopCallback()
}
