package xcli

import (
	"github.com/kardianos/service"
)

type SystemService struct {
	startCallback func()
	stopCallback  func() error
}

func (ss *SystemService) Start(s service.Service) error {
	go ss.run()
	return nil
}

func (ss *SystemService) run() {
	ss.startCallback()
}

func (ss *SystemService) Stop(s service.Service) error {
	return ss.stopCallback()
}
